package evt

import (
	"strings"
	"time"
)

type FilterFunc func(record Record) bool

func (r Records) Filter(f FilterFunc) Records {
	re := Records{}
	for _, record := range r {
		if f(record) {
			re = append(re, record)
		}
	}
	return re
}

func (r Records) FilterText(text string) Records {
	return r.Filter(func(record Record) bool {
		return strings.Contains(record.SourceName, text) ||
			strings.Contains(record.ComputerName, text) ||
			strings.Contains(record.SID, text) ||
			strings.Contains(record.Strings, text)
	})
}

func (r Records) FilterAfterTime(t time.Time) Records {
	return r.Filter(func(record Record) bool {
		return record.CreationTime.After(t)
	})
}

func (r Records) FilterBeforeTime(t time.Time) Records {
	return r.Filter(func(record Record) bool {
		return record.CreationTime.Before(t)
	})
}

func (r Records) FilterType(t ...EventType) Records {
	return r.Filter(func(record Record) bool {
		match := false
		for _, tt := range t {
			match = match || (record.Type == tt)
		}
		return match
	})
}

func (r Records) FilterSeverity(min, max Severity) Records {
	return r.Filter(func(record Record) bool {
		return record.Identifier.Severity >= min && record.Identifier.Severity <= max
	})
}

func (r Records) FilterCodes(codes ...uint16) Records {
	return r.Filter(func(record Record) bool {
		match := false
		for _, code := range codes {
			match = match || (record.Identifier.Code == code)
		}
		return match
	})
}

func (r Records) FilterFacility(min, max uint16) Records {
	return r.Filter(func(record Record) bool {
		return record.Identifier.Facility >= min && record.Identifier.Facility <= max
	})
}
