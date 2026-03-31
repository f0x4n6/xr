package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const ChunkSize = 65536
const EpochZero = 116444736000000000

func Debug(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, a...)
}

func FileTime(t uint64) time.Time {
	return time.Unix(0, (int64(t)-EpochZero)*100)
}

func Unicode(s []byte) []byte {
	b := bytes.Repeat([]byte{0}, len(s)/2)

	for i := 0; i < len(s)/2; i++ {
		b[i] = s[i*2]
	}

	return b
}

func FromUtf16(s string) []byte {
	b := bytes.Repeat([]byte{0}, (len(s)*2)+4)

	binary.LittleEndian.PutUint16(b[0:2], uint16(len(s)))

	for i, c := range []byte(s) {
		b[2+(i*2)] = c
	}

	return b
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, 8))
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, 4))
}

func ReadUint16(r io.Reader) uint16 {
	return binary.LittleEndian.Uint16(ReadBytes(r, 2))
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
