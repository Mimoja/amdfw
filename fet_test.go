package amdfw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testImage16MB = 16777216

/** Demo FET
00020000  aa 55 aa 55 00 00 00 00  00 00 00 00 00 10 02 ff  |.U.U............|
00020010  00 10 10 ff 00 00 10 00  00 10 1c ff 00 d0 2b ff  |..............+.|
00020020  00 90 3e 00 fe ff ff ff  ff ff ff ff ff ff ff ff  |..>.............|
00020030  00 00 00 00 00 00 9a ff  ff ff ff ff ff ff ff ff  |................|
00020040  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  |................|
*/

var (
	fetBytes = []byte{
		0xaa, 0x55, 0xaa, 0x55, 0x00, 0x00, 0x00, 0x00, 0x67, 0x45, 0x23, 0x01, 0x00, 0x10, 0x02, 0xff,
		0x00, 0x10, 0x10, 0xff, 0x00, 0x00, 0x10, 0x00, 0x00, 0x10, 0x1c, 0xff, 0x00, 0xd0, 0x2b, 0xff,
	}

	testImcRomBase    uint32 = 0
	testGecRomBase    uint32 = 0x01234567
	testXHCRomBase    uint32 = 0xFF021000
	testPSPDirBase    uint32 = 0xFF101000
	testNewPSPDirBase uint32 = 0x100000
	testBHDDirBase    uint32 = 0xFF1C1000
	testUnknown1      uint32 = 0xFF2BD000

	testFet = FirmwareEntryTable{
		Signature:     0x55AA55AA,
		ImcRomBase:    &testImcRomBase,
		GecRomBase:    &testGecRomBase,
		XHCRomBase:    &testXHCRomBase,
		PSPDirBase:    &testPSPDirBase,
		NewPSPDirBase: &testNewPSPDirBase,
		BHDDirBase:    &testBHDDirBase,
		NewBHDDirBase: &testUnknown1,
		Location:      0x20000,
	}
)

func mockFetImage() []byte {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[FETDefaultOffset:], fetBytes)

	// Copy $PSP Header
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], PSPCOOCKIE)
	return imageBytes
}

func TestCheckValidFET(t *testing.T) {

	err := checkValidFirmwareEntryTable(mockFetImage(), FETDefaultOffset)

	assert.Nil(t, err)
}

func TestCheckValidFETFailSize(t *testing.T) {
	imageToSmall := make([]byte, FETDefaultOffset/2)

	err := checkValidFirmwareEntryTable(imageToSmall, FETDefaultOffset)

	assert.EqualError(t, err, "Not AMD Table Header: Address out of bounds")
}

func TestCheckValidFETFailInvalid(t *testing.T) {

	err := checkValidFirmwareEntryTable(mockFetImage(), 0x1234)

	assert.EqualError(t, err, "Not AMD Table Header")
}

func TestFindFirmwareEntryTableByScanFail(t *testing.T) {
	smallImage := make([]byte, 500)

	offset, err := FindFirmwareEntryTable(smallImage)

	assert.EqualError(t, err, "No FirmwareTable found")
	assert.Equal(t, uint32(0), offset)
}

func TestFindFirmwareEntryTableFail(t *testing.T) {
	smallImage := make([]byte, testImage16MB)

	offset, err := FindFirmwareEntryTable(smallImage)

	assert.EqualError(t, err, "No FirmwareTable found")
	assert.Equal(t, uint32(0), offset)
}

func TestFindFirmwareEntryTableByScanFETDefaultOffset(t *testing.T) {
	offset, err := FindFirmwareEntryTableByScan(mockFetImage())

	assert.Nil(t, err)
	assert.Equal(t, FETDefaultOffset, offset)
}

func TestFindFirmwareEntryTableFETDefaultOffset(t *testing.T) {
	offset, err := FindFirmwareEntryTable(mockFetImage())

	assert.Nil(t, err)
	assert.Equal(t, FETDefaultOffset, offset)
}

func TestParseFirmwareEntryTable(t *testing.T) {
	entryTable, err := ParseFirmwareEntryTable(fetBytes, 0)

	expectedFet := testFet
	expectedFet.Location = 0

	assert.Nil(t, err)
	assert.Equal(t, expectedFet, *entryTable)
}

func TestParseFirmwareEntryTableAtAddress(t *testing.T) {
	baseImage := mockFetImage()
	entryTable, err := ParseFirmwareEntryTable(baseImage, FETDefaultOffset)

	assert.Nil(t, err)
	assert.Equal(t, testFet, *entryTable)
}

func TestParseFirmwareEntryTableFailToSmall(t *testing.T) {
	entryTable, err := ParseFirmwareEntryTable(fetBytes, FETDefaultOffset)

	assert.EqualError(t, err, "Could not find FirmwareEntryTable Signature: Not AMD Table Header: Address out of bounds")
	assert.Nil(t, entryTable)
}

func TestParseFirmwareEntryTableWithWrongSignature(t *testing.T) {
	wrongBytes := []byte{
		0x12, 0x23, 0x34, 0x45,
		0x67, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	entryTable, err := ParseFirmwareEntryTable(wrongBytes, 0)

	assert.EqualError(t, err, "Could not find FirmwareEntryTable Signature: Not AMD Table Header")
	assert.Nil(t, entryTable)
}

func TestFirmwareEntryTable_Write(t *testing.T) {
	imageBytes := make([]byte, 500)
	expectedBytes := make([]byte, 500)

	err := testFet.Write(imageBytes, 100)

	copy(expectedBytes[100:], fetBytes)

	assert.Nil(t, err)
	assert.Equal(t, imageBytes, imageBytes)
}

func TestFirmwareEntryTable_WriteFailToSmall(t *testing.T) {
	imageBytes := make([]byte, FETDefaultOffset+5)

	err := testFet.Write(imageBytes, FETDefaultOffset)

	assert.EqualError(t, err, "BaseImage to small to insert FET")
}
