// Experimental record analyzer.
//
// Usage:
//
//	INPUT | xr | OUTPUT
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
const Magic = 0x00002A2A
const Header = 28
const Layout = "2006.01.02T15:04:05.0000000"

var cache = make(map[string]string)
var buf0 = make([]byte, 0, Chunk)
var buf4 = make([]byte, 4)
var buf8 = make([]byte, 8)

func main() {
	var b []byte
	var ev uint16
	var sz, n uint32
	var id, ts uint64
	var k, cn string

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
			os.Exit(1)
		}
	}()

	for r := bufio.NewReaderSize(os.Stdin, Chunk); ReadUntil(r); {
		if sz = ReadUint32(r); sz < Header || sz > Chunk {
			continue // check sane size
		}

		if id = ReadUint64(r); id == 0 {
			continue // check valid record id
		}

		if ts = ReadUint64(r); ts == 0 {
			continue // check valid time
		}

		if b = ReadBytes(r, buf0[:sz-Header]); len(b) < 30 {
			continue // check valid binxml length
		}

		if sz != ReadUint32(r) {
			continue // check size equals copy
		}

		if n = binary.LittleEndian.Uint32(b[14:]); n > 20 {
			continue // check substitution array length
		}

		if b[28] == 6 && b[29] == 0 {
			ev = binary.LittleEndian.Uint16(b[(n*4)+22:])
		}

		if k = string(b[6:10]); k == string(b[18:22]) {
			if i := bytes.Index(b, []byte{0x3B, 0x6E}); i >= 0 {
				if j := bytes.IndexByte(b[i:], 0x02); j >= 0 {
					cache[k] = Utf16String(b[i+j+3:])
				}
			}
		}

		if v, hit := cache[k]; hit {
			cn = v // if not cached, use last computer name
		}

		fmt.Printf("XR|%s|%s|%d\n", FileTime(ts).Format(Layout), cn, ev)
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100).UTC()
}

func Utf16String(s []byte) string {
	return string(s[2 : 2+2*binary.LittleEndian.Uint16(s)])
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, buf8))
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, buf4))
}

func ReadBytes(r io.Reader, b []byte) []byte {
	if n, err := r.Read(b); err == nil {
		return b[:n]
	} else {
		panic(err)
	}
}

func ReadUntil(r io.Reader) bool {
	for binary.LittleEndian.Uint32(buf4) != Magic {
		switch _, err := io.ReadFull(r, buf4); {
		case errors.Is(err, io.ErrUnexpectedEOF):
			return false
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}

	return true
}

// Sources:
// https://github.com/libyal/libevtx/blob/main/documentation/Windows%20XML%20Event%20Log%20(EVTX).asciidoc
// https://github.com/libyal/libfwnt/blob/main/documentation/Security%20Descriptor.asciidoc
// https://www.researchgate.net/publication/222426407_Introducing_the_Microsoft_Vista_event_log_file_format
// https://blog.fox-it.com/2019/06/04/export-corrupts-windows-event-log-files/
// https://blog.fox-it.com/2017/12/08/detection-and-recovery-of-nsas-covered-up-tracks/
// https://parsiya.net/blog/2018-11-01-windows-filetime-timestamps-and-byte-wrangling-with-go/
// https://ernw.de/download/EventManipulation.pdf
// https://wtf-8.codeberg.page/
