package export

import (
	"strings"
	_ "time/tzdata"
)

const (
	tagTitle        = "title"
	tagType         = "type"
	typeRFC3339     = "RFC3339"
	typeArrayString = "ArrayString"
)

type ExportTags struct {
	Name string
	Type string
}

func (e *ExportTags) Convert(list string) {
	f := strings.Split(list, ",")
	for _, g := range f {
		h := strings.Split(g, ":")
		switch h[0] {
		case "title":
			e.Name = h[1]
		case "type":
			e.Type = h[1]
		}
	}
}
