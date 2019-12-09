package sid

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type SID struct {
	Bytes []byte

	revision       uint8
	subAuthCount   uint8
	authority      uint64
	subAuthorities []uint32
}

func ParseSID(b []byte) (string, error) {
	sid, err := NewSID(b)
	if err != nil {
		return "", err
	}
	return sid.String(), nil
}

func NewSID(b []byte) (SID, error) {
	errwf := errors.New("wrong SID format")
	if len(b) < 8 {
		return SID{}, errwf
	}
	sid := SID{
		revision:       b[0],
		subAuthCount:   b[1],
		authority:      binary.BigEndian.Uint64([]byte{0, 0, b[2], b[3], b[4], b[5], b[6], b[7]}),
		subAuthorities: make([]uint32, 0),
	}
	if len(b) < int(8+sid.subAuthCount*4) {
		return sid, errwf
	}

	for i := 0; i < int(sid.subAuthCount); i++ {
		sid.subAuthorities = append(sid.subAuthorities, binary.LittleEndian.Uint32(b[8+i*4:8+i*4+4]))
	}

	sid.Bytes = b

	return sid, nil
}

func (sid SID) String() string {
	if len(sid.Bytes) == 0 {
		return ""
	}
	s := fmt.Sprintf("S-%d-%d", sid.revision, sid.authority)
	for _, b := range sid.subAuthorities {
		s += fmt.Sprintf("-%d", b)
	}

	return s
}
