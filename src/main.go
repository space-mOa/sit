package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
)

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

type SiteInfo struct {
	XMLName    xml.Name      `xml:"siteinfo"`
	SiteName   string        `xml:"sitename"`
	DbName     string        `xml:"dbname"`
	Base       string        `xml:"base"`
	Generator  string        `xml:"generator"`
	Case       string        `xml:"case"`
	Namespaces []Namespacese `xml:"namespaces"`
}

type Namespacese struct {
	XMLName xml.Name    `xml:"namespaces"`
	Name    []NameSpace `xml:"namespace"`
}

type NameSpace struct {
	XMLName xml.Name `xml:"namespace"`
	Key     string   `xml:"key"`
	Case    string   `xml:"case"`
	Name    string   `xml:"namespace"`
}

type Page struct {
	XMLName  xml.Name  `xml:"page"`
	Title    string    `xml:"title"`
	Ns       string    `xml:"ns"`
	Redirect Redirecte `xml:"redirect"`
	Revision Revisione `xml:"revision"`
}

type Redirecte struct {
	XMLName xml.Name `xml:"redirect"`
	Title   string   `xml:"title,attr"`
}

type Contributore struct {
	XMLName  xml.Name `xml:"contributor"`
	UserName string   `xml:"username"`
}

type Revisione struct {
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

func getData(d *WikiData) {
	// get authors
	for a := range d.Page {
		fmt.Println(a)
	}
}

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DÉLKA", len(data.Page)-1, "\n")
	fmt.Println(string(data.Page[300].Revision.Text))
	re := regexp.MustCompile(`(<span lang=.*)(<\/span>)`)
	found := re.FindAll(data.Page[300].Revision.Text, -1)
	fmt.Println("\n\nNašlo se:", (string(found[1])), string(found[0]), string(found[2]))
}
