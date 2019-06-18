package amdfw

import (
	"bytes"
	"fmt"
)

const DefaultFlashMapping = uint32(0xFF000000)

func GetFlashMapping(firmwareBytes []byte, fet *FirmwareEntryTable) (uint32, error) {

	type mappingMagic struct {
		addr  *uint32
		magic []string
	}

	for _, s := range []mappingMagic{{
		addr:  fet.PSPDirBase,
		magic: []string{PSPCOOCKIE, DUALPSPCOOCKIE},
	}, {
		addr:  fet.NewPSPDirBase,
		magic: []string{PSPCOOCKIE, DUALPSPCOOCKIE},
	}, {
		addr:  fet.BHDDirBase,
		magic: []string{BHDCOOCKIE},
	}, {
		addr:  fet.BHDDirBase,
		magic: []string{BHDCOOCKIE},
	},
	} {
		for _, m := range s.magic {

			if s.addr != nil && *s.addr != 0 {
				mapping, err := testMapping(firmwareBytes, *s.addr, m)
				if err == nil {
					return mapping, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("No valid mapping found!")
}

func testMapping(firmwareBytes []byte, address uint32, expected string) (uint32, error) {

	for _, mapping := range []uint32{
		DefaultFlashMapping + 0x000000, //16M
		DefaultFlashMapping + 0x800000, // 8M
		DefaultFlashMapping + 0xB00000, // 4M
		DefaultFlashMapping + 0xD00000, // 2M
		DefaultFlashMapping + 0xE00000, // 1M
		DefaultFlashMapping + 0xE80000, // 512K
	} {

		expectedBytes := []byte(expected)
		testAddr := address - mapping
		if int(testAddr) > len(firmwareBytes) {
			continue
		}

		if bytes.Equal(firmwareBytes[testAddr:testAddr+4], expectedBytes) {
			return mapping, nil
		}
	}
	return 0, fmt.Errorf("No Default Mapping fits")
}
