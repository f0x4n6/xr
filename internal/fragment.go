package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const HeaderSize = 14
const TemplateOffset1 = 6
const TemplateOffset2 = 18

type Fragment struct {
	TemplateId1 uint32
	TemplateId2 uint32
	Computer    []byte
	EventId     uint16
	Items       []Item
	Stream      []byte
}

type Item struct {
	Size uint16
	Type byte
	Null byte
}

func NewFragment(stream []byte) *Fragment {
	fragment := &Fragment{
		TemplateId1: binary.BigEndian.Uint32(stream[TemplateOffset1 : TemplateOffset1+4]),
		TemplateId2: binary.BigEndian.Uint32(stream[TemplateOffset2 : TemplateOffset2+4]),
		Stream:      stream,
	}

	// skip fragment header (unused)
	r := bytes.NewReader(stream[HeaderSize:])

	fragment.Items = make([]Item, ReadUint32(r))

	// invalid substitution array length
	if !fragment.IsItemsValid() {
		return fragment
	}

	for i := 0; i < len(fragment.Items); i++ {
		_ = binary.Read(r, binary.LittleEndian, &fragment.Items[i])
	}

	if len(fragment.Items) >= 3 {
		fragment.EventId = binary.LittleEndian.Uint16(fragment.GetItemData(2))
	}

	if !fragment.IsTemplate() {
		return fragment
	}

	computer := FromUtf16("Computer")

	Debug("%x\n", computer)

	if ReadUntil(r, computer) && ReadUntil(r, []byte{0x05, 0x01}) {
		length := ReadUint16(r)

		Debug("Length [%04x]", length)

		fragment.Computer = ReadBytes(r, uint32(length))
	}

	// 000002e0  04 00 00|3b 6e|08 00|43  00 6f 00 6d 00 70 00 75  |...;n..C.o.m.p.u|
	// 000002f0  00 74 00 65 00 72 00|00  00|02|05|01|0f 00|57 00  |.t.e.r........W.|
	// 00000300  49 00 4e 00 2d 00 54 00  49 00 50 00 33 00 4e 00  |I.N.-.T.I.P.3.N.|
	// 00000310  39 00 30 00 4b 00 4b 00  37 00 34 00|04|41 ff ff  |9.0.K.K.7.4..A..|

	return fragment
}

func (f *Fragment) IsItemsValid() bool {
	return len(f.Items) <= 20
}

func (f *Fragment) IsTemplate() bool {
	return f.TemplateId1 == f.TemplateId2
}

func (f *Fragment) GetItemData(n int) []byte {
	offset := uint16((len(f.Items) * 4) + HeaderSize + 4)

	for i := 0; i <= n; i++ {
		offset += f.Items[i].Size
	}

	return f.Stream[offset : offset+f.Items[n].Size]
}

func (f *Fragment) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("EventID   %d\n", f.EventId))
	sb.WriteString(fmt.Sprintf("Computer  %s\n", f.Computer))

	for i, v := range f.Items {
		sb.WriteString(fmt.Sprintf("Item #%02d  %04x %02x %02x\n", i, v.Size, v.Type, v.Null))
	}

	return sb.String()
}
