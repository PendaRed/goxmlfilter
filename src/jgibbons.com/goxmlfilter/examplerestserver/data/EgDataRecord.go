package data

import (
	"encoding/json"
	"encoding/xml"
)

type EgDataRecord struct {
	Forename string `xml:"forename" json:"forename"`
	Surname  string `xml:"surname" json:"surname"`
	Gender   string `xml:"gender" json:"gender"`
	House    string `xml:"house" json:"house"`
}

func (r *EgDataRecord) AsJson() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r *EgDataRecord) AsXml() string {
	b, err := xml.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
