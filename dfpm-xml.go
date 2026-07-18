package main

import (
	"encoding/xml"
	"io"
	"strings"
)

type html struct {
	XMLName xml.Name `xml:"html"`
	Body    htmlbody `xml:"body"`
}

type htmlbody struct {
	XMLName xml.Name   `xml:"body"`
	Main    htmlmain   `xml:"main"`
	Script  htmlscript `xml:"script"`
}

type htmlscript struct {
	XMLName xml.Name `xml:"script"`
	Script  string   `xml:",innerxml"`
}

type htmlmain struct {
	XMLName  xml.Name      `xml:"main"`
	Sections []htmlsection `xml:"section"`
}

type htmlsection struct {
	XMLName xml.Name  `xml:"section"`
	Class   string    `xml:"class,attr"`
	Id      string    `xml:"id,attr"`
	Divs    []htmldiv `xml:"div"`
}

type htmldiv struct {
	XMLName  xml.Name      `xml:"div"`
	Class    string        `xml:"class,attr"`
	Articles []htmlarticle `xml:"article"`
}

type htmlarticle struct {
	XMLName xml.Name    `xml:"article"`
	Id      string      `xml:"id,attr"`
	Search  string      `xml:"data-search,attr"`
	Details htmldetails `xml:"details"`
}

type htmldetails struct {
	XMLName xml.Name    `xml:"details"`
	Summary htmlsummary `xml:"summary"`
}

type htmlsummary struct {
	XMLName xml.Name `xml:"summary"`
	Link    htmla    `xml:"a"`
}

type htmla struct {
	Href    string `xml:"href,attr"`
	Version string `xml:",innerxml"`
}

func ParseHtmlXml(r io.Reader, search string) (packs []string) {
	// f, err := os.Open(infile)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }

	fd := xml.NewDecoder(r)

	start := SkipToData(fd)
	var h htmlmain
	err := fd.DecodeElement(&h, &start)
	// var h html
	// err := fd.Decode(&h)
	if err != nil {
		println(err.Error())
		return
	}

	// for _, s := range h.Body.Main.Sections {
	// 	if s.Class == "pack-section" && s.Id == "microchip-dfp" {
	// 		println(s.Class, s.Id)
	// 		for _, a := range s.Divs[1].Articles {
	// 			if search != "" {
	// 				if strings.Contains(strings.ToLower(a.Search), search) {
	// 					println("*", a.Details.Summary.Link.Href, a.Details.Summary.Link.Version)
	// 				}
	// 			} else {
	// 				println(" ", a.Details.Summary.Link.Href, a.Details.Summary.Link.Version)
	// 			}
	// 		}
	// 	}
	// }

	for _, s := range h.Sections {
		if s.Class == "pack-section" {
			// println(s.Class, s.Id)
			for _, a := range s.Divs[1].Articles {
				if search != "" {
					if strings.Contains(strings.ToLower(a.Search), strings.ToLower(search)) {
						// println("*", a.Details.Summary.Link.Href, a.Details.Summary.Link.Version)
						packs = append(packs, a.Details.Summary.Link.Href)
					}
				} else {
					// println(" ", a.Details.Summary.Link.Href, a.Details.Summary.Link.Version)
					packs = append(packs, a.Details.Summary.Link.Href)
				}
			}
		}
	}

	return
}

func SkipToData(d *xml.Decoder) xml.StartElement {
	for {
		t, err := d.Token()
		if err != nil {
			return xml.StartElement{}
		}
		switch v := t.(type) {
		case xml.StartElement:
			// println("start", v.Name.Space, v.Name.Local)
			if v.Name.Local == "main" {
				return v.Copy()
			}
		case xml.EndElement:
			// println("end", v.Name.Space, v.Name.Local)
			// if v.Name.Local == "header" {
			// 	return
			// }
			// case xml.CharData:
			// 	println("char", string(v))
			// case xml.Comment:
			// 	println("comment", string(v))
			// case xml.ProcInst:
			// 	println("procinst", v.Target)
			// case xml.Directive:
			// 	println("directive", string(v))
		}
	}
}
