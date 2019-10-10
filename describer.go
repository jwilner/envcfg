package envcfg

import "encoding/json"

func Describe(configurable Configurable) []Description {
	d := new(Describer)
	configurable.Configure(d)
	return d.descriptions
}

type Description struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Optional bool                   `json:"optional"`
	Default  *DefaultValDescription `json:"default,omitempty"`
	Params   interface{}            `json:"params,omitempty"`
	Comment  string                 `json:"comment,omitempty"`
}

type DefaultValDescription struct {
	Value interface{}
}

func (d DefaultValDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

type Describer struct {
	descriptions []Description
}

func (d *Describer) Has(name string) bool {
	d.addDescription(Description{
		Name:     name,
		Type:     "bool",
		Optional: true,
		Default:  &DefaultValDescription{false},
	})
	return false
}

func (d *Describer) HasNot(name string) bool {
	d.addDescription(Description{
		Name:     name,
		Type:     "bool",
		Optional: true,
		Default:  &DefaultValDescription{true},
	})
	return false
}

func (d *Describer) describe(s *spec) {
	desc := Description{
		Name:     s.name,
		Type:     s.typeName,
		Optional: s.flags&flagOptional > 0,
		Comment:  s.comment,
		Params:   s.parser.describe(),
	}
	if s.flags&flagDefaultVal > 0 {
		desc.Default = &DefaultValDescription{s.defaultVal}
	}
	d.addDescription(desc)
}

func (d *Describer) addDescription(desc Description) {
	d.descriptions = append(d.descriptions, desc)
}

func (d *Describer) Describe() []Description {
	return d.descriptions
}

var _ Configurer = new(Describer)
