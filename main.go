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

const Slack = 28
const Chunk = 65536

func main() {
	var x uint64
	var y uint32
	var b []byte

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
			os.Exit(1)
		}
	}()

	for r := bufio.NewReaderSize(os.Stdin, Chunk); ReadUntil(r); {
		if y = ReadUint32(r); y < Slack || y > Chunk {
			continue // check sane size
		}

		if x = ReadUint64(r); x == 0 {
			continue // check valid record id
		}

		if x = ReadUint64(r); x == 0 {
			continue // check valid time
		}

		if b = ReadBytes(r, y-Slack); len(b) < 18 {
			continue // check valid stream length
		}

		if y != ReadUint32(r) {
			continue // check size equals copy
		}

		if y = binary.LittleEndian.Uint32(b[14:]); y > 20 {
			continue // check substitution items
		}

		if b[28] != 0x06 || b[29] != 0 {
			continue // check event id type and null
		}

		fmt.Printf("Record Time: %s Event ID: %d\n",
			FileTime(x).Format("2006-01-02 15:04:05.0000000Z"),
			binary.LittleEndian.Uint16(b[(y*4)+22:]),
		)
	}
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-116444736000000000)*100).UTC()
}

func ReadUint64(r *bufio.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, 8))
}

func ReadUint32(r *bufio.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, 4))
}

func ReadBytes(r *bufio.Reader, n uint32) []byte {
	if b, err := r.Peek(int(n)); err == nil {
		_, _ = r.Discard(int(n))
		return b
	} else {
		panic(err)
	}
}

func ReadUntil(r *bufio.Reader) bool {
	for {
		switch b, err := r.Peek(4); {
		default:
			_, _ = r.Discard(1)
		case string(b) == "**\x00\x00":
			_, _ = r.Discard(4)
			return true
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}
}
