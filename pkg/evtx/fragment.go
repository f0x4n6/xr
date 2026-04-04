package evtx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"go.foxforensics.dev/tri/pkg/utils"
)

const HeaderSize = 14
const TemplateOffset1 = 6
const TemplateOffset2 = 18
const TypeUInt16 = 0x06
const TypeSID = 0x13

type Fragment struct {
	TemplateId1 uint32
	TemplateId2 uint32
	Computer    string
	EventId     uint16
	UserId      []byte
	Datas       []Item
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

	// set computer by cached template
	if v, ok := Computers[fragment.TemplateId1]; ok {
		fragment.Computer = v
	} else {
		fragment.Computer = LastComputer + "?"

		if Debug {
			utils.Debug("[!] No entry for template [%08x]\n", fragment.TemplateId1)
		}
	}

	// skip fragment header (unused)
	r := bytes.NewReader(stream[HeaderSize:])

	fragment.Items = make([]Item, utils.ReadUint32(r))

	// invalid substitution array length
	if !fragment.IsItemsValid() {
		return fragment
	}

	for i := 0; i < len(fragment.Items); i++ {
		_ = binary.Read(r, binary.LittleEndian, &fragment.Items[i])
	}

	if len(fragment.Items) > 2 && fragment.Items[2].Type == TypeUInt16 {
		fragment.EventId = binary.LittleEndian.Uint16(fragment.GetItemData(2))
	}

	/*
		if len(fragment.Items) > 12 && fragment.Items[12].Type == TypeSID {
			a := fragment.GetItemData(12)
			var v uint64

			for _, b := range a[2:8] {
				v = (v << 8) | uint64(b)
			}

			// "S-1-5-19"
			// 01 01 00 00 00 13
			// 01 02 00 00 00 00 00 05 13 00 00 00 (?)
			// ??  0f 01 01 00 0c 01 59 41 b6 26 4d 76
			//utils.Debug("SID %d S-%d-%d [%x]\n", a[1], a[0], v, a)

				case 0x13: // SID
					str := "S"
					str += fmt.Sprintf("-%d", ctx.ConsumeUint8())

					ctx.ConsumeUint8()
					v_q := uint64(0)
					for _, b := range ctx.ConsumeBytes(6) {
					v_q = (v_q << 8) | uint64(b)
					}

					str += fmt.Sprintf("-%d", v_q)
					for idx := 0; idx < arg.argLen-8; idx += 4 {
					str += fmt.Sprintf("-%d", ctx.ConsumeUint32())
					}
					arg_values[idx] = str

			// fragment.UserId = 0 //binary.LittleEndian.Uint16(fragment.GetItemData(12))
		}
	*/

	if !fragment.IsTemplate() {
		return fragment
	}

	fmt.Printf("%x\n", stream[:16])

	DataOffset := binary.LittleEndian.Uint16(stream[0x0A:0x0C])
	DataSize := binary.LittleEndian.Uint32(stream[0x22:0x26])

	fmt.Printf("DataOffset %08x %d\n", DataOffset, DataOffset)
	fmt.Printf("DataSize   %08x %d\n", DataSize, DataSize)

	fragment.Datas = make([]Item, DataSize)

	r2 := bytes.NewReader(stream[DataOffset:])

	for i := 0; i < len(fragment.Datas); i++ {
		_ = binary.Read(r2, binary.LittleEndian, &fragment.Datas[i])
	}

	// 000003b0  00 00 00 00 00 00 00 00  00 00 00 00 00 f3 00 21  |...............!|
	// 000003c0  00 04 00 00 00[7b 17] 0  80 00 00 00 00 00 00 80  |.....{..........|

	// TODO: Search by NameString hash? 3b 6e

	// 000002e0  04 00 00|3b 6e|08 00|43  00 6f 00 6d 00 70 00 75  |...;n..C.o.m.p.u|
	// 000002f0  00 74 00 65 00 72 00|00  00|02|05|01|0f 00|57 00  |.t.e.r........W.|
	// 00000300  49 00 4e 00 2d 00 54 00  49 00 50 00 33 00 4e 00  |I.N.-.T.I.P.3.N.|
	// 00000310  39 00 30 00 4b 00 4b 00  37 00 34 00|04|41 ff ff  |9.0.K.K.7.4..A..|

	// record is a template instance and carries a computer name to be cached
	if i := bytes.Index(fragment.Stream, utils.ToUtf16("Computer")); i >= 0 {
		if j := bytes.Index(fragment.Stream[i:], []byte{0x02}); j >= 0 {
			l := int(binary.LittleEndian.Uint16(fragment.Stream[i+j+3:i+j+5]) * 2)
			LastComputer = utils.FromUtf16(fragment.Stream[i+j+5 : i+j+5+l])
			Computers[fragment.TemplateId1] = LastComputer
			fragment.Computer = LastComputer
		}
	}

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

	if len(f.Stream) <= int(offset+f.Items[n].Size) {
		return []byte("err")
	}

	return f.Stream[offset : offset+f.Items[n].Size]
}

func (f *Fragment) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[+] TemplateID1  [%08x]\n", f.TemplateId1))
	sb.WriteString(fmt.Sprintf("[+] TemplateID2  [%08x]\n", f.TemplateId2))
	sb.WriteString(fmt.Sprintf("[+] Computer     %s\n", f.Computer))
	sb.WriteString(fmt.Sprintf("[+] EventID      %d\n", f.EventId))
	sb.WriteString(fmt.Sprintf("[+] UserID       %s\n", f.UserId))
	sb.WriteString(fmt.Sprintf("[+] Datas        %d\n", len(f.Datas)))

	for i, v := range f.Datas {
		sb.WriteString(fmt.Sprintf("[+]  #%02d %04x %02x %02x = %x\n", i+1, v.Size, v.Type, v.Null, 0))
	}

	sb.WriteString(fmt.Sprintf("[+] Items        %d\n", len(f.Items)))

	for i, v := range f.Items {
		sb.WriteString(fmt.Sprintf("[+]  #%02d %04x %02x %02x = %x\n", i+1, v.Size, v.Type, v.Null, f.GetItemData(i)))
	}

	return sb.String()
}
