// Experimental record analyzer.
//
// Usage:
//
//	cat FILE | xr | uniq | sort
package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const Chunk = 65536
const Slack = 28

var buf0 = make([]byte, 0, Chunk)
var buf2 = make([]byte, 2)
var buf4 = make([]byte, 4)
var buf8 = make([]byte, 8)

func main() {
	var t uint64
	var n uint32
	var b []byte

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
			os.Exit(1)
		}
	}()

	for r := bufio.NewReaderSize(os.Stdin, Chunk); ReadUntil(r); {
		if n = ReadUint32(r); n < Slack || n > Chunk {
			continue // check sane size
		}

		if t = ReadUint64(r); t == 0 {
			continue // check valid record id
		}

		if t = ReadUint64(r); t == 0 {
			continue // check valid time
		}

		if b = ReadBytes(r, buf0[:n-Slack]); len(b) < 18 {
			continue // check valid binxml length
		}

		if n != ReadUint32(r) {
			continue // check size equals copy
		}

		if n = binary.LittleEndian.Uint32(b[14:]); n > 20 {
			continue // check substitution array length
		}

		if b[28] != 0x06 || b[29] != 0 {
			continue // check event id type and null
		}

		fmt.Printf("XR|%s|%d\n", FileTime(t).Format("2006-01-02 15:04:05.0000000Z"), binary.LittleEndian.Uint16(b[(n*4)+22:]))
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100).UTC()
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, buf8))
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, buf4))
}

func ReadUint16(r io.Reader) uint16 {
	return binary.LittleEndian.Uint16(ReadBytes(r, buf2))
}

func ReadBytes(r io.Reader, b []byte) []byte {
	if n, err := r.Read(b); err == nil {
		return b[:n]
	} else {
		panic(err)
	}
}

func ReadUntil(r io.Reader) bool {
	for string(buf2) != "**" {
		switch _, err := io.ReadFull(r, buf2); {
		case errors.Is(err, io.ErrUnexpectedEOF):
			return false
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}

	ReadUint16(r)
	return true
}
