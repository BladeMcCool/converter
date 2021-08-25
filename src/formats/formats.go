package formats

type formatSpec struct {
	name           string
	implementation FormatImporterExporter
}

type FormatImporterExporter interface{
	SetDelimiters(string, string) error
	Import([]byte) ([]Segment, error)
	Export([]Segment) []byte
}
//recordSeparator
//valueSeparator

type Segment struct {
	kind string
	values []string
}

//var formatSpecs = []formatSpec{
//	formatSpec{"format1"},
//}


func gogo() bool {
	return true
}