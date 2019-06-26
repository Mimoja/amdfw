package main

import (
	"MimojaFirmwareToolkit/pkg/amdfw"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func main() {

	imageBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal("Could not read file: ", err)
	}

	image, err := amdfw.ParseImage(imageBytes)

	if err != nil {
		log.Println("Error while parse Image: ", err.Error())
	}

	renderFET(*image)

	for _, rom := range image.Roms {
		println()
		renderRom(*rom)
	}

}

func renderRom(rom amdfw.Rom) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)

	t.AppendHeader(table.Row{rom.Type})
	t.Render()

	for _, directory := range rom.Directories {

		t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleColoredBright)
		t.AppendHeader(table.Row{"Field", "Value"})

		checksum := "✓"
		if valid, should := directory.ValidateChecksum(); !valid {
			checksum = fmt.Sprintf("✕ (0x%08X)", should)
		}
		t.AppendRows([]table.Row{
			{"Magic", fmt.Sprintf("%s (0x%08X)", string(directory.Header.Cookie[:]), directory.Header.Cookie)},
			{"Checksum", fmt.Sprintf("0x%08X %s", directory.Header.Checksum, checksum)},
			{"Number of Entries", fmt.Sprintf("0x%08X", directory.Header.TotalEntries)},
			{"Reserved", fmt.Sprintf("0x%08X", directory.Header.Reserved)},
		})
		t.Render()

		t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleColoredBright)
		t.AppendHeader(table.Row{
			"Index",
			"Type",
			"Location",
			"Size",
			"Name",
			"ID",
			"SizeSigned",
			"Signed",
			"SigFingerprint",
			"Zipped",
			"FullSize",
			"Version",
			"SizePacked",
		})
		for entryID, entry := range directory.Entries {
			name := ""
			if entry.TypeInfo != nil {
				name = entry.TypeInfo.Name
			}

			nextRow := table.Row{
				fmt.Sprintf("0x%04X", entryID),
				fmt.Sprintf("0x%04X", entry.DirectoryEntry.Type),
				fmt.Sprintf("0x%08X", entry.DirectoryEntry.Location),
				fmt.Sprintf("0x%08X", entry.DirectoryEntry.Size),
				name,
			}
			if entry.Header != nil {
				nextRow = append(nextRow,
					fmt.Sprintf("0x%08X", entry.Header.ID),
					fmt.Sprintf("0x%08X", entry.Header.SizeSigned),
					fmt.Sprintf("0x%08X", entry.Header.IsSigned),
					fmt.Sprintf("0x%08X", entry.Header.SigFingerprint),
					fmt.Sprintf("%X", entry.Header.IsCompressed),
					fmt.Sprintf("0x%08X", entry.Header.FullSize),
					entry.Version,
					fmt.Sprintf("0x%08X", entry.Header.SizePacked),
				)
			}

			t.AppendRow(nextRow)
		}
		t.Render()

		for entryID, entry := range directory.Entries {

			if entry.Header != nil {
				t = table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				//t.SetStyle(table.StyleColoredBright)

				t.AppendHeader(table.Row{entryID, fmt.Sprintf("@ 0x%X", entry.DirectoryEntry.Location)})

				reflectVal := reflect.Indirect(reflect.ValueOf(entry.Header))
				for i := 0; i < reflectVal.Type().NumField(); i++ {
					fieldName := reflectVal.Type().Field(i).Name
					fieldValue := reflectVal.Field(i)
					t.AppendRow(table.Row{fieldName, fmt.Sprintf("0x%X", fieldValue)})
				}
				t.Render()
			}
		}

	}
}

func renderFET(image amdfw.Image) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)

	t.AppendHeader(table.Row{"Firmware Entry Table"})
	t.Render()

	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Field", "Value"})
	t.AppendRows([]table.Row{
		{"Signature", fmt.Sprintf("0x%08X", image.FET.Signature)},
		{"ImcRomBase", fmt.Sprintf("0x%08X", *image.FET.ImcRomBase)},
		{"GecRomBase", fmt.Sprintf("0x%08X", *image.FET.GecRomBase)},
		{"XHCRomBase", fmt.Sprintf("0x%08X", *image.FET.XHCRomBase)},
		{"PSPDirBase", fmt.Sprintf("0x%08X", *image.FET.PSPDirBase)},
		{"NewPSPDirBase", fmt.Sprintf("0x%08X", *image.FET.NewPSPDirBase)},
		{"BHDDirBase", fmt.Sprintf("0x%08X", *image.FET.BHDDirBase)},
		{"NewBHDDirBase", fmt.Sprintf("0x%08X", *image.FET.NewBHDDirBase)},
	})
	t.Render()

}
