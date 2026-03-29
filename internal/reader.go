package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

func ReadTime(r io.Reader) time.Time {
	nsec := int64(ReadUint64(r))
	nsec -= 116444736000000000
	nsec *= 100

	return time.Unix(0, nsec)
}

func ReadUint64(r io.Reader) uint64 {
	return binary.LittleEndian.Uint64(ReadBytes(r, 8))
}

func ReadUint32(r io.Reader) uint32 {
	return binary.LittleEndian.Uint32(ReadBytes(r, 4))
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
		case err != nil:
			panic(err)
		}
	}

	return true
}
