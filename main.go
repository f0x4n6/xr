// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://ernw.de/download/EventManipulation.pdf
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

var Signature = []byte{0x2A, 0x2A, 0, 0}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for ReadUntil(reader, Signature) {
		size1 := ReadUint32(reader)
		event := ReadBytes(reader, size1-4-4-4)
		size2 := ReadUint32(reader)

		if size1 == size2 {
			fmt.Printf("[+] found record : %d [0x%x] : %d [0x%x]\n", size1, size1, size2, size2)
		} else {
			fmt.Printf("[!] found record : %d [0x%x] : %d [0x%x]\n", size1, size1, size2, size2)
		}

		fmt.Printf("\n%s\n", hex.Dump(event))
	}
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, 4))
}

func ReadBytes(r io.Reader, n uint32) []byte {
	b := make([]byte, n)

	if n, err := r.Read(b); err == nil {
		return b[:n]
	} else {
		panic(err)
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
