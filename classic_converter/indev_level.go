package classic_converter

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "os"
    "time"

    "github.com/BJTMastermind/go-nbt"
)

type IndevLevel struct {
    // Root - About
    CreatedOn int64
    Name string
    Author string
    // Root - Environment
    TimeOfDay int16
    SkyBrightness int8
    SkyColor int32
    FogColor int32
    CloudColor int32
    CloudHeight int16
    SurroundingGroundType int8
    SurroundingGroundHeight int16
    SurroundingWaterType int8
    SurroundingWaterHeight int16
    // Root - Map
    Width int16
    Length int16
    Height int16
    Spawn [3]int16
    Blocks []int8
    Data []int8
    // Root
    Entities []nbt.Compound
    TileEntities []nbt.Compound
}

func (indev_level *IndevLevel) InitWithDefaults() *IndevLevel {
    indev_level.CreatedOn = time.Now().Unix()
    indev_level.Name = "A Nice World"
    indev_level.Author = ""
    indev_level.TimeOfDay = 0
    indev_level.SkyBrightness = 15
    indev_level.SkyColor = 10079487
    indev_level.FogColor = 16777215
    indev_level.CloudColor = 16777215
    indev_level.CloudHeight = 66
    indev_level.SurroundingGroundType = 2 //7
    indev_level.SurroundingGroundHeight = 23 //31
    indev_level.SurroundingWaterType = 8
    indev_level.SurroundingWaterHeight = 32 //2
    indev_level.Width = 256
    indev_level.Length = 256
    indev_level.Height = 64
    indev_level.Spawn = [3]int16{int16(rand.Intn(256)), int16(rand.Intn(64 - 1) + 1), int16(rand.Intn(256))}
    indev_level.Blocks = make([]int8, 256*256*64)
    indev_level.Data = make([]int8, 256*256*64)
    indev_level.Entities = []nbt.Compound{}
    indev_level.TileEntities = []nbt.Compound{}

    return indev_level
}

func (indev_level *IndevLevel) FindSpawn() {
    i := 0

    var x int32
    var y int32
    var z int32
    for ok := true; ok; ok = (y <= int32(indev_level.Height) / 2) { // do while loop
        i++
        x = rand.Int31n(int32(indev_level.Width) / 2) + int32(indev_level.Width) / 4
        z = rand.Int31n(int32(indev_level.Length) / 2) + int32(indev_level.Length) / 4
        y = indev_level.getHightestTile(x, z) + 1
        if i == 10000 {
            indev_level.Spawn[0] = int16(x)
            indev_level.Spawn[1] = -100
            indev_level.Spawn[2] = int16(z)
            return
        }
    }

    indev_level.Spawn[0] = int16(x)
    indev_level.Spawn[1] = int16(y)
    indev_level.Spawn[2] = int16(z)
}

func (indev_level *IndevLevel) WriteToFile(filename string) {
    root := nbt.NewCompoundTag("MinecraftLevel", map[string]nbt.Tag{
        "About": &nbt.Compound{
            Value: map[string]nbt.Tag{
                "CreatedOn": &nbt.Long{
                    Value: indev_level.CreatedOn,
                },
                "Name": &nbt.String{
                    Value: indev_level.Name,
                },
                "Author": &nbt.String{
                    Value: indev_level.Author,
                },
            },
        },
        "Environment": &nbt.Compound{
            Value: map[string]nbt.Tag{
                "TimeOfDay": &nbt.Short{
                    Value: indev_level.TimeOfDay,
                },
                "SkyBrightness": &nbt.Byte{
                    Value: indev_level.SkyBrightness,
                },
                "SkyColor": &nbt.Int{
                    Value: indev_level.SkyColor,
                },
                "FogColor": &nbt.Int{
                    Value: indev_level.FogColor,
                },
                "CloudColor": &nbt.Int{
                    Value: indev_level.CloudColor,
                },
                "CloudHeight": &nbt.Short{
                    Value: indev_level.CloudHeight,
                },
                "SurroundingGroundType": &nbt.Byte{
                    Value: indev_level.SurroundingGroundType,
                },
                "SurroundingGroundHeight": &nbt.Short{
                    Value: indev_level.SurroundingGroundHeight,
                },
                "SurroundingWaterType": &nbt.Byte{
                    Value: indev_level.SurroundingWaterType,
                },
                "SurroundingWaterHeight": &nbt.Short{
                    Value: indev_level.SurroundingWaterHeight,
                },
            },
        },
        "Map": &nbt.Compound{
            Value: map[string]nbt.Tag{
                "Width": &nbt.Short{
                    Value: indev_level.Width,
                },
                "Length": &nbt.Short{
                    Value: indev_level.Length,
                },
                "Height": &nbt.Short{
                    Value: indev_level.Height,
                },
                "Spawn": &nbt.List{
                    Value: []nbt.Tag{
                        &nbt.Short{
                            Value: indev_level.Spawn[0],
                        },
                        &nbt.Short{
                            Value: indev_level.Spawn[1],
                        },
                        &nbt.Short{
                            Value: indev_level.Spawn[2],
                        },
                    },
                    ListType: nbt.IDTagShort,
                },
                "Blocks": &nbt.ByteArray{
                    Value: indev_level.Blocks,
                },
                "Data": &nbt.ByteArray{
                    Value: indev_level.Data,
                },
            },
        },
        "Entities": &nbt.List{
            Value: compoundArrayToTagArray(indev_level.Entities),
            ListType: ternary[int8](len(indev_level.Entities) == 0, nbt.IDTagEnd, nbt.IDTagCompound),
        },
        "TileEntities": &nbt.List{
            ListType: nbt.IDTagEnd,
        },
    })

    stream := nbt.NewStream(nbt.BigEndian)

    err := stream.WriteTag(root)
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

func (indev_level *IndevLevel) getHightestTile(x int32, z int32) int32 {
    for y := int32(indev_level.Height) - 1; y >= 0; y-- {
        index := (y * int32(indev_level.Length) + z) * int32(indev_level.Width) + x
        if indev_level.Blocks[index] != 0 && indev_level.Blocks[index] != 8 && indev_level.Blocks[index] != 9 && indev_level.Blocks[index] != 10 && indev_level.Blocks[index] != 11 {
            return y
        }
    }
    return int32(indev_level.Height) / 2
}
