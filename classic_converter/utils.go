package classic_converter

import (
    "github.com/BJTMastermind/Go-MC-Classic-Parser"
    "github.com/BJTMastermind/go-nbt"
)

func ternary[T any](condition bool, a T, b T) T {
    if condition {
        return a
    }
    return b
}

func compoundArrayToTagArray(compoundArray []nbt.Compound) []nbt.Tag {
    out := make([]nbt.Tag, len(compoundArray))
    for i, compound := range compoundArray {
        copy := compound
        out[i] = &copy
    }
    return out
}

func textureName2Id(textureName string) string {
    switch textureName {
        case "/mob/zombie.png":
            return "Zombie"
        case "/mob/skeleton.png":
            return "Skeleton"
        case "/mob/creeper.png":
            return "Creeper"
        case "/mob/spider.png":
            return "Spider"
        case "/mob/pig.png":
            return "Pig"
        case "/mob/sheep.png":
            return "Sheep"
        default:
            return "Zombie"
    }
}

func inventory2List(player mc_classic_parser.ClassicPlayer) []nbt.Tag {
    var out []nbt.Tag

    inventory := player.Inventory
    slots := inventory["slots"].([]int32)
    count := inventory["count"].([]int32)

    for i := 0; i < 9; i++ {
        if slots[i] == -1 {
            continue
        }

        out = append(out, &nbt.Compound{
            Value: map[string]nbt.Tag{
                "Slot": &nbt.Byte{
                    Value: int8(i),
                },
                "id": &nbt.Short{
                    Value: int16(slots[i]),
                },
                "Damage": &nbt.Short{
                    Value: 0,
                },
                "Count": &nbt.Byte{
                    Value: int8(count[i]),
                },
            },
        })
    }
    // Add Arrows to players inventory
    if player.Arrows > 0 {
        out = append(out, &nbt.Compound{
            Value: map[string]nbt.Tag{
                "Slot": &nbt.Byte{
                    Value: 9,
                },
                "id": &nbt.Short{
                    Value: 262,
                },
                "Damage": &nbt.Short{
                    Value: 0,
                },
                "Count": &nbt.Byte{
                    Value: int8(player.Arrows),
                },
            },
        })
    }
    // Add Bow to players inventory
    out = append(out, &nbt.Compound{
        Value: map[string]nbt.Tag{
            "Slot": &nbt.Byte{
                Value: 10,
            },
            "id": &nbt.Short{
                Value: 261,
            },
            "Damage": &nbt.Short{
                Value: 0,
            },
            "Count": &nbt.Byte{
                Value: 1,
            },
        },
    })

    return out
}

func ByteArray2Int8Array(array []byte) []int8 {
    out := make([]int8, len(array))
    for i, b := range array {
        out[i] = int8(b)
    }
    return out
}

func ClassicEntity2Compound(entity mc_classic_parser.ClassicEntity, schematic bool) nbt.Compound {
    output := nbt.Compound{
        Value: map[string]nbt.Tag{
            "id": &nbt.String{
                Value: textureName2Id(entity.TextureName),
            },
            "Pos": &nbt.List{
                Value: []nbt.Tag{
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.X)}, &nbt.Float{Value: entity.X}),
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.Y)}, &nbt.Float{Value: entity.Y}),
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.Z)}, &nbt.Float{Value: entity.Z}),
                },
                ListType: ternary[int8](schematic, nbt.IDTagDouble, nbt.IDTagFloat),
            },
            "Rotation": &nbt.List{
                Value: []nbt.Tag{
                    &nbt.Float{Value: entity.YRot},
                    &nbt.Float{Value: entity.XRot},
                },
                ListType: nbt.IDTagFloat,
            },
            "Motion": &nbt.List{
                Value: []nbt.Tag{
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.Xd)}, &nbt.Float{Value: entity.Xd}),
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.Yd)}, &nbt.Float{Value: entity.Yd}),
                    ternary[nbt.Tag](schematic, &nbt.Double{Value: float64(entity.Zd)}, &nbt.Float{Value: entity.Zd}),
                },
                ListType: ternary[int8](schematic, nbt.IDTagDouble, nbt.IDTagFloat),
            },
            "FallDistance": &nbt.Float{
                Value: entity.FallDistance,
            },
            "Health": &nbt.Short{
                Value: int16(entity.Health),
            },
            "AttackTime": &nbt.Short{
                Value: int16(entity.AttackTime),
            },
            "HurtTime": &nbt.Short{
                Value: int16(entity.HurtTime),
            },
            "DeathTime": &nbt.Short{
                Value: int16(entity.DeathTime),
            },
            "Air": &nbt.Short{
                Value: int16(entity.AirSupply),
            },
            "Fire": &nbt.Short{
                Value: -1,
            },
        },
    }
    if schematic {
        output.Value["OnGround"] = &nbt.Byte{
            Value: ternary[int8](entity.OnGround, 1, 0),
        }
    }
    if textureName2Id(entity.TextureName) == "Sheep" {
        output.Value["Sheared"] = &nbt.Byte{
            Value: ternary[int8](entity.HasHair, 1, 0),
        }
    }

    return output
}

func ClassicPlayer2Compound(player mc_classic_parser.ClassicPlayer) nbt.Compound {
    return nbt.Compound{
        Value: map[string]nbt.Tag{
            "id": &nbt.String{
                Value: "LocalPlayer",
            },
            "Pos": &nbt.List{
                Value: []nbt.Tag{
                    &nbt.Float{Value: player.X},
                    &nbt.Float{Value: player.Y},
                    &nbt.Float{Value: player.Z},
                },
                ListType: nbt.IDTagFloat,
            },
            "Rotation": &nbt.List{
                Value: []nbt.Tag{
                    &nbt.Float{Value: player.YRot},
                    &nbt.Float{Value: player.XRot},
                },
                ListType: nbt.IDTagFloat,
            },
            "Motion": &nbt.List{
                Value: []nbt.Tag{
                    &nbt.Float{Value: player.Xd},
                    &nbt.Float{Value: player.Yd},
                    &nbt.Float{Value: player.Zd},
                },
                ListType: nbt.IDTagFloat,
            },
            "FallDistance": &nbt.Float{
                Value: player.FallDistance,
            },
            "Health": &nbt.Short{
                Value: int16(player.Health),
            },
            "AttackTime": &nbt.Short{
                Value: int16(player.AttackTime),
            },
            "HurtTime": &nbt.Short{
                Value: int16(player.HurtTime),
            },
            "DeathTime": &nbt.Short{
                Value: int16(player.DeathTime),
            },
            "Air": &nbt.Short{
                Value: int16(player.AirSupply),
            },
            "Fire": &nbt.Short{
                Value: -1,
            },
            "Score": &nbt.Int{
                Value: player.Score,
            },
            "Inventory": &nbt.List{
                Value: inventory2List(player),
                ListType: nbt.IDTagCompound,
            },
        },
    }
}
