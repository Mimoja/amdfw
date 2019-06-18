package amdfw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const PSPCOOCKIE = "$PSP"
const DUALPSPCOOCKIE = "2PSP"
const BHDCOOCKIE = "$BHD"
const SECONDPSPCOOCKIE = "$PL2"
const SECONDBHDCOOCKIE = "$BL2"

type (
	Directory struct {
		Header   DirectoryHeader
		Entries  []Entry
		Location uint32
	}

	DirectoryHeader struct {
		Cookie       [4]byte
		Checksum     uint32
		TotalEntries uint32
		Reserved     uint32
	}

	DirectoryEntry struct {
		Type     uint32
		Size     uint32
		Location uint32
		Reserved uint32
		Unknown  *uint64
	}

	binaryDirectoryEntry struct {
		Type     uint32
		Size     uint32
		Location uint32
		Reserved uint32
	}
)

func ParseDirectory(firmwareBytes []byte, address uint32, flashMapping uint32) (*Directory, error) {
	directory := Directory{}

	if address > flashMapping {
		address -= flashMapping
	}

	if address > uint32(len(firmwareBytes)) {
		return nil, fmt.Errorf("Firmwarebytes not long enough for reading directory header..")
	}

	directory.Location = address
	directoryBytes := firmwareBytes[address:]
	directoryReader := bytes.NewReader(directoryBytes)

	directoryHeader := DirectoryHeader{}
	err := binary.Read(directoryReader, binary.LittleEndian, &directoryHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not read directory header: %v", err)
	}

	cookie := string(directoryHeader.Cookie[:])

	isCookieKnown := false

	for _, c := range []string{PSPCOOCKIE, DUALPSPCOOCKIE, SECONDPSPCOOCKIE, BHDCOOCKIE, SECONDBHDCOOCKIE} {
		if isCookieKnown = c == cookie; isCookieKnown {
			break
		}
	}

	if !isCookieKnown {
		return nil, fmt.Errorf("No Valid Cookie at start of directory: %v", directoryHeader.Cookie[:])
	}

	directory.Header = directoryHeader

	if int(directory.Header.TotalEntries)*4*8 > len(firmwareBytes[address+8:]) {
		return nil, fmt.Errorf("Could not read directory entries from image: Too many entries(%d)!", directory.Header.TotalEntries)
	}

	if cookie == DUALPSPCOOCKIE {
		if _, err := directoryReader.Seek(16, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("Could not read 2PSP directory header: %v", err)
		}
	}

	directory.Entries = make([]Entry, directory.Header.TotalEntries)

	for i := 0; i < int(directory.Header.TotalEntries); i++ {
		directoryEntry := DirectoryEntry{}
		binDirEntry := binaryDirectoryEntry{}

		if err := binary.Read(directoryReader, binary.LittleEndian, &binDirEntry); err != nil {
			return nil, fmt.Errorf("Coulf not read directory directoryEntry")
		}

		directoryEntry.Type = binDirEntry.Type
		directoryEntry.Size = binDirEntry.Size
		directoryEntry.Location = binDirEntry.Location
		directoryEntry.Reserved = binDirEntry.Reserved

		entry, _ := ParseEntry(firmwareBytes, directoryEntry, flashMapping)

		if cookie == DUALPSPCOOCKIE {
			entry.TypeInfo.Name = "PSP_DIRECTORY"
			entry.TypeInfo.Comment = "Full PSP Directory"
		} else if cookie == BHDCOOCKIE || cookie == SECONDBHDCOOCKIE {
			//BHD Entries adds 2 additional bytes
			unknownBytes := make([]byte, 8)
			if c, err := directoryReader.Read(unknownBytes); err != nil || c != 8 {
				return nil, fmt.Errorf("Could not read BHD directory entry: %v", err)
			}
			val := binary.LittleEndian.Uint64(unknownBytes)

			entry.DirectoryEntry.Unknown = &val
		}

		directory.Entries[i] = *entry
	}

	return &directory, nil
}

func (header *DirectoryHeader) Write(baseImage []byte, address uint32) error {
	buf := new(bytes.Buffer)

	if int(address)+binary.Size(header) > len(baseImage) {
		return fmt.Errorf("Writing DirectoryHeader failed: BaseImage to small")

	}

	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		return fmt.Errorf("Writing binary failed: %v", err)
	}

	copy(baseImage[address:], buf.Bytes())

	return nil
}

func (entry *DirectoryEntry) Write(baseImage []byte, address uint32) error {
	buf := new(bytes.Buffer)

	binEntry := binaryDirectoryEntry{
		Type:     entry.Type,
		Size:     entry.Size,
		Location: entry.Location,
		Reserved: entry.Reserved,
	}

	if int(address)+binary.Size(binEntry) > len(baseImage) {
		return fmt.Errorf("Writing DirectoryHeader failed: BaseImage to small")

	}

	err := binary.Write(buf, binary.LittleEndian, binEntry)
	if err != nil {
		return fmt.Errorf("Writing binary failed: %v", err)
	}

	bytesNeeded := binary.Size(binaryDirectoryEntry{})

	if entry.Unknown != nil {
		binary.Write(buf, binary.LittleEndian, *entry.Unknown)
		bytesNeeded += binary.Size(*entry.Unknown)
	}

	bytesCopied := copy(baseImage[address:], buf.Bytes())

	if bytesCopied != bytesNeeded {
		return fmt.Errorf("Writing binary failed: Not all Bytes copied!")
	}

	return nil
}

func (directory *Directory) Write(baseImage []byte, flashMapping uint32) error {
	if int(directory.Location) > len(baseImage) {
		return fmt.Errorf("BaseImage to small to insert Directory")
	}

	err := directory.Header.Write(baseImage, directory.Location)
	if err != nil {
		return err
	}

	location := directory.Location + 0x10

	cookie := string(directory.Header.Cookie[:])

	if cookie == DUALPSPCOOCKIE {
		location += 0x10
	}

	for i, entry := range directory.Entries {

		entryLength := uint32(16)

		if entry.DirectoryEntry.Unknown != nil {
			entryLength += 8
		}

		entryAddress := location + uint32(i)*entryLength
		err := entry.DirectoryEntry.Write(baseImage, entryAddress)
		if err != nil {
			return err
		}

		entryLocation := entry.DirectoryEntry.Location

		err = entry.Write(baseImage, entryLocation&^flashMapping)
		if err != nil {
			return err
		}
	}
	return nil
}

func fletcher32(data []byte) uint32 {
	c0 := 0xFFFF
	c1 := 0xFFFF

	count := len(data)
	for index := 0; index < count; index += 2 {
		next := binary.LittleEndian.Uint16(data[index:])
		c0 += int(next)
		c1 += c0
		if index > 255 || index == count-2 {
			c0 = (c0 & 0xFFFF) + (c0 >> 16)
			c1 = (c1 & 0xFFFF) + (c1 >> 16)
		}
	}
	return uint32((c1 << 16) | c0)
}

// Validates the Directory Checksum and return the actual value
func (directory *Directory) ValidateChecksum() (valid bool, actual uint32) {
	uint32Buffer := make([]byte, 4)
	buf := new(bytes.Buffer)

	//Rest of the header

	binary.LittleEndian.PutUint32(uint32Buffer, directory.Header.TotalEntries)
	buf.Write(uint32Buffer)

	binary.LittleEndian.PutUint32(uint32Buffer, directory.Header.Reserved)
	buf.Write(uint32Buffer)

	cookie := string(directory.Header.Cookie[:])
	if cookie == DUALPSPCOOCKIE {
		buf.Write(bytes.Repeat([]byte{0}, 16))
	}

	for _, entry := range directory.Entries {

		binary.LittleEndian.PutUint32(uint32Buffer, entry.DirectoryEntry.Type)
		buf.Write(uint32Buffer)

		binary.LittleEndian.PutUint32(uint32Buffer, entry.DirectoryEntry.Size)
		buf.Write(uint32Buffer)

		binary.LittleEndian.PutUint32(uint32Buffer, entry.DirectoryEntry.Location)
		buf.Write(uint32Buffer)

		binary.LittleEndian.PutUint32(uint32Buffer, entry.DirectoryEntry.Reserved)
		buf.Write(uint32Buffer)

		if cookie == BHDCOOCKIE || cookie == SECONDBHDCOOCKIE {
			uint64Buffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(uint64Buffer, *entry.DirectoryEntry.Unknown)
			buf.Write(uint64Buffer)
		}
	}

	sum := fletcher32(buf.Bytes())
	return sum == directory.Header.Checksum, sum
}
