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
const Magic = "**\x00\x00"

var buf0 = make([]byte, 0, Chunk)
var buf4 = make([]byte, 4)
var buf8 = make([]byte, 8)

func main() {
	var x uint64
	var y uint32
	var b []byte

	/*
		defer func() {
			if err := recover(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "xr:", err)
				os.Exit(1)
			}
		}()
	*/

	for r := bufio.NewReaderSize(os.Stdin, Chunk); ReadUntil(r); {
		if y = ReadUint32(r); y < Slack || y > Chunk {
			println("abort: size", y)
			continue // check sane size
		}

		if x = ReadUint64(r); x == 0 {
			println("abort: record id", x)
			continue // check valid record id
		}

		if x = ReadUint64(r); x == 0 {
			println("abort: time", x)
			continue // check valid time
		}

		if b = ReadBytes(r, buf0[:y-Slack]); len(b) < 18 {
			println("abort: stream", len(b))
			continue // check valid stream length
		}

		if y != ReadUint32(r) {
			println("abort: copy", y)
			continue // check size equals copy
		}

		if y = binary.LittleEndian.Uint32(b[14:]); y > 20 {
			println("abort: items", y)
			continue // check substitution items
		}

		if b[28] != 0x06 || b[29] != 0 {
			println("abort: event id")
			continue // check event id type and null
		}

		fmt.Printf("Record Time: %s EventID: %d\n", FileTime(x).Format("2006-01-02 15:04:05.0000000Z"), binary.LittleEndian.Uint16(b[(y*4)+22:]))
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

func ReadUntil(r *bufio.Reader) bool {
	for {
		switch b, err := r.Peek(4); {
		default:
			_, _ = r.Discard(1)
		case string(b) == Magic:
			_, _ = r.Discard(4)
			return true
		case errors.Is(err, io.EOF):
			return false
		case err != nil:
			panic(err)
		}
	}
}
