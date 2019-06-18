package amdfw

import (
	"bytes"
	"fmt"
)

const (
	IMCRom    RomType = "IMC"
	GECRom    RomType = "GEC"
	XHCIRom   RomType = "XHCI"
	PSPRom    RomType = "PSP"
	NewPSPRom RomType = "NEWPSP"
	BHDRom    RomType = "BHD"
	NewBHDRom RomType = "NEWBDH"
)

type (
	RomType string

	Rom struct {
		Type        RomType
		Directories []*Directory
		Raw         []byte
	}
)

func ParseRoms(firmwareBytes []byte, table *FirmwareEntryTable, flashMapping uint32) ([]*Rom, []error) {
	var roms []*Rom

	var errors []error
	// Parsing ROMs without length is annoying. Not going to implement soon...
	//TODO IMC
	//TODO GEC
	//TODO XHCI

	// PSP
	rom, err := ParsePSPRom(firmwareBytes, table, flashMapping)
	if err != nil {
		errors = append(errors, fmt.Errorf("Could not parse psp rom: %v", err))
	}
	if rom != nil {
		roms = append(roms, rom)
	}

	// newPSP
	rom, err = ParseNewPSPRom(firmwareBytes, table, flashMapping)
	if err != nil {
		errors = append(errors, fmt.Errorf("Could not parse newpsp rom: %v", err))
	}
	if rom != nil {
		roms = append(roms, rom)
	}

	// BHD
	rom, err = ParseBHDRom(firmwareBytes, table, flashMapping)
	if err != nil {
		errors = append(errors, fmt.Errorf("Could not parse bhd rom: %v", err))
	}
	if rom != nil {
		roms = append(roms, rom)
	}

	// newBHD
	rom, err = ParseNewBHDRom(firmwareBytes, table, flashMapping)
	if err != nil {
		errors = append(errors, fmt.Errorf("Could not parse new bhd rom: %v", err))
	}
	if rom != nil {
		roms = append(roms, rom)
	}

	return roms, errors
}

func ParsePSPRom(firmwareBytes []byte, table *FirmwareEntryTable, flashMapping uint32) (*Rom, error) {
	return parseDirectoryRom(firmwareBytes, table.PSPDirBase, flashMapping, PSPRom)
}

func ParseNewPSPRom(firmwareBytes []byte, table *FirmwareEntryTable, flashMapping uint32) (*Rom, error) {
	return parseDirectoryRom(firmwareBytes, table.NewPSPDirBase, flashMapping, NewPSPRom)
}

func ParseBHDRom(firmwareBytes []byte, table *FirmwareEntryTable, flashMapping uint32) (*Rom, error) {
	return parseDirectoryRom(firmwareBytes, table.BHDDirBase, flashMapping, BHDRom)
}

func ParseNewBHDRom(firmwareBytes []byte, table *FirmwareEntryTable, flashMapping uint32) (*Rom, error) {
	return parseDirectoryRom(firmwareBytes, table.NewBHDDirBase, flashMapping, NewBHDRom)
}

func parseDirectoryRom(firmwareBytes []byte, address *uint32, flashMapping uint32, romType RomType) (*Rom, error) {
	rom := Rom{
		Type: romType,
	}

	if address == nil {
		return nil, fmt.Errorf("No %s offset available", romType)
	}

	directory, err := ParseDirectory(firmwareBytes, *address, flashMapping)

	if err != nil {
		return nil, fmt.Errorf("Could not read %s Rom: %v", romType, err)
	}

	rom.Directories = append(rom.Directories, directory)
	others, err := recursiveDirectories(firmwareBytes, directory, flashMapping)

	rom.Directories = append(rom.Directories, others...)
	return &rom, err

}

func recursiveDirectories(firmwareBytes []byte, directory *Directory, flashMapping uint32) ([]*Directory, error) {
	//TODO Test
	var directories []*Directory
	for _, entry := range directory.Entries {
		if entry.DirectoryEntry.Type == 0x40 ||
			entry.DirectoryEntry.Type == 0x70 ||
			bytes.Equal(directory.Header.Cookie[:], []byte(DUALPSPCOOCKIE)) {

			newDirectory, err := ParseDirectory(firmwareBytes, entry.DirectoryEntry.Location, flashMapping)

			if err != nil {
				return directories, fmt.Errorf("Could not read Directory: %v", err)
			}

			directories = append(directories, newDirectory)
			others, err := recursiveDirectories(firmwareBytes, newDirectory, flashMapping)
			if err != nil {
				return directories, fmt.Errorf("Could not read Directory: %v", err)
			}
			directories = append(directories, others...)
		}
	}
	return directories, nil
}

func GetAddressFromTable(romType RomType, table *FirmwareEntryTable) (uint32, error) {
	switch romType {
	case PSPRom:
		return *table.PSPDirBase, nil
	case NewPSPRom:
		return *table.NewPSPDirBase, nil
	case BHDRom:
		return *table.BHDDirBase, nil
	case NewBHDRom:
		return *table.NewBHDDirBase, nil
	case GECRom:
		return *table.GecRomBase, nil
	case IMCRom:
		return *table.ImcRomBase, nil
	case XHCIRom:
		return *table.XHCRomBase, nil
	default:
		return 0, fmt.Errorf("Cannot get Address: Unknown Type")
	}

}

func (rom Rom) Write(baseImage []byte, table *FirmwareEntryTable, flashMapping uint32) error {
	var err error
	if rom.Raw != nil {
		address, err := GetAddressFromTable(rom.Type, table)

		if err != nil {
			return fmt.Errorf("Cannot Write: Unknown Type")
		}

		address = address &^ flashMapping

		//TODO Test
		if int(address) > len(baseImage) {
			return fmt.Errorf("Cannot write Rom: Invalid address in FET")
		}

		copy(baseImage[address:], rom.Raw)
	} else if rom.Directories != nil {
		//TODO Test
		for _, directory := range rom.Directories {
			err = directory.Write(baseImage, flashMapping)
			if err != nil {
				//TODO Test
				return fmt.Errorf("Cannot Write Rom: %v", err)
			}
		}
	}
	return nil
}
