package amdfw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testRawRomBytes = []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7}
	testRawRom      = Rom{
		Type: XHCIRom,
		Raw:  testRawRomBytes,
	}
	testDirectoryRom = Rom{
		Type: PSPRom,
		Directories: []*Directory{
			&testPSPDirectory,
		},
	}
)

func mockRawRom(address uint32) []byte {
	baseImage := make([]byte, testImage16MB)

	copy(baseImage[address:], testRawRomBytes)
	return baseImage
}

func TestRom_WriteRaw(t *testing.T) {
	baseImage := make([]byte, testImage16MB)

	expectedImage := mockRawRom(testXHCRomBase - DefaultFlashMapping)

	err := testRawRom.Write(baseImage, &testFet, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestRom_WriteDir(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	testPSPDirectory.Write(expectedImage, testPSPDirBase)

	err := testDirectoryRom.Write(baseImage, &testFet, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestGetAddressFromTable(t *testing.T) {
	addressFromTable, err := GetAddressFromTable(PSPRom, &testFet)

	assert.Nil(t, err)
	assert.Equal(t, testPSPDirBase, addressFromTable)
}

func TestGetAddressFromTableUnknown(t *testing.T) {

	addressFromTable, err := GetAddressFromTable("FooBar", &testFet)

	assert.EqualError(t, err, "Cannot get Address: Unknown Type")
	assert.Equal(t, addressFromTable, uint32(0))
}
