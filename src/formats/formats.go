package formats

type FormatSpec struct {
	Name               string
	Mimetype           string
	RequiresDelimiters bool
	Implementation     FormatImporterExporter
}

type FormatImporterExporter interface{
	SetDelimiters(string, string) //this naively assumes any format that uses delimiters requires exactly two of them.
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
	Implementation: &Format1ImporterExporter{},
},{
	Name: "Format2",
	Mimetype: "application/json",
	Implementation: &Format2ImporterExporter{},
},{
	Name: "Format3",
	Mimetype: "application/xml",
	Implementation: &Format3ImporterExporter{},
}}


