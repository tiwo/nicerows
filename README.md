# nicerows

Get results from database/sql as slices or maps, ready to be serialized to JSON.

Example:

```go
    result, err := db.Query(`Select "userid", "name" from "users";`)
    it := nicerows.New(result, err).IterateMaps()
    for row := range it {
        fmt.Printf("Row: %#v", row)
    }
```