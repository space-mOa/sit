package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
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

func (w *WikiData) getData(terms map[string]string) {
	reAut := regexp.MustCompile(terms["authors"])
	reLinks := regexp.MustCompile(terms["links"])
	reAut2 := regexp.MustCompile(`:\s[A-Za-z\sěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9]*`)
	for _, d := range w.Page {
		fmt.Println("Název hesla:", d.Title)
		fnAut := reAut.FindAll(d.Revision.Text, -1)
		fnLinks := reLinks.FindAll(d.Revision.Text, -1)
		for _, f := range fnAut {
			fmt.Println("Autor:", string(f))
			l := reAut2.FindAll(f, -1)
			for _, n := range l {
				n := strings.TrimPrefix(string(n), `: `)
				fmt.Println(string(n))
			}
		}
		for _, f := range fnLinks {
			fmt.Println("Link:", string(f))
		}
		fmt.Println("\n")
	}
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

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DÉLKA", len(data.Page)-1, "\n")
	fmt.Println(string(data.Page[300].Revision.Text))
	/*
		re := regexp.MustCompile(`(<span lang=.*)(<\/span>)`)
		found := re.FindAll(data.Page[300].Revision.Text, -1)
		fmt.Println("\n\nNašlo se:", (string(found[1])), string(found[0]), string(found[2]))
	*/
	searchTerms := map[string]string{
		"authors": `(\[\[Kategorie:Aut:.*\]\])`,
		"lang":    `(<span lang=.*)(<\/span>)`,
		"links":   `(\[\[([A-Za-zěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9|])*\]\])`}
	data.getData(searchTerms)
}
