// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://ernw.de/download/EventManipulation.pdf
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

const BLOCK = 65536 // chunk size

var Magic = []byte{0x2A, 0x2A, 0, 0}

func main() {
	last := make([]byte, len(Magic))
	buffer := make([]byte, BLOCK)
	reader := bufio.NewReaderSize(os.Stdin, BLOCK)

	for {
		n, err := io.ReadFull(reader, buffer)

		buffer = buffer[:n] // shrink

		if n >= 4 {
			copy(last, buffer[n-4:])
		} else {
			last = last[:0]
		}

		fmt.Printf("Read %d Buffer [%d] Last [%d] %x\n", n, len(buffer), len(last), last)

		//fmt.Printf("found\n")

		if errors.Is(err, io.EOF) {
			break
		}
	}
}
