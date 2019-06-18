package amdfw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	pspentries "github.com/Mimoja/PSP-Entry-Types"
)

type (
	Entry struct {
		DirectoryEntry DirectoryEntry
		Header         *EntryHeader
		Raw            []byte
		Signature      []byte
		Comment        []string
		TypeInfo       *TypeInfo
	}

	EntryHeader struct {
		Unknown00      [0x10]byte // 0x00
		ID             uint32     // 0x10
		SizeSigned     uint32     // 0x14
		IsEncrypted    uint32     // 0x18
		Unknown1C      uint32     // 0x1C
		EncFingerprint [0x10]byte // 0x20
		IsSigned       uint32     // 0x30
		Unknown34      uint32     // 0x34
		SigFingerprint [0x10]byte // 0x38
		IsCompressed   uint32     // 0x48
		Unknown4C      uint32     // 0x4C
		FullSize       uint32     // 0x50
		Unknown54      uint32     // 0x54
		Unknown58      [0x08]byte // 0x58
		Version        [0x04]byte // 0x60
		Unknown64      uint32     // 0x64
		Unknown68      uint32     // 0x68
		SizePacked     uint32     // 0x6C
		Unknown70      [0x10]byte // 0x70
		Unknown80      [0x10]byte // 0x80
		Unknown90      uint32     // 0x90
		Unknown94      uint32     // 0x94
		Unknown98      uint32     // 0x98
		Unknown9C      uint32     // 0x9C
		UnknownA0      uint32     // 0xA0
		UnknownA4      uint32     // 0xA4
		UnknownA8      uint32     // 0xA8
		UnknownAC      uint32     // 0xAC
		UnknownB0      [0x50]byte // 0xB0

	}

	TypeInfo struct {
		Name    string
		Comment string
	}
)

var knownTypes = pspentries.Types()

func ParseEntry(firmwareBytes []byte, directoryEntry DirectoryEntry, flashMapping uint32) (*Entry, error) {
	entry := Entry{
		DirectoryEntry: directoryEntry,
	}

	/**
	 *	Typechecking
	 */
	for _, knownType := range knownTypes {
		if knownType.Type == directoryEntry.Type {
			name := knownType.Name
			if name == "" {
				name = knownType.ProposedName
			}
			info := TypeInfo{
				Name:    name,
				Comment: knownType.Comment,
			}
			entry.TypeInfo = &info
			break
		}
	}

	if entry.TypeInfo == nil {
		errorAndComment(&entry, fmt.Errorf("Unknown Type: 0x%08X", directoryEntry.Type))
	}

	/**
	 * Raw Data logic
	 */

	location := directoryEntry.Location
	size := directoryEntry.Size

	if location >= flashMapping {
		location -= flashMapping
	}

	if int(location) > len(firmwareBytes) {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: Location out of bounds (0x%08X)", location))
	}

	if int(size) > len(firmwareBytes)-int(location) {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: Size to big (0x%08X)", size))
	}

	entryBytes := firmwareBytes[location : location+size]
	entry.Raw = entryBytes

	/**
	 * Header Parsing
	 */

	if size < 0x100 && size > 0 {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: Entry to small for header parsing: (0x08%X) bytes", size))
	}

	headerBytes := entryBytes[:0x100]

	if allOneValue(headerBytes) {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: All Fields are 0x%02X", headerBytes[1]))
	}

	header := EntryHeader{}
	err := binary.Read(bytes.NewReader(entryBytes), binary.LittleEndian, &header)

	if err != nil {
		return errorAndComment(&entry, fmt.Errorf("Error: Could not read header: %v", err))
	}

	if header.IsCompressed > 1 {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: Compressed Filed is 0x%02X", header.IsCompressed))
	}

	if header.SizePacked == 0 &&
		header.SizeSigned == 0 &&
		header.FullSize == 0 {
		return errorAndComment(&entry, fmt.Errorf("Not a parsable Entry: Size Values not reasonable"))
	}

	entry.Header = &header
	entry.Signature = entryBytes[size-2048/8:]

	return &entry, nil
}

func (entry Entry) Write(baseImage []byte, address uint32) error {
	copied := copy(baseImage[address:], entry.Raw)

	if copied != len(entry.Raw) {
		return fmt.Errorf("Could not write Entry: Failed after 0x08%X Bytes", copied)
	}
	return nil
}

func allOneValue(s []byte) bool {
	reference := s[0]
	for _, v := range s {
		if v != reference {
			return false
		}
	}
	return true
}

func errorAndComment(entry *Entry, err error) (*Entry, error) {
	entry.Comment = append(entry.Comment, err.Error())
	return entry, err
}
