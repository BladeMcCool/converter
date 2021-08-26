package formats

import "net/url"

type FormatSpec struct {
	Name               string
	Mimetype           string
	RequiresDelimiters bool
	Implementation     func() FormatImporterExporter
}

type FormatImporterExporter interface{
	SetDelimiters(fieldPrefix string, form *url.Values) error
	Import([]byte) ([]Segment, error)
	Export([]Segment) []byte
}

type Segment struct {
	kind string
	values []string
}

var FormatSpecs = []FormatSpec{{
	Name: "Format1",
	Mimetype: "text/plain",
	RequiresDelimiters: true,
	Implementation: func()FormatImporterExporter { return &Format1ImporterExporter{} },
},{
	Name: "Format2",
	Mimetype: "application/json",
	Implementation: func()FormatImporterExporter { return &Format2ImporterExporter{} },
},{
	Name: "Format3",
	Mimetype: "application/xml",
	Implementation: func()FormatImporterExporter { return &Format3ImporterExporter{} },
}}


