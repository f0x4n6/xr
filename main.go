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
var buf4 = make([]byte, 4)
var buf8 = make([]byte, 8)

func main() {
	var b []byte
	var l, m uint32
	var n, t uint64

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
			os.Exit(1)
		}
	}()

	for r := bufio.NewReaderSize(os.Stdin, Chunk); ReadUntil(r); {
		if l = ReadUint32(r); l < Slack || l > Chunk {
			continue // check sane size
		}

		if n = ReadUint64(r); n == 0 {
			continue // check valid record id
		}

		if t = ReadUint64(r); t == 0 {
			continue // check valid time
		}

		if b = ReadBytes(r, buf0[:l-Slack]); len(b) < 18 {
			continue // check valid binxml length
		}

		if l != ReadUint32(r) {
			continue // check size equals copy
		}

		if m = binary.LittleEndian.Uint32(b[14:]); m > 20 {
			continue // check substitution array length
		}

		if b[28] != 0x06 || b[29] != 0 {
			continue
		}

		fmt.Printf("XR|%s|%d\n", FileTime(t).Format(time.RFC3339Nano), binary.LittleEndian.Uint16(b[(m*4)+22:]))
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

func ReadBytes(r io.Reader, b []byte) []byte {
	if n, err := r.Read(b); err == nil {
		return b[:n]
	} else {
		panic(err)
	}
}

func ReadUntil(r io.Reader) bool {
	for binary.LittleEndian.Uint32(buf4) != 0x00002A2A {
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
