package amdfw

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const FETSignature = uint32(0x55AA55AA)
const FETDefaultOffset = uint32(0x20000)

type FirmwareEntryTable struct {
	Location uint32

	Signature     uint32
	ImcRomBase    *uint32
	GecRomBase    *uint32
	XHCRomBase    *uint32
	PSPDirBase    *uint32
	NewPSPDirBase *uint32
	BHDDirBase    *uint32
	NewBHDDirBase *uint32
}

type binaryFet struct {
	Signature     uint32
	ImcRomBase    uint32
	GecRomBase    uint32
	XHCRomBase    uint32
	PSPDirBase    uint32
	NewPSPDirBase uint32
	BHDDirBase    uint32
	NewBHDDirBase uint32
}

// Looks for the FET Signature at the often used offsets.
func FindFirmwareEntryTable(firmware []byte) (uint32, error) {

	for _, addr := range []uint32{FETDefaultOffset, 0, 0x820000, 0xC20000, 0xE20000, 0xF20000} {
		err := checkValidFirmwareEntryTable(firmware, uint32(addr))
		if err == nil {
			return addr, nil
		}
	}
	return 0, fmt.Errorf("No FirmwareTable found")
}

// Looks for the FET Signature everywhere (slow)
func FindFirmwareEntryTableByScan(firmware []byte) (uint32, error) {

	end := len(firmware)
	for addr := 0; addr <= end-4; addr++ {
		err := checkValidFirmwareEntryTable(firmware, uint32(addr))
		if err == nil {
			return uint32(addr), nil
		}
	}
	return 0, fmt.Errorf("No FirmwareTable found")
}

func checkValidFirmwareEntryTable(firmware []byte, address uint32) error {

	if int(address) > len(firmware)+binary.Size(FirmwareEntryTable{}) {
		return fmt.Errorf("Not AMD Table Header: Address out of bounds")
	}
	potentialMagic := binary.LittleEndian.Uint32(firmware[address:])
	if potentialMagic == FETSignature {
		return nil
	}

	return fmt.Errorf("Not AMD Table Header")
}

// Parses the FET at the given Address
// Default address is 0x20000 but some PSP versions appear to be able to search on different offsets
func ParseFirmwareEntryTable(firmware []byte, address uint32) (*FirmwareEntryTable, error) {

	if err := checkValidFirmwareEntryTable(firmware, address); err != nil {
		return nil, fmt.Errorf("Could not find FirmwareEntryTable Signature: %v", err)
	}

	tempTable := binaryFet{}

	if err := binary.Read(bytes.NewReader(firmware[address:]), binary.LittleEndian, &tempTable); err != nil {
		return nil, fmt.Errorf("Could not read FirmwareEntryTable: %v", err)
	}

	//TODO check validity
	fet := FirmwareEntryTable{
		Signature:     tempTable.Signature,
		ImcRomBase:    &tempTable.ImcRomBase,
		GecRomBase:    &tempTable.GecRomBase,
		XHCRomBase:    &tempTable.XHCRomBase,
		PSPDirBase:    &tempTable.PSPDirBase,
		NewPSPDirBase: &tempTable.NewPSPDirBase,
		BHDDirBase:    &tempTable.BHDDirBase,
		NewBHDDirBase: &tempTable.NewBHDDirBase,
		Location:      uint32(address),
	}

	return &fet, nil
}

// Writes FET into existing image
func (fet *FirmwareEntryTable) Write(baseImage []byte, address uint32) error {

	tempTable := binaryFet{
		Signature:     fet.Signature,
		ImcRomBase:    *fet.ImcRomBase,
		GecRomBase:    *fet.GecRomBase,
		XHCRomBase:    *fet.XHCRomBase,
		PSPDirBase:    *fet.PSPDirBase,
		NewPSPDirBase: *fet.NewPSPDirBase,
		BHDDirBase:    *fet.BHDDirBase,
		NewBHDDirBase: *fet.NewBHDDirBase,
	}

	fetSize := binary.Size(tempTable)

	if len(baseImage) < int(address)+fetSize {
		return fmt.Errorf("BaseImage to small to insert FET")
	}

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, tempTable)
	if err != nil {
		return fmt.Errorf("Writing binary failed: %v", err)
	}

	bytesCopied := copy(baseImage[address:], buf.Bytes())

	if bytesCopied != fetSize {
		return fmt.Errorf("Writing binary failed: Not all Bytes copied!")
	}

	return nil
}
