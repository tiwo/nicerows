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


## To Do

- [ ] Expose the error state (via an Err() error method)
- [ ] Make the first parameter of New(), sqlresult, an interface to remove the
      database/sql dependency
- [ ] Create a mock sqlresult type that can simulate syntax/database/connection
      errors
- [ ] Use that mock type in testing.