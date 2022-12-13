package metadata

import (
	"github.com/siongui/gojianfan"
)

type Artist struct {
	Name     string
	SortName string
	Aliases  []string
}

func (this *Artist) ToSimplified() {
	this.Name = gojianfan.T2S(this.Name)
	for i := range this.Aliases {
		this.Aliases[i] = gojianfan.T2S(this.Aliases[i])
	}
}
