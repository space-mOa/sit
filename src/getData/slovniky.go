package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// WikiData je struct pro XML soubor
type WikiData struct {
	XMLName  xml.Name   `xml:"mediawiki"`
	XMLNS    string     `xml:"xmlns,attr"`
	XSI      string     `xml:"xsi,attr"`
	Location string     `xml:"schemaLocation,attr"`
	Version  string     `xml:"version,attr"`
	Lang     string     `xml:"lang,attr"`
	SiteInfo []SiteInfo `xml:"siteinfo"`
	Page     []Page     `xml:"page"`
}

// SiteInfo obsahuje informace o stránce a Namespace
type SiteInfo struct {
	XMLName    xml.Name    `xml:"siteinfo"`
	SiteName   string      `xml:"sitename"`
	DbName     string      `xml:"dbname"`
	Base       string      `xml:"base"`
	Generator  string      `xml:"generator"`
	Case       string      `xml:"case"`
	Namespaces []Namespace `xml:"namespaces"`
}

type Namespace struct {
	XMLName xml.Name    `xml:"namespaces"`
	Name    []NameSpace `xml:"namespace"`
}

type NameSpace struct {
	XMLName xml.Name `xml:"namespace"`
	Key     string   `xml:"key"`
	Case    string   `xml:"case"`
	Name    string   `xml:"namespace"`
}

// Page obsahuje např. Název a Redirect, který v sobě má text článku
type Page struct {
	XMLName  xml.Name `xml:"page"`
	Title    string   `xml:"title"`
	Ns       string   `xml:"ns"`
	Redirect Redirect `xml:"redirect"`
	Revision Revision `xml:"revision"`
}

// Redirect přesměrování na jiné heslo v Soc. encyklopedii
type Redirect struct {
	XMLName xml.Name `xml:"redirect"`
	Title   string   `xml:"title,attr"`
}

type Contributore struct {
	XMLName  xml.Name `xml:"contributor"`
	UserName string   `xml:"username"`
}

// Revision je text hesla
type Revision struct {
	XMLName     xml.Name     `xml:"revision"`
	Contributor Contributore `xml:"contributor"`
	Text        []byte       `xml:"text"`
}

func unpackFile(v *WikiData, name string) (err error) {
	c, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(c, &v)
	if err != nil {
		return err
	}

	return nil
}

/* Získá každý slovník/encyklopedii */

func (w *WikiData) getCategory(name string, head []string, r string) (ctg Node) {
	rgx := getRegexp(
		r,               // 0 reg
		`\[\[[^]]*\]\]`, // 1 Odkazy
		`Soubor:`,       // 2 Vybrání souboru pro jeho zahození
		`#REDIRECT`,     // 3 REDIRECT odstraní duplikovaná hesla
	)
	ctg.name = name
	ctg.head = head
	var index uint64 = 0
	for _, piece := range w.Page {
		if len(rgx[0].FindAll(piece.Revision.Text, -1)) == 0 || len(rgx[3].FindAll(piece.Revision.Text, -1)) != 0 {
			continue
		}
		index++
		fndLinks := rgx[1].FindAll(piece.Revision.Text, -1)
		var links []string
		for _, link := range fndLinks {
			if rgx[2].FindAllString(string(link), -1) == nil {
				links = append(links, trimLink(string(link), string(piece.Title)))
			}
		}
		ctg.addNode([]string{strconv.FormatUint(index, 10), string(piece.Title)}, links, piece.Revision.Text)
	}
	return ctg
}

func (w *WikiData) getCategoryModify(name string, head []string, r string, f func(piece Page) []string) (ctg Node) {
	rgx := getRegexp(
		r,               // 0 reg
		`\[\[[^]]*\]\]`, // 1 Odkazy
		`Soubor:`,       // 2 Vybrání souboru pro jeho zahození
		`#REDIRECT`,     // 3 REDIRECT odstraní duplikovaná hesla
	)
	ctg.name = name
	ctg.head = head
	var index uint64 = 0
	for _, piece := range w.Page {
		if len(rgx[0].FindAll(piece.Revision.Text, -1)) == 0 || len(rgx[3].FindAll(piece.Revision.Text, -1)) != 0 {
			continue
		}
		if attributes := f(piece); attributes != nil {
			index++
			fndLinks := rgx[1].FindAll(piece.Revision.Text, -1)
			var links []string
			for _, link := range fndLinks {
				if rgx[2].FindAllString(string(link), -1) == nil {
					links = append(links, trimLink(string(link), string(piece.Title)))
				}
			}
			vals := []string{strconv.FormatUint(index, 10)}
			for _, a := range attributes {
				vals = append(vals, a)
			}
			ctg.addNode(vals, links, piece.Revision.Text)
		}

	}
	return ctg
}

func (w *WikiData) allLinks() {
	rgx := getRegexp(
		`\[\[[^]]*\]\]`,
	)
	for _, chunk := range w.Page {
		fndLinks := rgx[0].FindAll(chunk.Revision.Text, -1)
		for _, link := range fndLinks {
			fmt.Println(string(link), " ->  ", string(chunk.Title))
		}
	}
}

// Pomáhá očistit odkazy
func trimLink(str string, title string) (newStr string) {
	new := strings.Split(string(str), "|")
	if len(new) > 2 {
		fmt.Println("\nERROR:\nOld:", str, "\nNew:", new, "\nin:", title)
		panic("Více jak dva stringy ve splitu. Funkce trimLink()")
	}
	newStr = strings.TrimPrefix(new[0], `[[`)
	newStr = strings.TrimRight(newStr, "]]")
	return newStr
}
