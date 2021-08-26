package formats

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type jsonFormat map[string]map[string]string

type Format2ImporterExporter struct {
}
func (f *Format2ImporterExporter) SetDelimiters(_ string, _ *url.Values) (err error) {return}

func (f *Format2ImporterExporter) Import(rawData []byte) (segments []Segment, err error) {
	jsonSegments := &jsonFormat{}
	err = json.Unmarshal(rawData, jsonSegments)

	if err != nil {
		return []Segment{}, err
	}

	var hadValidSegment = false
	for segmentKind, segmentValues := range *jsonSegments {
		curSegment := Segment{
			kind: segmentKind,
			values: []string{},
		}
		moreValues := true
		checkEntry := 1
		for moreValues == true {
			checkKey := segmentKind + strconv.Itoa(checkEntry)
			val, ok := segmentValues[checkKey]
			if !ok {
				break
			}
			checkEntry++
			curSegment.values = append(curSegment.values, strings.Trim(val, " "))
		}
		if len(curSegment.values) == 0 {
			continue
		}

		segments = append(segments, curSegment)
		hadValidSegment = true
	}

	if !hadValidSegment { return []Segment{}, errors.New("invalid input: no segments with at least one value were found")}

	return
}

func (f *Format2ImporterExporter) Export(segments []Segment) (jsonExport []byte) {
	segmentsExport := jsonFormat{}
	for _, segment := range segments {
		segmentsExport[segment.kind] = map[string]string{}
		for i, v := range segment.values {
			segmentsExport[segment.kind][segment.kind + strconv.Itoa(i+1)] = v
		}
	}
	//a JSON object -- An object is an unordered set of name/value pairs
	// if order preservation is required then a custom marshaller, or defined structs, could be possible.
	jsonExport, err := json.Marshal(segmentsExport)
	if err != nil { return []byte{} }
	return
}
