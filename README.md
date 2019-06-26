# amdfw

Golang library for reading and writing AMD firmware components

Credit goes to @cweling for his [psptool](https://github.com/cwerling/psptool)

## amddump
cmd/amddump is a small tool, that dumps all informations known to this library on a specfic image.


```
amddump ryzeimage.rom

```

## Current Limitations
- Always assumes valid FirmwareEntryTable. 
  - Some AM1 CPUs are not using it.
  - Older FETs might be parsed wrong
- Non Directory-Based Firmware (IMC, GEC, XHCI) cannot be extracted
- All Offsets are treated as absolute. Partial Images often can't be read.

## Usage

See cmd/amddump.go for read-only example code.

```golang

func main() {

	imageBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal("Could not read file: ", err)
	}

	image, err := amdfw.ParseImage(imageBytes)

	if err != nil {
		log.Println("Error while parse Image: ", err.Error())
	}

	targetAddress := uint32(0x1C1000)
	image.FET.ImcRomBase = &targetAddress

	image.FET.Write(imageBytes, image.FET.Location)

	ioutil.WriteFile(os.Args[1], imageBytes, 666)
}

```
