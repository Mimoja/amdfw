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

func TestRom_WriteRawToSmall(t *testing.T) {
	baseImage := make([]byte, testXHCRomBase-DefaultFlashMapping)
	untouched := make([]byte, testXHCRomBase-DefaultFlashMapping)

	err := testRawRom.Write(baseImage, &testFet, DefaultFlashMapping)

	assert.EqualError(t, err, "Cannot write Rom: Invalid address in FET")
	assert.Equal(t, baseImage, untouched)
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

func TestRom_WriteDirFailed(t *testing.T) {
	baseImage := make([]byte, testPSPDirBase-5-DefaultFlashMapping)
	untouchedImage := make([]byte, testPSPDirBase-5-DefaultFlashMapping)

	err := testDirectoryRom.Write(baseImage, &testFet, DefaultFlashMapping)

	assert.EqualError(t, err, "Cannot Write Rom: BaseImage to small to insert Directory")
	assert.Equal(t, untouchedImage, baseImage)
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

func TestParseRomsRecursiv(t *testing.T) {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], test2PSPDirectoryBytes)

	for _, entry := range test2PSPDirectory.Entries {

		copy(imageBytes[entry.DirectoryEntry.Location&^DefaultFlashMapping:], testPSPMiniDirectoryBytes)
	}

	directories, err := recursiveDirectories(imageBytes, &test2PSPDirectory, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, len(directories), 4)

	for _, dir := range directories {
		assert.Equal(t, testPSPMiniDirectory.Header, dir.Header)
		for i, _ := range dir.Entries {
			assert.Equal(t, testPSPMiniDirectory.Entries[i].DirectoryEntry, dir.Entries[i].DirectoryEntry)
		}
	}
}
