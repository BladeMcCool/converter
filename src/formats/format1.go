package formats

import (
	"errors"
	"strings"
)

type Format1Formatter struct {
	recordSeparator string
	valueSeparator string
}


func (f *Format1Formatter) SetDelimiters(recordSeparator, valueSeparator string) {
	f.recordSeparator = recordSeparator
	f.valueSeparator = valueSeparator
}

func (f *Format1Formatter) Import(rawData []byte) ([]Segment, error) {
	//note: given the statements that the 3 different sample inputs are equivalent
	//      and given that there are actually inconsistencies in the leading/trailing spaces
	//       - txt vs json/xml for AddressID eg: " 14 " vs " 14"
	//       - txt vs json/xml for ContactId eg: " 59" vs "59"

	rawSegments := strings.Split(string(rawData), f.recordSeparator)
	var hadValidSegment = false
	segments := []Segment{}
	for _, rawSegment := range rawSegments {
		rawValues := strings.Split(rawSegment, f.valueSeparator)
		if len(rawValues) < 1 {
			break
		}
		curSegment := Segment{
			kind: rawValues[0],
			values: []string{},
		}

		for i := 1; i < len(rawValues); i++ {
			trimmed := strings.Trim(rawValues[i], " ")
			curSegment.values = append(curSegment.values, trimmed)
		}
		if len(curSegment.values) == 0 {
			break
		}

		segments = append(segments, curSegment)
		hadValidSegment = true
	}
	if !hadValidSegment { return []Segment{}, errors.New("invalid input: no segments with at least one value were found")}

	return segments, nil
}

func (f *Format1Formatter) Export(segments []Segment) []byte {
	segmentsExport := []string{}
	for _, segment := range segments {
		exportLine := segment.kind + f.valueSeparator + strings.Join(segment.values, f.valueSeparator)
		segmentsExport = append(segmentsExport, exportLine)
	}
	return []byte(strings.Join(segmentsExport, f.recordSeparator) + f.recordSeparator)
}
