package evt

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	HeaderSize = 48
)

// Header describes .evt file
type Header struct {
	Version Version

	endOfFileRecordOffset uint32
}

// Version shows version of .evt file format
type Version struct {
	Major uint32
	Minor uint32
}

func (h Header) String() string {
	b, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(b)
}

func ParseEvt(file *os.File) (Header, Records, error) {
	h, err := getHeader(file)
	if err != nil {
		return Header{}, nil, err
	}

	records := Records{}

	offset := uint32(HeaderSize)

	for offset < h.endOfFileRecordOffset {
		record, err := getRecord(file, int64(offset))
		if record.size == 0 {
			return h, records, fmt.Errorf("can't parse record: %v", err)
		}
		offset += record.size
		if err != nil {
			fmt.Printf("can't parse record: %v. Continue...", err)
			continue
		}
		records = append(records, record)
	}

	return h, records, nil
}

func getHeader(file *os.File) (Header, error) {
	h := Header{}

	b := make([]byte, HeaderSize)
	_, err := file.ReadAt(b, 0)
	if err != nil {
		return h, err
	}

	err = read(b, 8, 4, &h.Version.Major)
	if err != nil {
		return h, err
	}

	err = read(b, 12, 4, &h.Version.Minor)
	if err != nil {
		return h, err
	}

	err = read(b, 20, 4, &h.endOfFileRecordOffset)
	if err != nil {
		return h, err
	}

	return h, nil
}
