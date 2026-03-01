// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://ernw.de/download/EventManipulation.pdf
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
// https://cs.opensource.google/go/x/sys/+/refs/tags/v0.41.0:windows/types_windows.go;l=803
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
	"time"
)

var Signature = []byte{0x2A, 0x2A, 0, 0}

func main() {
	var c byte
	var v, i uint64

	reader := bufio.NewReader(os.Stdin)

	for ReadUntil(reader, Signature) {
		size1 := ReadUint32(reader)
		recid := ReadUint64(reader)
		ftime := ReadTime(reader)
		event := ReadBytes(reader, size1-4-4-4-8-8)
		size2 := ReadUint32(reader)

		if size1 == size2 {
			c = '+'
			v++
		} else {
			c = '!'
			i++
		}

		fmt.Printf("[%c] found record #%d @%s : %d [0x%x] : %d [0x%x]\n",
			c,
			recid, ftime.Format(time.RFC3339),
			size1, size1,
			size2, size2,
		)

		fmt.Printf("\n%s\n", hex.Dump(event))
	}

	fmt.Printf("[=] found %d (valid) / %d (invalid) records\n", v, i)
}

func ReadTime(r io.Reader) time.Time {
	t := ReadBytes(r, 8)

	lt := binary.LittleEndian.Uint32(t[:4])
	ht := binary.LittleEndian.Uint32(t[4:])

	// 100-nanosecond intervals since January 1, 1601
	nsec := int64(ht)<<32 + int64(lt)
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= 116444736000000000
	// convert into nanoseconds
	nsec *= 100

	return time.Unix(0, nsec)
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, 8))
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
