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
	Id       uint64
	Time     uint64
	Size     uint32
	Copy     uint32
	Stream   []byte
	Fragment *Fragment
}

func NewRecord(r io.Reader) *Record {
	record := Record{
		Size: ReadUint32(r),
		Id:   ReadUint64(r),
		Time: ReadUint64(r),
	}

	record.Stream = ReadBytes(r, record.Size-4-4-4-8-8)
	record.Copy = ReadUint32(r)

	// invalid record size
	if !record.IsSizeValid() {
		return &record
	}

	// invalid stream length
	if !record.IsStreamValid() {
		return &record
	}

	record.Fragment = NewFragment(record.Stream)

	return &record
}

func (r *Record) IsSizeValid() bool {
	return r.Size == r.Copy && r.Size > HeaderSize && r.Size < ChunkSize
}

func (r *Record) IsTimeValid() bool {
	return r.Time > 0
}

func (r *Record) IsStreamValid() bool {
	return len(r.Stream) >= 22
}

func (r *Record) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Id    [%04x] %d\n", r.Id, r.Id))
	sb.WriteString(fmt.Sprintf("Size  [%04x] %d\n", r.Size, r.Size))
	sb.WriteString(fmt.Sprintf("Copy  [%04x] %d\n", r.Copy, r.Copy))
	sb.WriteString(fmt.Sprintf("Time  %s\n", FileTime(r.Time).Format(time.RFC3339)))
	sb.WriteString(r.Fragment.String())
	sb.WriteString(hex.Dump(r.Stream))

	return sb.String()
}
