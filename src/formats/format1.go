package formats

import (
	"errors"
	"net/url"
	"strings"
)

type Format1ImporterExporter struct {
	recordSeparator string
	valueSeparator string
}

func (f *Format1ImporterExporter) SetDelimiters(fieldPrefix string, form *url.Values) (err error) {
	lineSep := form.Get(fieldPrefix + "LineSeparator")
	elSep := form.Get(fieldPrefix + "ElementSeparator")
	//assuming these cannot be blank...
	if lineSep == "" || elSep == "" {
		return errors.New("missing one or more separator")
	}
	f.recordSeparator = lineSep
	f.valueSeparator = elSep
	return
}



func (f *Format1ImporterExporter) Import(rawData []byte) (segments []Segment, err error) {
	// given the statements that the 3 different sample inputs formats are equivalent to each other
	// and given that there are actually inconsistencies in the leading/trailing spaces
	//  - txt vs json/xml for AddressID eg: " 14 " vs " 14"
	//  - txt vs json/xml for ContactId eg: " 59" vs "59"
	// so we will trim the values in order for them to in fact be equivalent.
	//
	// if export should use different spaces padding for certain fields depending on the
	// export type, then that logic will need to be defined clearly before we can accommodate it.

	rawSegments := strings.Split(string(rawData), f.recordSeparator)
	var hadValidSegment = false
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

	return
}

func (f *Format1ImporterExporter) Export(segments []Segment) []byte {
	segmentsExport := []string{}
	for _, segment := range segments {
		exportLine := segment.kind + f.valueSeparator + strings.Join(segment.values, f.valueSeparator)
		segmentsExport = append(segmentsExport, exportLine)
	}
	return []byte(strings.Join(segmentsExport, f.recordSeparator) + f.recordSeparator)
}
