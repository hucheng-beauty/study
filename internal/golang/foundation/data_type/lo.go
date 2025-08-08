package main

import (
    "fmt"

    "github.com/samber/lo"
)

func Lo() {
    var x int

    fmt.Println()
    fmt.Println("[lo.Ternary]")
    fmt.Println("x ==> ", x)
    fmt.Printf("ternary ==> %v\n", lo.Ternary(x > 10, true, false))

    m := map[string]int{"a": 1, "b": 2, "c": 3}

    fmt.Println()
    fmt.Println("[lo.Keys、lo.Values]")
    keys := lo.Keys(m)
    values := lo.Values(m)
    fmt.Println("map ==> ", m)
    fmt.Printf("keys ==> %v\n", keys)
    fmt.Printf("values ==> %v\n ", values)

    fmt.Println()
    fmt.Println("[lo.MapKeys、lo.MapValues]")
    fmt.Println("map ==> ", m)
    mapKeys := lo.MapKeys(m, func(value int, key string) string {
        return fmt.Sprintf("key_%s", key)
    })
    fmt.Printf("mapKeys ==> %v\n", mapKeys)

    mapValues := lo.MapValues(m, func(value int, key string) string {
        return fmt.Sprintf("value_%d", value)
    })
    fmt.Printf("mapValues ==> %v\n", mapValues)

    s := []string{"c++", "c", "golang", "java", "python", "golang"}

    fmt.Println()
    fmt.Println("[lo.Uniq]")
    uniqSli := lo.Uniq(s)
    fmt.Println("sli ==> ", s)
    fmt.Println("uniqSli ==> ", uniqSli)

    fmt.Println()
    fmt.Println("[lo.ForEach slice]")
    fmt.Println("sli ==> ", s)
    lo.ForEach(s, func(item string, index int) {
        fmt.Printf("index ==> %d, item ==> %s\n", index, item)
    })

    fmt.Println()
    fmt.Println("[lo.Entries] map => entries")
    entries := lo.Entries(m) // map => slice
    fmt.Printf("map ==> %v \n", m)
    fmt.Printf("entries ==> %v\n", entries)

    fmt.Println()
    fmt.Println("[lo.ForEach map]")
    fmt.Println("map ==> ", entries)
    lo.ForEach(entries, func(e lo.Entry[string, int], _ int) {
        fmt.Printf("key ==> %s, value ==> %d\n", e.Key, e.Value)
    })

    fmt.Println()
    fmt.Println("[lo.MapEntries]")
    mapEntries := lo.MapEntries(m, func(key string, value int) (string, int) {
        return fmt.Sprintf("key_%s", key), value
    })
    fmt.Printf("mapEntries ==> %v\n", mapEntries)

    fmt.Println()
    fmt.Println("[lo.Map]")
    mapSli := lo.Map(s, func(item string, _ int) bool {
        return item == "golang"
    })
    fmt.Println("sli ==> ", s)
    fmt.Printf("mapSli ==> %v\n", mapSli)

    fmt.Println()
    fmt.Println("[lo.Filter]")
    filterSli := lo.Filter(s, func(item string, _ int) bool {
        return item == "golang"
    })
    fmt.Println("sli ==> ", s)
    fmt.Printf("filterSli ==> %v\n", filterSli)

    fmt.Println()
    fmt.Println("[lo.Contains]")
    fmt.Println("sli ==> ", s)
    contains := lo.Contains(s, "golang")
    fmt.Println("lo.Contains(s,'golang') ==>", contains)
}
