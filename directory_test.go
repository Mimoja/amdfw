package amdfw

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (

	/**
	  00100000  32 50 53 50 ad 33 60 7f  04 00 00 00 00 00 00 00  |2PSP.3`.........|
	  00100010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
	  00100020  00 00 00 00 00 05 0b bc  00 90 2e 00 00 00 00 00  |................|
	  00100030  00 00 00 00 00 00 0a bc  00 d0 1d ff 00 00 00 00  |................|
	  00100040  00 00 00 00 00 01 0a bc  00 d0 1d ff 00 00 00 00  |................|
	  00100050  00 00 00 00 00 00 09 bc  00 10 11 ff 00 00 00 00  |................|
	  00100060  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  |................|
	*/

	test2PSPDirectoryBytes = []byte{
		0x32, 0x50, 0x53, 0x50, 0xad, 0x33, 0x60, 0x7f, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x0b, 0xbc, 0x00, 0x90, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0xbc, 0x00, 0xd0, 0x1d, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x0a, 0xbc, 0x00, 0xd0, 0x1d, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0xbc, 0x00, 0x10, 0x11, 0xff, 0x00, 0x00, 0x00, 0x00,
	}

	test2PSPDirectory = Directory{
		Header: DirectoryHeader{
			Cookie:       [4]uint8{0x32, 0x50, 0x53, 0x50},
			Checksum:     0x7f6033ad,
			TotalEntries: 0x4,
			Reserved:     0x0,
		},
		Location: testNewPSPDirBase,
		Entries: []Entry{
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0xbc0b0500, Location: 0x002e9000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0xbc0a0000, Location: 0xff1dd000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0xbc0a0100, Location: 0xff1dd000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0xbc090000, Location: 0xff111000, Reserved: 0x0}},
		},
	}
	/*
	   000c0000  24 50 53 50 70 2e 88 ea  14 00 00 00 01 0e 18 80  |$PSPp...........|
	   000c0010  00 00 00 00 40 02 00 00  00 10 0c ff 00 00 00 00  |....@...........|
	   000c0020  01 00 00 00 00 80 00 00  00 10 18 ff 00 00 00 00  |................|
	   000c0030  08 00 00 00 00 40 01 00  00 90 18 ff 00 00 00 00  |.....@..........|
	   000c0040  03 00 00 00 00 60 00 00  00 20 0c ff 00 00 00 00  |.....`... ......|
	   000c0050  05 00 00 00 40 03 00 00  00 80 0c ff 00 00 00 00  |....@...........|
	   000c0060  06 00 00 00 00 10 00 00  00 f0 ff ff 00 00 00 00  |................|
	   000c0070  02 00 00 00 00 e0 01 00  00 d0 19 ff 00 00 00 00  |................|
	   000c0080  04 00 00 00 00 00 01 00  00 00 09 ff 00 00 00 00  |................|
	   000c0090  08 01 00 00 00 40 01 00  00 b0 1b ff 00 00 00 00  |.....@..........|
	   000c00a0  09 00 00 00 40 03 00 00  00 90 0c ff 00 00 00 00  |....@...........|
	   000c00b0  0b 00 00 00 ff ff ff ff  01 00 00 00 00 00 00 00  |................|
	   000c00c0  0c 00 00 00 00 a0 01 00  00 f0 1c ff 00 00 00 00  |................|
	   000c00d0  0d 00 00 00 40 03 00 00  00 a0 0c ff 00 00 00 00  |....@...........|
	   000c00e0  10 00 00 00 00 80 00 00  00 90 1e ff 00 00 00 00  |................|
	   000c00f0  12 00 00 00 00 b0 00 00  00 10 1f ff 00 00 00 00  |................|
	   000c0100  14 00 00 00 00 00 02 00  00 c0 1f ff 00 00 00 00  |................|
	   000c0110  12 01 00 00 00 b0 00 00  00 c0 21 ff 00 00 00 00  |..........!.....|
	   000c0120  5f 00 00 00 00 10 00 00  00 70 22 ff 00 00 00 00  |_........p".....|
	   000c0130  5f 01 00 00 00 10 00 00  00 80 22 ff 00 00 00 00  |_.........".....|
	   000c0140  1a 00 00 00 00 30 00 00  00 90 22 ff 00 00 00 00  |.....0....".....|
	   000c0150  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  |................|
	*/

	testPSPDirectoryBytes = []byte{
		0x24, 0x50, 0x53, 0x50, 0x70, 0x2e, 0x88, 0xea, 0x14, 0x00, 0x00, 0x00, 0x01, 0x0e, 0x18, 0x80,
		0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0x00, 0x00, 0x00, 0x10, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x10, 0x18, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x08, 0x00, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x90, 0x18, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 0x20, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x05, 0x00, 0x00, 0x00, 0x40, 0x03, 0x00, 0x00, 0x00, 0x80, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x06, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x01, 0x00, 0x00, 0xd0, 0x19, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x09, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x08, 0x01, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0xb0, 0x1b, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x09, 0x00, 0x00, 0x00, 0x40, 0x03, 0x00, 0x00, 0x00, 0x90, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x0b, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x01, 0x00, 0x00, 0xf0, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x0d, 0x00, 0x00, 0x00, 0x40, 0x03, 0x00, 0x00, 0x00, 0xa0, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x90, 0x1e, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x12, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x10, 0x1f, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0xc0, 0x1f, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x12, 0x01, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0xc0, 0x21, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x5f, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x70, 0x22, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x5f, 0x01, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x80, 0x22, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x1a, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x90, 0x22, 0xff, 0x00, 0x00, 0x00, 0x00,
	}

	testPSPDirectory = Directory{
		Header: DirectoryHeader{
			Cookie:       [4]byte{0x24, 0x50, 0x53, 0x50},
			Checksum:     0xea882e70,
			TotalEntries: 0x14,
			Reserved:     0x80180e01,
		},
		Location: testPSPDirBase - DefaultFlashMapping,
		Entries: []Entry{
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0x240, Location: 0xff0c1000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x1, Size: 0x8000, Location: 0xff181000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x8, Size: 0x14000, Location: 0xff189000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x3, Size: 0x6000, Location: 0xff0c2000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x5, Size: 0x340, Location: 0xff0c8000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x6, Size: 0x1000, Location: 0xfffff000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x2, Size: 0x1e000, Location: 0xff19d000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x4, Size: 0x10000, Location: 0xff090000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x108, Size: 0x14000, Location: 0xff1bb000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x9, Size: 0x340, Location: 0xff0c9000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0xb, Size: 0xffffffff, Location: 0x1, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0xc, Size: 0x1a000, Location: 0xff1cf000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0xd, Size: 0x340, Location: 0xff0ca000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x10, Size: 0x8000, Location: 0xff1e9000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x12, Size: 0xb000, Location: 0xff1f1000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x14, Size: 0x20000, Location: 0xff1fc000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x112, Size: 0xb000, Location: 0xff21c000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x5f, Size: 0x1000, Location: 0xff227000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x15f, Size: 0x1000, Location: 0xff228000, Reserved: 0x0}},
			{DirectoryEntry: DirectoryEntry{Type: 0x1a, Size: 0x3000, Location: 0xff229000, Reserved: 0x0}}},
	}

	testPSPDirectoryHeaderBytes = testPSPDirectoryBytes[:16]
	testPSPDirectoryEntryBytes  = testPSPDirectoryBytes[16:32]

	testPSPMiniDirectoryBytes = []byte{
		0x24, 0x50, 0x53, 0x50, 0x70, 0x2e, 0x88, 0xea, 0x01, 0x00, 0x00, 0x00, 0x01, 0x0e, 0x18, 0x80,
		0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0x00, 0x00, 0x00, 0x10, 0x0c, 0xff, 0x00, 0x00, 0x00, 0x00,
	}

	testPSPMiniDirectory = Directory{
		Header: DirectoryHeader{
			Cookie:       [4]byte{0x24, 0x50, 0x53, 0x50},
			Checksum:     0xea882e70, // Not checked and wrong
			TotalEntries: 0x1,
			Reserved:     0x80180e01,
		},
		Location: testPSPDirBase - DefaultFlashMapping,
		Entries: []Entry{
			{DirectoryEntry: DirectoryEntry{Type: 0x0, Size: 0x240, Location: 0xff0c1000, Reserved: 0x0}},
		},
	}

	testBHDDirectoryBytes = []byte{
		0x24, 0x42, 0x48, 0x44, 0x52, 0xf8, 0x63, 0x9e, 0x0b, 0x00, 0x00, 0x00, 0x1c, 0x04, 0x00, 0x00,
		0x60, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x60, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00,
		0x00, 0x40, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x68, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x60, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x68, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00,
		0x00, 0x80, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x20, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x62, 0x00, 0x03, 0x00, 0x00, 0x00, 0x20, 0x00,
		0x00, 0x00, 0xe0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x09, 0x00, 0x00, 0x00, 0x00,
		0x64, 0x00, 0x10, 0x00, 0x40, 0x3c, 0x00, 0x00, 0x00, 0xa0, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x65, 0x00, 0x10, 0x00, 0x30, 0x03, 0x00, 0x00,
		0x00, 0xdd, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x64, 0x00, 0x40, 0x00, 0x10, 0x46, 0x00, 0x00, 0x00, 0xe1, 0x1c, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x65, 0x00, 0x40, 0x00, 0x20, 0x03, 0x00, 0x00,
		0x00, 0x28, 0x1d, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x70, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x10, 0x64, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}

	testBHDDirectoryHeaderBytes = testBHDDirectoryBytes[0:16]
	testBHDDirectoryEntryBytes  = testBHDDirectoryBytes[16 : 16+4*6]

	testBHDDirectoryHeader     = DirectoryHeader{Cookie: [4]uint8{0x24, 0x42, 0x48, 0x44}, Checksum: 0x9e63f852, TotalEntries: 0xb, Reserved: 0x41c}
	testBHDDirEntryUnknown0xFF = uint64(0xffffffffffffffff)
	testBHDDirEntryKnown0xA2   = uint64(0xa200000)
	testBHDDirEntryKnown0x9E   = uint64(0x9e00000)

	testBHDDirectory = Directory{
		Header: testBHDDirectoryHeader,
		Entries: []Entry{
			{DirectoryEntry: DirectoryEntry{Type: 0x60, Size: 0x2000, Location: 0xff1c2000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x200060, Size: 0x2000, Location: 0xff1c4000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x68, Size: 0x2000, Location: 0xff1c6000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x200068, Size: 0x2000, Location: 0xff1c8000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x61, Size: 0x0, Location: 0x0, Reserved: 0x0, Unknown: &testBHDDirEntryKnown0xA2}},
			{DirectoryEntry: DirectoryEntry{Type: 0x30062, Size: 0x200000, Location: 0xffe00000, Reserved: 0x0, Unknown: &testBHDDirEntryKnown0x9E}},
			{DirectoryEntry: DirectoryEntry{Type: 0x100064, Size: 0x3c40, Location: 0xff1ca000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x100065, Size: 0x330, Location: 0xff1cdd00, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x400064, Size: 0x4610, Location: 0xff1ce100, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x400065, Size: 0x320, Location: 0xff1d2800, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
			{DirectoryEntry: DirectoryEntry{Type: 0x70, Size: 0x400, Location: 0xff641000, Reserved: 0x0, Unknown: &testBHDDirEntryUnknown0xFF}},
		},
		Location: testBHDDirBase - DefaultFlashMapping,
	}
)

func TestParseDirectoryHeaderToSmall(t *testing.T) {
	imageToSmall := make([]byte, FETDefaultOffset/2)

	dir, err := ParseDirectory(imageToSmall, DefaultFlashMapping+FETDefaultOffset, DefaultFlashMapping)

	assert.EqualError(t, err, "Firmwarebytes not long enough for reading directory header..")
	assert.Nil(t, dir)

}

func TestParseDirectoryHeaderCookieWrong(t *testing.T) {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], "$INV")

	dir, err := ParseDirectory(imageBytes, uint32(testPSPDirBase), DefaultFlashMapping)

	assert.EqualError(t, err, "No Valid Cookie at start of directory: [36 73 78 86]")
	assert.Nil(t, dir)
}

func TestParseDirectoryHeaderToManyEntries(t *testing.T) {
	overSizedDirectory := testPSPDirectory
	overSizedDirectory.Header.TotalEntries = 1

	imgSize := testPSPDirBase - DefaultFlashMapping + 16
	imageBytes := make([]byte, imgSize)
	overSizedDirectory.Header.Write(imageBytes, testPSPDirBase-DefaultFlashMapping)

	// Read with mapping
	dir, err := ParseDirectory(imageBytes, testPSPDirBase, DefaultFlashMapping)

	assert.EqualError(t, err, "Could not read directory entries from image: Too many entries(1)!")
	assert.Nil(t, dir)

}

func TestParseDirectoryPass(t *testing.T) {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], testPSPDirectoryBytes)

	// Read with mapping
	dir, err := ParseDirectory(imageBytes, testPSPDirBase, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.NotNil(t, dir)

	// And without
	dir2, err2 := ParseDirectory(imageBytes, testPSPDirBase, DefaultFlashMapping)

	assert.Nil(t, err2)
	assert.NotNil(t, dir2)
	assert.Equal(t, *dir, *dir2)

	assert.Equal(t, testPSPDirectory.Header, dir.Header)
	assert.Equal(t, len(testPSPDirectory.Entries), len(dir.Entries))

	for i := 0; i < len(dir.Entries); i++ {
		assert.Equal(t, testPSPDirectory.Entries[i].DirectoryEntry, dir.Entries[i].DirectoryEntry)
	}
}

func TestParse2PSPDirectory(t *testing.T) {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[testPSPDirBase-DefaultFlashMapping:], test2PSPDirectoryBytes)

	// Read with mapping
	dir, err := ParseDirectory(imageBytes, testPSPDirBase, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.NotNil(t, dir)

	assert.Equal(t, test2PSPDirectory.Header, dir.Header)
	assert.Equal(t, len(test2PSPDirectory.Entries), len(dir.Entries))

	for i := 0; i < len(dir.Entries); i++ {
		assert.Equal(t, test2PSPDirectory.Entries[i].DirectoryEntry, dir.Entries[i].DirectoryEntry)
	}
}

func TestParseBHDDirectory(t *testing.T) {
	imageBytes := make([]byte, testImage16MB)
	copy(imageBytes[testBHDDirBase-DefaultFlashMapping:], testBHDDirectoryBytes)

	// Read with mapping
	dir, err := ParseDirectory(imageBytes, testBHDDirBase, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.NotNil(t, dir)

	assert.Equal(t, testBHDDirectory.Header, dir.Header)
	assert.Equal(t, len(testBHDDirectory.Entries), len(dir.Entries))
	assert.Equal(t, testBHDDirBase-DefaultFlashMapping, dir.Location)

	for i := 0; i < len(dir.Entries); i++ {
		assert.Equal(t, *testBHDDirectory.Entries[i].DirectoryEntry.Unknown, *dir.Entries[i].DirectoryEntry.Unknown, fmt.Sprint("Entries padding does not match: ", i))

		// Whipe unknown Pointer
		dir.Entries[i].DirectoryEntry.Unknown = testBHDDirectory.Entries[i].DirectoryEntry.Unknown
		assert.Equal(t, testBHDDirectory.Entries[i].DirectoryEntry, dir.Entries[i].DirectoryEntry)
	}
}

func TestDirectory_ValidateChecksum(t *testing.T) {
	valid, actual := testPSPDirectory.ValidateChecksum()

	assert.Equal(t, true, valid)
	assert.Equal(t, testPSPDirectory.Header.Checksum, actual)
}

func TestDirectory_ValidateChecksum2PSP(t *testing.T) {
	valid, actual := test2PSPDirectory.ValidateChecksum()

	assert.Equal(t, true, valid)
	assert.Equal(t, test2PSPDirectory.Header.Checksum, actual)
}

func TestDirectory_ValidateChecksumBHD(t *testing.T) {
	valid, actual := testBHDDirectory.ValidateChecksum()

	assert.Equal(t, true, valid)
	assert.Equal(t, testBHDDirectory.Header.Checksum, actual)
}

func TestDirectory_Write(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	copy(expectedImage[testPSPDirectory.Location:], testPSPDirectoryBytes)

	err := testPSPDirectory.Write(baseImage, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestDirectoryHeader_Write(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	copy(expectedImage[testPSPDirectory.Location:], testPSPDirectoryHeaderBytes)

	err := testPSPDirectory.Header.Write(baseImage, testPSPDirectory.Location)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestDirectoryEntry_Write(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	copy(expectedImage[testPSPDirectory.Location+16:], testPSPDirectoryEntryBytes)

	err := testPSPDirectory.Entries[0].DirectoryEntry.Write(baseImage, testPSPDirectory.Location+16)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestDirectoryEntry_WriteBHD(t *testing.T) {
	baseImage := make([]byte, testBHDDirectory.Location+0x100)
	expectedImage := make([]byte, testBHDDirectory.Location+0x100)
	copy(expectedImage[testBHDDirectory.Location+16:], testBHDDirectoryEntryBytes)

	err := testBHDDirectory.Entries[0].DirectoryEntry.Write(baseImage, testBHDDirectory.Location+16)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestDirectory_WriteToSmall(t *testing.T) {
	baseImage := make([]byte, FETDefaultOffset+5)

	err := testPSPDirectory.Write(baseImage, DefaultFlashMapping)

	assert.EqualError(t, err, "BaseImage to small to insert Directory")
}

func TestDirectory_WriteBHD(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	copy(expectedImage[testBHDDirectory.Location:], testBHDDirectoryBytes)

	err := testBHDDirectory.Write(baseImage, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}

func TestDirectory_Write2PSP(t *testing.T) {
	baseImage := make([]byte, testImage16MB)
	expectedImage := make([]byte, testImage16MB)
	copy(expectedImage[test2PSPDirectory.Location:], test2PSPDirectoryBytes)

	err := test2PSPDirectory.Write(baseImage, DefaultFlashMapping)

	assert.Nil(t, err)
	assert.Equal(t, expectedImage, baseImage)
}
