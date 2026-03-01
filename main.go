// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://ernw.de/download/EventManipulation.pdf
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
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
	"strings"
	"time"
)

var Signature = []byte{0x2A, 0x2A, 0, 0}

type Record struct {
	Id    uint64
	Time  time.Time
	Size  uint32
	Copy  uint32
	Event []byte
}

func NewRecord(reader io.Reader) *Record {
	record := Record{
		Size: ReadUint32(reader),
		Id:   ReadUint64(reader),
		Time: ReadTime(reader),
	}

	record.Event = ReadBytes(reader, record.Size-4-4-4-8-8)
	record.Copy = ReadUint32(reader)

	// @0x04 = 0x01 "normal" = alles zwischen 01 und 02 ist der Name als String ohne NULL
	// 		   01 name 02 ...04 00 = <name>...</name>

	// @0x04 = 0x0C template instance
	//

	// Plan: Alle token types parsen, map aufbauen mit Namen als Key und byte array als value
	// 		 Wenn template definiert wird, dann template map aufbauen mit key und value
	//       Wenn template verwendet wird, dann template auflösen und einsetzen.

	switch record.Event[0x04] {
	case 0x01:
		// normal
	case 0x0C:
		// template instance
	}

	return &record
}

func (r *Record) String() string {
	return hex.Dump(r.Event)
}

func (r *Record) Header() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ID:   %d\n", r.Id))
	sb.WriteString(fmt.Sprintf("Time: %s\n", r.Time.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Size: %d [0x%x]\n", r.Size, r.Size))
	sb.WriteString(fmt.Sprintf("Copy: %d [0x%x]\n", r.Copy, r.Copy))

	return sb.String()
}

func (r *Record) IsValid() bool {
	return r.Size == r.Copy
}

func main() {
	var c byte
	var v, i uint64

	reader := bufio.NewReader(os.Stdin)

	for ReadUntil(reader, Signature) {
		record := NewRecord(reader)

		if record.IsValid() {
			c, v = '+', v+1
		} else {
			c, i = '!', i+1
		}

		fmt.Printf("[%c] found record #%d\n", c, record.Id)

		fmt.Printf("\n%s\n%s\n", record.Header(), record.String())
	}

	fmt.Printf("[=] found %d (valid) / %d (invalid) records\n", v, i)
}

func ReadTime(r io.Reader) time.Time {
	t := int64(ReadUint64(r))
	t -= 116444736000000000
	t *= 100

	return time.Unix(0, t)
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
