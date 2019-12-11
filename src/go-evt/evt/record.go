package evt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"time"

	"go-evt/sid"
)

// Record describes one evt record
type Record struct {
	Number        uint32
	CreationTime  time.Time
	WrittenTime   time.Time
	Identifier    EventIdentifier
	Type          EventType
	EventCategory uint16
	SourceName    string
	ComputerName  string
	SID           string `json:",omitempty"`
	Strings       string
	Data          []byte `json:",omitempty"`

	size          uint32
	stringsNumber uint16
	stringsOffset uint32
	sidSize       uint32
	sidOffset     uint32
	dataSize      uint32
	dataOffset    uint32
}

type Records []Record

type EventIdentifier struct {
	Code     uint16
	Facility uint16
	Customer CustomerFlag
	Severity Severity
}

type CustomerFlag uint8
type Severity uint8
type EventType uint16

const (
	SystemCode   CustomerFlag = 0x0
	CustomerCode CustomerFlag = 0x1

	SeveritySuccess       Severity = 0x0
	SeverityInformational Severity = 0x1
	SeverityWarning       Severity = 0x2
	SeverityError         Severity = 0x3

	TypeError         EventType = 0x1
	TypeWarning       EventType = 0x2
	TypeInformational EventType = 0x4
	TypeAuditSuccess  EventType = 0x8
	TypeAuditFailure  EventType = 0x10
)

func getRecord(r io.ReaderAt, recordOffset int64) (Record, error) {
	record := Record{}

	// read record size
	b := make([]byte, 4)
	_, err := r.ReadAt(b, recordOffset)
	if err != nil {
		return record, err
	}
	err = read(b, 0, 4, &record.size)
	if err != nil {
		return record, err
	}

	// read whole record to []byte
	b = make([]byte, record.size)
	_, err = r.ReadAt(b, recordOffset)
	if err != nil {
		return record, err
	}

	err = read(b, 8, 4, &record.Number)
	if err != nil {
		return record, err
	}

	t, err := getTime(b, 12)
	if err != nil {
		return record, err
	}
	record.CreationTime = t

	t, err = getTime(b, 16)
	if err != nil {
		return record, err
	}
	record.WrittenTime = t

	err = read(b, 20, 2, &record.Identifier.Code)
	if err != nil {
		return record, err
	}

	var d uint16
	err = binary.Read(bytes.NewReader(b[22:24]), binary.BigEndian, &d)
	if err != nil {
		return record, err
	}
	// facility is high 12 bit of d
	record.Identifier.Facility = d >> 4
	record.Identifier.Customer = CustomerFlag((d << 13) >> 15)
	record.Identifier.Severity = Severity((d << 14) >> 14)

	err = read(b, 24, 2, &record.Type)
	if err != nil {
		return record, err
	}

	err = read(b, 26, 2, &record.stringsNumber)
	if err != nil {
		return record, err
	}

	err = read(b, 28, 2, &record.EventCategory)
	if err != nil {
		return record, err
	}

	err = read(b, 36, 4, &record.stringsOffset)
	if err != nil {
		return record, err
	}

	err = read(b, 40, 4, &record.sidSize)
	if err != nil {
		return record, err
	}

	err = read(b, 44, 4, &record.sidOffset)
	if err != nil {
		return record, err
	}

	err = read(b, 48, 4, &record.dataSize)
	if err != nil {
		return record, err
	}

	err = read(b, 52, 4, &record.dataOffset)
	if err != nil {
		return record, err
	}

	s, l, err := getString(b, 56)
	if err != nil {
		return record, err
	}
	record.SourceName = s

	offset := 56 + l
	s, _, err = getString(b, offset)
	if err != nil {
		return record, err
	}
	record.ComputerName = s

	if record.sidOffset != 0 && record.sidSize != 0 {
		record.SID, err = sid.ParseSID(b[record.sidOffset : record.sidOffset+record.sidSize])
		if err != nil {
			return record, err
		}
	}

	if record.stringsOffset != 0 {
		offset = int64(record.stringsOffset)
		for i := 0; i < int(record.stringsNumber); i++ {
			s, l, err := getString(b, offset)
			if err != nil {
				return record, err
			}
			record.Strings += s + "\n"
			offset += l
		}
	}

	if record.dataOffset != 0 && record.dataSize != 0 {
		record.Data = b[record.dataOffset : record.dataOffset+record.dataSize]
	}

	return record, nil
}

func (record Record) String() string {
	b, err := json.Marshal(record)
	if err != nil {
		return ""
	}
	return string(b)
}
