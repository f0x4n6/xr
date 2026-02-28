// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
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
	buffer := make([]byte, 4)
	reader := bufio.NewReader(os.Stdin)

	for offset := 0; ; offset += 4 {
		if _, err := reader.Read(buffer); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			panic(err)
		}

		if bytes.Equal(buffer, Magic) {
			println(offset)
		}
	}
}
