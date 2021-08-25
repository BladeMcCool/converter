package formats

import (
	"bytes"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

type Format3ImporterExporter struct {
}
func (f *Format3ImporterExporter) SetDelimiters(_, _ string) {}

func (f *Format3ImporterExporter) Import(rawData []byte) (segments []Segment, err error) {
	var hadValidSegment = false

	tokenparser := xml.NewDecoder(bytes.NewReader(rawData)) //hrm, maybe would be better to keep it as a reader in the first place.

	//reachedRoot := false
	depth := 0
	var curSegment Segment
	for {
		token, _ := tokenparser.Token()
		if token == nil { break }
		switch element := token.(type) {
		case xml.StartElement:
			depth++
			//log.Printf("starttag depth: %d", depth)
			if depth == 2 {
				curSegment = Segment{
					kind: element.Name.Local,
					values: []string{},
				}
			}
			if depth == 3 && element.Name.Local == curSegment.kind + strconv.Itoa(len(curSegment.values)+1) {
				//we would basically stop extracting if the expected record sequence was broken though.
				//thats a limitation here.
				token, err := tokenparser.Token()
				if err != nil {
					return []Segment{}, err
				}
				textValue, ok := token.(xml.CharData)
				if !ok {
					continue
				}
				curSegment.values = append(curSegment.values, strings.Trim(string(textValue), " "))
				hadValidSegment = true
				//log.Printf("%# v\n", pretty.Formatter(element))
			}

		case xml.EndElement:
			if depth == 2 && element.Name.Local == curSegment.kind {
				//log.Printf("endtag depth: %d", depth)
				//log.Printf("%# v\n", pretty.Formatter(element))
				segments = append(segments, curSegment)
			}
			depth--
		}

	}

	if !hadValidSegment { return []Segment{}, errors.New("invalid input: no segments with at least one value were found")}
	return
}

func (f *Format3ImporterExporter) Export(segments []Segment) (xmlExport []byte) {
	var err error
	if err != nil { return []byte{} }
	return
}
