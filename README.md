# Classic Converter

A tool to convert Minecraft pre-classic and classic worlds to a Minecraft indev world or a schematic file

## How To Build

**Prerequisites**

You will need the following in order to build the converter:

* A computer
* Go 1.20 or later

**Building**

1. Open a terminal
2. Clone the repo with: `git clone https://github.com/BJTMastermind/Classic-Converter.git`
3. cd into the Classic-Converter folder
4. Run: `go build`

You have now built the Classic Converter!

## How To Use

1. Open a terminal
2. Run `Classic-Converter -i /path/to/classic_world.mine -f "indev_level"` (Only `.dat` and `.mine` files are accepted)
3. You now have a converted classic world in the specified format. (File keeps the same name as original)<br>Outputed file would be: `classic_world.mclevel`.

## Language(s) Used

* Go 1.20
