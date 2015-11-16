# uilive [![GoDoc](https://godoc.org/github.com/gosuri/uilive?status.svg)](https://godoc.org/github.com/gosuri/uilive)

uilive is a go library for refreshing ui in a timed interval for terminal applications

## Example

Full source for the below example is in [example/main.go](example/main.go).

```go
writer := uilive.New()
writer.Start() // start listening

for i := 0; i <= 100; i++ {
  fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i, 100)
  time.Sleep(time.Millisecond * 5)
}

fmt.Fprintln(writer, "Finished: Downloaded 100GB")
writer.Wait() // wait for writer to finish writing
```

The above will render

![example](doc/example.gif)
