package utils

import (
	"encoding/xml"
	"io"
)

type xmlResult struct {
	m map[string]string
}

type xmlMapEntry struct {
	XMLName xml.Name
	XMLAttr xml.Attr
	Value   string `xml:",chardata"`
}

// NewXMLResult ...
func NewXMLResult() xmlResult {
	return xmlResult{
		m: map[string]string{},
	}
}

// Get ...
func (xr xmlResult) Get(k string) string {
	if v, ok := xr.m[k]; ok {
		return v
	}
	return ""
}

// Set ...
func (xr xmlResult) SetKey(k string, r string) {
	if v, ok := xr.m[k]; ok {
		xr.m[r] = v
		delete(xr.m, k)
	}
}

// GetAll ...
func (xr xmlResult) GetAll() map[string]string {
	return xr.m

}

// UnmarshalXML ...
func (xr *xmlResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		xr.m[e.XMLName.Local] = e.Value
	}
	return nil
}


func (xr *xmlResult) UnmarshalXMLOriginal(d *xml.Decoder, start xml.StartElement) error {
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		xr.m[e.XMLName.Local] = e.Value
	}
	return nil
}