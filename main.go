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

const Chunk = 65536
const Header = 14
const TempID1 = 6
const TempID2 = 18
const EventData = Header + 12
const EventType = 0x06
const EventNull = 0x00

var Magic = []byte{'*', '*', 0, 0}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	var cache = map[string]string{}

	var timestamp uint64
	var computer string
	var eventId uint16
	var key string

	r := bufio.NewReaderSize(os.Stdin, Chunk)

	for ReadUntil(r, Magic) {
		size := ReadUint32(r)

		// check if size is sane
		if size < Header || size > Chunk {
			continue
		}

		_ = ReadUint64(r)

		timestamp = ReadUint64(r)

		// check if time is valid
		if timestamp == 0 {
			continue
		}

		b := ReadBytes(r, size-4-4-4-8-8)

		// check if stream length is valid
		if len(b) < 22 {
			continue
		}

		// check if copy equals size
		if size != ReadUint32(r) {
			continue
		}

		// skip fragment header (unused)
		n := binary.LittleEndian.Uint32(b[Header : Header+4])

		// check substitution array length
		if n > 20 {
			continue
		}

		// check if EventData type and null is valid
		if b[EventData+2] == EventType && b[EventData+3] == EventNull {
			i := Header + 4 + (n * 4) + 2 + 1 + 1
			eventId = binary.LittleEndian.Uint16(b[i : i+2])
		}

		key = string(b[TempID1 : TempID1+4])

		// record is a template instance and carries a computer name to be cached
		if key == string(b[TempID2:TempID2+4]) {
			if i := bytes.Index(b, []byte{0x3B, 0x6E}); i >= 0 {
				if j := bytes.IndexByte(b[i:], 0x02); j >= 0 {
					l := int(binary.LittleEndian.Uint16(b[i+j+3:i+j+5]) * 2)
					cache[key] = string(b[i+j+5 : i+j+5+l])
				}
			}
		}

		// set computer by cached template
		if v, ok := cache[key]; ok {
			computer = v
		}

		fmt.Printf("%s %s %d\n", FileTime(timestamp).UTC().Format(time.RFC3339), computer, eventId)
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100)
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
		case errors.Is(err, io.ErrUnexpectedEOF):
			return false
		case err != nil:
			panic(err)
		}
	}

	return true
}
