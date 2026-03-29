package internal

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"
)

var Signature = []byte{0x2A, 0x2A, 0x00, 0x00}

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
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Id:   %d\n", r.Id))
	sb.WriteString(fmt.Sprintf("Time: %s\n", r.Time.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Size: %d [0x%x]\n", r.Size, r.Size))
	sb.WriteString(fmt.Sprintf("Copy: %d [0x%x]\n\n", r.Copy, r.Copy))
	sb.WriteString(hex.Dump(r.Event))

	return sb.String()
}

func (r *Record) IsValid() bool {
	return r.Size == r.Copy
}
