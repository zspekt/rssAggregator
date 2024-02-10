package xmldecoding

import (
	"encoding/xml"
	"io"
)

func DecodeXml(r io.Reader, rss *Rss) error {
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(rss)
	if err != nil {
		return err
	}
	return nil
}
