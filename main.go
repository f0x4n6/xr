// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

var Magic = []byte{0x2A, 0x2A, 0, 0}

func main() {
	b := make([]byte, 4)
	r := bufio.NewReader(os.Stdin)

	for i := 0; ; i += 4 {
		if _, err := r.Read(b); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			panic(err)
		}

		if bytes.Equal(b, Magic) {
			println(i)
		}
	}
}
