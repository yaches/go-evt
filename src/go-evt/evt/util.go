package evt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
	"unicode/utf16"
)

func read(b []byte, offset, size int64, target interface{}) error {
	if offset+size > int64(len(b)) {
		return errors.New("out of range")
	}
	return binary.Read(bytes.NewReader(b[offset:offset+size]), binary.LittleEndian, target)
}

func getTime(b []byte, offset int64) (time.Time, error) {
	var epoch uint32
	err := read(b, offset, 4, &epoch)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(epoch), 0), nil
}

func getString(b []byte, start int64) (string, int64, error) {
	offset := start
	for {
		if offset >= int64(len(b)) {
			return "", 0, errors.New("out of range")
		}
		b1 := b[offset]
		b2 := b[offset+1]
		if b1 == 0x00 && b2 == 0x00 {
			break
		}
		offset += 2
	}

	ints := make([]uint16, (offset-start)/2)
	err := binary.Read(bytes.NewReader(b[start:offset]), binary.LittleEndian, &ints)
	if err != nil {
		return "", 0, err
	}

	return string(utf16.Decode(ints)), offset - start + 2, nil
}
