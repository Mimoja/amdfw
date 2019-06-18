package amdfw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockPspHeaderImage() []byte {

	imageBytes := mockFetImage()

	// Copy $PSP Header
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], PSPCOOCKIE)
	return imageBytes
}

func TestGetFlashMappingValid(t *testing.T) {

	mapping, err := GetFlashMapping(mockPspHeaderImage(), &testFet)

	assert.Nil(t, err)
	assert.Equal(t, DefaultFlashMapping, mapping)
}

func TestGetFlashMappingInvalid(t *testing.T) {
	empty := make([]byte, testImage16MB)
	mapping, err := GetFlashMapping(empty, &testFet)

	assert.EqualError(t, err, "No valid mapping found!")
	assert.Equal(t, uint32(0), mapping)
}
