// Experimental record analysis.
//
// Usage:
//
//	INPUT | xr | OUTPUT
//
// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://github.com/libyal/libfwnt/blob/main/documentation/Security%20Descriptor.asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://blog.fox-it.com/2017/12/08/detection-and-recovery-of-nsas-covered-up-tracks/
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
// https://ernw.de/download/EventManipulation.pdf
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	Chunk     = 65536
	Header    = 14
	TempID1   = 6
	TempID2   = 18
	EventId   = Header + 12
	MinStream = 22
	MaxItems  = 20
	Layout    = "2006.01.02T15:04:05.0000000"
)

var Hash = []byte{0x3B, 0x6E}

var cache = map[string]string{}

var u4 = make([]byte, 4)
var u8 = make([]byte, 8)

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
			os.Exit(1)
		}
	}()

	var timestamp uint64
	var template string
	var computer string
	var eventId uint16

	r := bufio.NewReaderSize(os.Stdin, Chunk)

	for ReadUntil(r, 0x00002A2A) {
		size := ReadUint32(r)

		// check if size is sane
		if size < Header || size > Chunk {
			continue
		}

		// skip record id
		_ = ReadUint64(r)

		timestamp = ReadUint64(r)

		// check if time is valid
		if timestamp == 0 {
			continue
		}

		// read binxml stream
		b := ReadBytes(r, make([]byte, size-4-4-4-8-8))

		// check if stream length is valid
		if len(b) < MinStream {
			continue
		}

		// check if size equals copy
		if size != ReadUint32(r) {
			continue
		}

		// get substitution array length
		items := binary.LittleEndian.Uint32(b[Header : Header+4])

		// check array length
		if items > MaxItems {
			continue
		}

		// check if type and null are valid
		if b[EventId+2] != 0x06 || b[EventId+3] != 0x00 {
			continue
		}

		offset := Header + 4 + (items * 4) + 2 + 1 + 1
		eventId = binary.LittleEndian.Uint16(b[offset : offset+2])
		template = string(b[TempID1 : TempID1+4])

		// record is a template instance and has a cacheable computer name
		if template == string(b[TempID2:TempID2+4]) {
			if i := bytes.Index(b, Hash); i >= 0 {
				if j := bytes.IndexByte(b[i:], 0x02); j >= 0 {
					offset := i + j + 5
					length := int(binary.LittleEndian.Uint16(b[offset-2:offset]) * 2)
					cache[template] = string(b[offset : offset+length])
				}
			}
		}

		// get cached computer name by template id
		if v, ok := cache[template]; ok {
			computer = v
		}

		if len(computer) > 255 {
			computer = computer[:255]
		}

		fmt.Printf("XR|%s|%s|%d\n", FileTime(timestamp).UTC().Format(Layout), computer, eventId)
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100)
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, u8))
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, u4))
}

func ReadBytes(r io.Reader, b []byte) []byte {
	if n, err := r.Read(b); err == nil {
		return b[:n]
	} else {
		panic(err)
	}
}

func ReadUntil(r io.Reader, v uint32) bool {
	buf := make([]byte, 4)

	for binary.LittleEndian.Uint32(buf) != v {
		switch _, err := io.ReadFull(r, buf); {
		case errors.Is(err, io.ErrUnexpectedEOF):
			fallthrough
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}

	return true
}
