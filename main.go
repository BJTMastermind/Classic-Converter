package main

import (
    "bytes"
    "compress/gzip"
    "encoding/binary"
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/BJTMastermind/Classic-Converter/classic_converter"
    "github.com/BJTMastermind/Go-MC-Classic-Parser"
    "github.com/BJTMastermind/go-nbt"
    "github.com/akamensky/argparse"
)

func main() {
    argparser := argparse.NewParser("Classic-Converter", "Convert classic worlds to indev worlds or schematics!")

    input := argparser.File("i", "input", os.O_RDONLY, os.ModePerm, &argparse.Options{Required: true, Help: "The classic world file to parse."})
    format := argparser.String("f", "format", &argparse.Options{Required: true, Help: "The output format to use. One of \"indev_level\" or \"schematic\"."})

    err := argparser.Parse(os.Args)
    if err != nil {
        fmt.Print(argparser.Usage(err))
        return
    }

    if !strings.HasSuffix(input.Name(), ".dat") && !strings.HasSuffix(input.Name(), ".mine") {
        fmt.Print(argparser.Usage("Input file must be a classic world file. (.dat or .mine)"))
        return
    }

    if *format != "indev_level" && *format != "schematic" {
        fmt.Print(argparser.Usage("Output format must be one of \"indev_level\" or \"schematic\"."))
        return
    }

    if *format == "indev_level" {
        fmt.Println("Converting to a indev level...")
        err := convertToIndevLevel(*input)
        if err != nil {
            fmt.Println(err)
        }
    } else {
        fmt.Println("Converting to a schematic...")
        err = convertToSchematic(*input)
        if err != nil {
            fmt.Println(err)
        }
    }
}

func convertToIndevLevel(inputFile os.File) (error) {
    // Check that given file is a gzipped file
    gzbytes, _ := os.ReadFile(inputFile.Name())

    if gzMagic := binary.BigEndian.Uint16(gzbytes[0:2]); gzMagic != 0x1f8b {
        return errors.New("Not a GZIP file.")
    }

    // Get uncompressed file size
    uncompressedSize := binary.LittleEndian.Uint32(gzbytes[len(gzbytes) - 4:])

    // Decompress Gzip
    gzReader, _ := gzip.NewReader(bytes.NewBuffer(gzbytes))
    defer gzReader.Close()

    var uncompressedBytes = make([]byte, uncompressedSize)
    binary.Read(gzReader, binary.BigEndian, &uncompressedBytes)

    reader := bytes.NewBuffer(uncompressedBytes)

    magic := int32(binary.BigEndian.Uint32(reader.Next(4)))
    version := reader.Next(1)[0]

    indevLevel := new(classic_converter.IndevLevel).InitWithDefaults()

    // Check if a classic world
    fmt.Println("Figuring out what classic version the world is...")
    if magic != 0x271bb788 {
        // Check if a pre classic world
        if len(uncompressedBytes) != (256*256*64) {
            return errors.New("error: Not a vaild Minecraft Pre-Classic save, Byte array is not equal to 4,194,304 bytes.")
        }

        for i := 0; i < (256*256*64); i++ {
            if uncompressedBytes[i] < 0 || uncompressedBytes[i] > 49 {
                return errors.New("error: Not a vaild Minecraft Pre-Classic save, Byte array contains block IDs greater then 49.")
            }
        }

        // Vaild Pre-Classic save
        fmt.Println("Found pre-classic world format!")

        indevLevel.Blocks = classic_converter.ByteArray2Int8Array(uncompressedBytes)

        indevLevel.FindSpawn()
        indevLevel.WriteToFile(strings.NewReplacer(".dat", ".mclevel", ".mine", ".mclevel").Replace(inputFile.Name()))

        return nil
    }

    if version != 0x01 && version != 0x02 {
        return errors.New(fmt.Sprintf("error: Not a supported classic format version. Got %d, Expected 1 or 2\n", version))
    }

    if version == 0x01 {
        fmt.Println("Found classic version 1 world format!")

        worldNameLength := binary.BigEndian.Uint16(reader.Next(2))
        indevLevel.Name = string(reader.Next(int(worldNameLength)))

        creatorNameLength := binary.BigEndian.Uint16(reader.Next(2))
        indevLevel.Author = string(reader.Next(int(creatorNameLength)))

        indevLevel.CreatedOn = int64(binary.BigEndian.Uint16(reader.Next(8)))
        indevLevel.Width = int16(binary.BigEndian.Uint16(reader.Next(2)))
        indevLevel.Length = int16(binary.BigEndian.Uint16(reader.Next(2)))
        indevLevel.Height = int16(binary.BigEndian.Uint16(reader.Next(2)))
        indevLevel.Blocks = classic_converter.ByteArray2Int8Array(reader.Bytes())

        indevLevel.FindSpawn()
        indevLevel.WriteToFile(strings.NewReplacer(".dat", ".mclevel", ".mine", ".mclevel").Replace(inputFile.Name()))
    } else if version == 0x02 {
        fmt.Println("Found classic version 2 world format!")

        parser := new(mc_classic_parser.ClassicParser)

        world, err := parser.ParseBytes(reader.Bytes())
        if err != nil {
            return err
        }

        indevLevel.CreatedOn = world.CreateTime
        indevLevel.Name = world.Name
        indevLevel.Author = world.Creator
        indevLevel.SkyColor = world.SkyColor
        indevLevel.FogColor = world.FogColor
        indevLevel.CloudColor = world.CloudColor
        indevLevel.Width = int16(world.Width)
        indevLevel.Length = int16(world.Depth)
        indevLevel.Height = int16(world.Height)
        indevLevel.Spawn = [3]int16{int16(world.XSpawn), int16(world.YSpawn), int16(world.ZSpawn)}
        indevLevel.Blocks = world.Blocks

        compoundEntities := []nbt.Compound{}
        for _, entity := range world.Entities {
            if entity.TextureName == "/char.png" {
                continue
            }

            compoundEntity := classic_converter.ClassicEntity2Compound(entity, false)
            compoundEntities = append(compoundEntities, compoundEntity)
        }
        compoundPlayer := classic_converter.ClassicPlayer2Compound(world.Player)
        compoundEntities = append(compoundEntities, compoundPlayer)

        indevLevel.Entities = compoundEntities

        indevLevel.FindSpawn()
        indevLevel.WriteToFile(strings.NewReplacer(".dat", ".mclevel", ".mine", ".mclevel").Replace(inputFile.Name()))
    }
    return nil
}

func convertToSchematic(inputFile os.File) (error) {
    // Check that given file is a gzipped file
    gzbytes, _ := os.ReadFile(inputFile.Name())

    if gzMagic := binary.BigEndian.Uint16(gzbytes[0:2]); gzMagic != 0x1f8b {
        return errors.New("Not a GZIP file.")
    }

    // Get uncompressed file size
    uncompressedSize := binary.LittleEndian.Uint32(gzbytes[len(gzbytes) - 4:])

    // Decompress Gzip
    gzReader, _ := gzip.NewReader(bytes.NewBuffer(gzbytes))
    defer gzReader.Close()

    var uncompressedBytes = make([]byte, uncompressedSize)
    binary.Read(gzReader, binary.BigEndian, &uncompressedBytes)

    reader := bytes.NewBuffer(uncompressedBytes)

    magic := int32(binary.BigEndian.Uint32(reader.Next(4)))
    version := reader.Next(1)[0]

    schematic := new(classic_converter.Schematic).InitWithDefaults()

    // Check if a classic world
    fmt.Println("Figuring out what classic version the world is...")
    if magic != 0x271bb788 {
        // Check if a pre classic world
        if len(uncompressedBytes) != (256*256*64) {
            return errors.New("error: Not a vaild Minecraft Pre-Classic save, Byte array is not equal to 4,194,304 bytes.")
        }

        for i := 0; i < (256*256*64); i++ {
            if uncompressedBytes[i] < 0 || uncompressedBytes[i] > 49 {
                return errors.New("error: Not a vaild Minecraft Pre-Classic save, Byte array contains block IDs greater then 49.")
            }
        }

        // Vaild Pre-Classic save
        fmt.Println("Found pre-classic world format!")

        schematic.Blocks = classic_converter.ByteArray2Int8Array(uncompressedBytes)

        schematic.WriteToFile(strings.NewReplacer(".dat", ".schematic", ".mine", ".schematic").Replace(inputFile.Name()))
    }

    if version != 0x01 && version != 0x02 {
        return errors.New(fmt.Sprintf("error: Not a supported classic format version. Got %d, Expected 1 or 2\n", version))
    }

    if version == 0x01 {
        fmt.Println("Found classic version 1 world format!")

        schematic.Width = int16(binary.BigEndian.Uint16(reader.Next(2)))
        schematic.Length = int16(binary.BigEndian.Uint16(reader.Next(2)))
        schematic.Height = int16(binary.BigEndian.Uint16(reader.Next(2)))
        schematic.Blocks = classic_converter.ByteArray2Int8Array(reader.Bytes())

        schematic.WriteToFile(strings.NewReplacer(".dat", ".schematic", ".mine", ".schematic").Replace(inputFile.Name()))

    } else if version == 0x02 {
        fmt.Println("Found classic version 2 world format!")

        parser := new(mc_classic_parser.ClassicParser)

        world, err := parser.ParseBytes(reader.Bytes())
        if err != nil {
            return err
        }

        schematic.Width = int16(world.Width)
        schematic.Length = int16(world.Depth)
        schematic.Height = int16(world.Height)
        schematic.Blocks = world.Blocks

        compoundEntities := []nbt.Compound{}
        for _, entity := range world.Entities {
            if entity.TextureName == "/char.png" {
                continue
            }

            compoundEntity := classic_converter.ClassicEntity2Compound(entity, true)
            compoundEntities = append(compoundEntities, compoundEntity)
        }
        schematic.Entities = compoundEntities

        schematic.WriteToFile(strings.NewReplacer(".dat", ".schematic", ".mine", ".schematic").Replace(inputFile.Name()))
    }
    return nil
}
