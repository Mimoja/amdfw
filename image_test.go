package amdfw

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testFlashMapping = DefaultFlashMapping
	testImage        = Image{
		FET:          &testFet,
		FlashMapping: &testFlashMapping,
		Roms: []*Rom{
			{
				Type: PSPRom,
				Directories: []*Directory{
					&testPSPMiniDirectory,
				},
				Raw: nil,
			},
			&testRawRom,
		},
	}
)

func mockImage() []byte {
	baseImage := mockFetImage()
	copy(baseImage[testPSPDirBase-DefaultFlashMapping:], testPSPMiniDirectoryBytes)
	copy(baseImage[testXHCRomBase-DefaultFlashMapping:], testRawRomBytes)
	return baseImage
}

func TestImage_Write(t *testing.T) {
	baseImage := make([]byte, testImage16MB)

	baseImage, err := testImage.Write(baseImage)

	assert.Nil(t, err)

	for pos, byteVal := range mockImage() {
		if baseImage[pos] != byteVal {
			assert.Equal(t, byteVal, baseImage[pos], fmt.Sprintf("Bytes not Equal at 0x%08X", pos))
		}
	}

}
