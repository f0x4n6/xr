// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://ernw.de/download/EventManipulation.pdf
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var Magic = []byte{0x2A, 0x2A, 0, 0}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for ReadUntil(reader, Magic) {
		fmt.Printf("found\n")
	}
}

func ReadUntil(r io.Reader, m []byte) bool {
	b := make([]byte, len(m))

	for !bytes.Equal(b, m) {
		switch _, err := io.ReadFull(r, b); {
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}

	return true
}
