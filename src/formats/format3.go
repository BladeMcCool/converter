package formats

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/beevik/etree"
	"net/url"
	"strconv"
	"strings"
)

type Format3ImporterExporter struct {
}
func (f *Format3ImporterExporter) SetDelimiters(_ string, _ *url.Values) (err error) {return}

func (f *Format3ImporterExporter) Import(rawData []byte) (segments []Segment, err error) {
	var hadValidSegment = false

	tokenparser := xml.NewDecoder(bytes.NewReader(rawData)) //hrm, maybe would be better to keep it as a reader in the first place.
	depth := 0
	var curSegment Segment
	for {
		token, err := tokenparser.Token()
		if token == nil {
			if err != nil && err.Error() != "EOF" {
				return []Segment{}, err
			}
			break
		}
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

func (f *Format3ImporterExporter) Export(segments []Segment) []byte {
	// during research for how to build the document the way i want i encountered a module called etree
	// using it here due to it being much easier to work with
	// the import code above can possibly be rewritten to use it too.
	outBuf := bytes.NewBuffer([]byte{})
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8" `)
	rootEl := doc.CreateElement("root")
	for _, segment := range segments {
		segmentEl := rootEl.CreateElement(segment.kind)
		for i, v := range segment.values {
			valueEl := segmentEl.CreateElement(segment.kind + strconv.Itoa(i+1))
			valueEl.SetText(v)
		}
	}
	doc.Indent(0)
	doc.WriteTo(outBuf)
	return outBuf.Bytes()
}
