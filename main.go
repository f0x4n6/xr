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

var computer = []byte{
	0x08, 0, 0x43, 0, 0x6f, 0, 0x6d, 0, 0x70, 0, 0x75, 0, 0x74, 0, 0x65, 0, 0x72, 0, 0, 0,
}

const chunk = 65536
const head = 14
const item = head + 12
const tid1 = 6
const tid2 = 18

func main() {
	var ts uint64
	var ev uint16
	var cn string

	var cache = map[string]string{}
	var last string

	r := bufio.NewReader(os.Stdin)

	temp := make([]byte, 4)

	for ReadUntil(r, []byte{0x2A, 0x2A, 0, 0}) {
		_, _ = r.Read(temp)

		size := binary.LittleEndian.Uint32(temp)

		rec := make([]byte, size)

		cp := binary.LittleEndian.Uint32(rec[len(rec)-4:])

		// check if copy equals size
		if size != cp {
			continue
		}

		// check if size is sane
		if size < head || size > chunk {
			continue
		}

		ts = binary.LittleEndian.Uint64(rec[12:20])

		// check if time is valid
		if ts == 0 {
			continue
		}

		// read event stream
		b := rec[20 : size-4-4-4-8-8]

		// check if stream length is valid
		if len(b) < 22 {
			continue
		}

		// skip fragment header (unused)
		items := binary.LittleEndian.Uint32(b[head : head+4])

		// invalid substitution array length
		if items > 20 {
			continue
		}

		// check if item type and null is valid
		if b[item+2] == 0x06 && b[item+3] == 0x00 {
			offset := head + 4 + (items * 4) + 2 + 1 + 1
			ev = uint16(b[offset]) | uint16(b[offset+1])<<8
		}

		t1 := string(b[tid1 : tid1+4])
		t2 := string(b[tid2 : tid2+4])

		// set computer by cached template
		if v, ok := cache[t1]; ok {
			cn = v
		} else {
			cn = last + "?"
		}

		// record is a template instance and carries a computer name to be cached
		if t1 == t2 {
			if i := bytes.Index(b, computer); i >= 0 {
				if j := bytes.Index(b[i:], []byte{0x05, 0x01}); j >= 0 {
					l := int(binary.LittleEndian.Uint16(b[i+j+2:i+j+4]) * 2)
					last = string(b[i+j+4 : i+j+4+l])
					cache[t1] = last
					cn = last
				}
			}
		}

		fmt.Printf("%s %s %d\n", FileTime(ts).UTC().Format(time.RFC3339), cn, ev)
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100)
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
