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

const ChunkSize = 65536
const HeaderSize = 14
const tOffset1 = 6
const tOffset2 = 18

var cache = map[uint32]string{}
var last string

type Record struct {
	Id       uint64
	Time     uint64
	Computer string
	EventId  uint16
}

func NewRecord(r io.Reader) *Record {
	size := ReadUint32(r)

	record := Record{
		Id:   ReadUint64(r),
		Time: ReadUint64(r),
	}

	// read event stream
	b := ReadBytes(r, size-4-4-4-8-8)

	// check if copy equals size
	if size != ReadUint32(r) {
		return nil
	}

	// check if size is sane
	if size < HeaderSize || size > ChunkSize {
		return nil
	}

	// check if time is valid
	if record.Time == 0 {
		return nil
	}

	// check if stream length is valid
	if len(b) < 22 {
		return nil
	}

	tid1 := binary.BigEndian.Uint32(b[tOffset1 : tOffset1+4])
	tid2 := binary.BigEndian.Uint32(b[tOffset2 : tOffset2+4])

	// set computer by cached template
	if v, ok := cache[tid1]; ok {
		record.Computer = v
	} else {
		record.Computer = last + "?"
	}

	// skip fragment header (unused)
	r2 := bytes.NewReader(b[HeaderSize:])

	items := ReadUint32(r2)

	// invalid substitution array length
	if items > 20 {
		return nil
	}

	itemOffset := HeaderSize + 4 + 4 + 4
	itemSize := binary.LittleEndian.Uint16(b[itemOffset : itemOffset+2])
	itemType := b[itemOffset+2 : itemOffset+3][0]
	itemNull := b[itemOffset+3 : itemOffset+4][0]

	if itemType == 0x06 && itemNull == 0 {
		if len(b) > itemOffset+int(itemSize) {
			offset := HeaderSize + 4 + (items * 4) + 4

			//fmt.Printf("%x %x %x %x\n", offset, itemSize, itemType, itemNull)

			record.EventId = binary.LittleEndian.Uint16(b[offset : int(offset)+int(itemSize)])
		}
	}

	if tid1 != tid2 {
		return &record
	}

	// record is a template instance and carries a computer name to be cached
	if i := bytes.Index(b, computer); i >= 0 {
		if j := bytes.Index(b[i:], []byte{0x05, 0x01}); j >= 0 {
			l := int(binary.LittleEndian.Uint16(b[i+j+2:i+j+4]) * 2)
			last = string(b[i+j+4 : i+j+4+l])
			cache[tid1] = last
			record.Computer = last
		}
	}

	return &record
}

func main() {
	r := bufio.NewReader(os.Stdin)

	for ReadUntil(r, []byte{0x2A, 0x2A, 0, 0}) {
		if record := NewRecord(r); record != nil {
			fmt.Printf("%s %s %d\n",
				FileTime(record.Time).UTC().Format(time.RFC3339),
				record.Computer,
				record.EventId,
			)
		}
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
