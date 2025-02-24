package filters

type FilterType int

const (
	TextFilter FilterType = iota
	DateFilter
	MultiSelectFilter
)

var filterTypeToString = map[FilterType]string{
	TextFilter:        "TextFilter",
	DateFilter:        "DateFilter",
	MultiSelectFilter: "MultiSelectFilter",
}

func (ft FilterType) String() string {
	s := ""
	if str, ok := filterTypeToString[ft]; ok {
		s = str
	} else {
		s = "Unknown"
	}
	return s
}

type Filter struct {
	Name  string
	State string
	Type  FilterType
}
