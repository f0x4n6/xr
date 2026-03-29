package internal

import (
	"encoding/binary"
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

func NewRecord(r io.Reader) *Record {
	record := Record{
		Size: ReadUint32(r),
		Id:   ReadUint64(r),
		Time: ReadTime(r),
	}

	record.Event = ReadBytes(r, record.Size-4-4-4-8-8)
	record.Copy = ReadUint32(r)

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

	// Each record starts with a system element
	/*
		<Event>
			<System>
			</System>
			<...>
		</Event>

		Fragment
			Fragment header
			0f header token
			01 major version
			01 minor version
			00 flags

			Element or Template instance
			| Template instance
				0C template instance token
				... temp definition
					[1] unknown
					[4] unknown
					[4] temp def data (or offset)
					[4] unknown (0 = not used)
					[16] temp id (guid)
					[4] data size
					... fragment header
					... element
					[1] end of file token (00)

				... temp inst data
					[4] number of temp values
					... array temp value desc
						[2] value size
						[1] value type
						00  unknown (empty?)

					... array temp value data

			| Element (empty of filled)
				| Empty
					element start
					close empty element token
				| Filled
					01|41 element start (01 = no element, 41 attribute list exists)
						? 0xffff (optional, dependency identifier)
						Data size
						? element name offset
						Attribute list
					close start element token
					content
					end element token

					Name
						[4] unknown
						[2] name hash
						[2] length (chars)
						UTF-16 little-endian string with an end-of-string character
							Idea: parse strings into map?

					Content
						element / string data / character entity reference / entity reference / CDATA / PI?

						unicode text string
							[2] length (chars)
							... UTF-16 little-endian string without an end-of-string character

						string data
							value test / substitution

							value text
							05|45 value token
							01    value type
							...   value data

							substitution
								normal / optional

							[4] normal
								0x0d normal subs token
								[2] identifier
								[1] value type

							[4] optional
								0x0d normal subs token
								[2] identifier
								[1] value type

	*/

	return &record
}

func (r *Record) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Id:   %d\n", r.Id))
	sb.WriteString(fmt.Sprintf("Time: %s\n", r.Time.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Size: %d [0x%04x]\n", r.Size, r.Size))
	sb.WriteString(fmt.Sprintf("Copy: %d [0x%04x]\n\n", r.Copy, r.Copy))
	sb.WriteString(hex.Dump(r.Event))

	eventId := binary.LittleEndian.Uint16(r.Event[0x66:0x68])
	sb.WriteString(fmt.Sprintf("\nEventID: %d [0x%04x] DEBUG\n", eventId, eventId))

	return sb.String()
}

func (r *Record) IsValid() bool {
	return r.Size == r.Copy && !r.Time.IsZero()
}
