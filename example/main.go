package main

import (
	"fmt"
	"time"

	"github.com/gosuri/uilive"
)

func main() {
	writer := uilive.New()

	// start listening for updates and render
	writer.Start()

	for _, f := range [][]string{{"Foo.zip", "Bar.iso"}, {"Baz.tar.gz", "Qux.img"}} {
		for i := 0; i <= 50; i++ {
			_, _ = fmt.Fprintf(writer, "Downloading %s.. (%d/%d) GB\n", f[0], i, 50)
			_, _ = fmt.Fprintf(writer.Newline(), "Downloading %s.. (%d/%d) GB\n", f[1], i, 50)
			time.Sleep(time.Millisecond * 25)
		}
		_, _ = fmt.Fprintf(writer.Bypass(), "Downloaded %s\n", f[0])
		_, _ = fmt.Fprintf(writer.Bypass(), "Downloaded %s\n", f[1])
	}
	_, _ = fmt.Fprintln(writer, "Finished: Downloaded 150GB")
	writer.Stop() // flush and stop rendering
}
