package amdfw

import (
	"fmt"
)

type (
	Image struct {
		FET          *FirmwareEntryTable
		FlashMapping *uint32
		Roms         []*Rom
	}
)

func ParseImage(firmwareBytes []byte) (*Image, error) {
	image := Image{}

	fetOffset, err := FindFirmwareEntryTable(firmwareBytes)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Image: %v", err)
	}

	fet, err := ParseFirmwareEntryTable(firmwareBytes, fetOffset)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Image: %v", err)
	}
	image.FET = fet

	mapping, err := GetFlashMapping(firmwareBytes, fet)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Image: %v", err)
	}
	image.FlashMapping = &mapping

	roms, errs := ParseRoms(firmwareBytes, fet, mapping)
	if len(errs) != 0 {
		err = fmt.Errorf("Errors parsing images %v", errs)
	} else {
		err = nil
	}

	image.Roms = roms
	return &image, err
}

func (image *Image) Write(baseImage []byte) ([]byte, error) {
	var err error
	if err = image.FET.Write(baseImage, image.FET.Location); err != nil {
		return nil, err
	}

	if image.FlashMapping != nil {
		for _, rom := range image.Roms {
			if err = rom.Write(baseImage, image.FET, *image.FlashMapping); err != nil {
				return nil, err
			}

		}
	}
	return baseImage, nil
}
