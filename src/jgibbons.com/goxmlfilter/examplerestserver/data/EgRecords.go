package data

import (
	"encoding/json"
	"encoding/xml"
)

type EgRecords struct {
	People []*EgDataRecord `xml:"people" json:"people"`
}

func (r *EgRecords) Add(rec *EgDataRecord) {
	r.People = append(r.People, rec)
}

func (r *EgRecords) AsJson() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r *EgRecords) AsXml() string {
	b, err := xml.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
