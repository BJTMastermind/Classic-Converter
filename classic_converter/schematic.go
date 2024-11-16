package classic_converter

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/BJTMastermind/go-nbt"
)

type Schematic struct {
    Width int16
    Height int16
    Length int16
    Blocks []int8
    Data []int8
    Entities []nbt.Compound
    TileEntities []nbt.Compound
}

func (schematic *Schematic) InitWithDefaults() *Schematic {
    schematic.Width = 256
    schematic.Height = 64
    schematic.Length = 256
    schematic.Blocks = make([]int8, 256*256*64)
    schematic.Data = make([]int8, 256*256*64)
    schematic.Entities = []nbt.Compound{}
    schematic.TileEntities = []nbt.Compound{}

    return schematic
}

func (schematic *Schematic) WriteToFile(filename string) {
    tag := nbt.NewCompoundTag("Schematic", map[string]nbt.Tag{
        "Width": &nbt.Short{
            Value: schematic.Width,
        },
        "Height": &nbt.Short{
            Value: schematic.Height,
        },
        "Length": &nbt.Short{
            Value: schematic.Length,
        },
        "Materials": &nbt.String{
            Value: "Alpha",
        },
        "Blocks": &nbt.ByteArray{
            Value: schematic.Blocks,
        },
        "Data": &nbt.ByteArray{
            Value: schematic.Data,
        },
        "Entities": &nbt.List{
            Value: compoundArrayToTagArray(schematic.Entities),
            ListType: ternary[int8](len(schematic.Entities) == 0, nbt.IDTagEnd, nbt.IDTagCompound),
        },
        "TileEntities": &nbt.List{
            ListType: nbt.IDTagEnd,
        },
    })

    stream := nbt.NewStream(nbt.BigEndian)

    err := stream.WriteTag(tag)
    if err != nil {
        panic(err)
    }

    data, err := nbt.Compress(stream, nbt.CompressGZip, nbt.DefaultCompressionLevel)
    if err != nil {
        panic(err)
    }

    ioutil.WriteFile(filename, data, os.ModePerm)

    fmt.Printf("Generated %s\n", filename)
}
