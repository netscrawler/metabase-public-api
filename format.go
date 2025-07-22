package metabase

// Format defines allowed export formats
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
	FormatXLSX Format = "xlsx"
)

func (f Format) Valid() bool {
	switch f {
	case FormatJSON, FormatCSV, FormatXLSX:
		return true
	default:
		return false
	}
}
