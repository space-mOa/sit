package main

import (
	"encoding/xml"
	"io/ioutil"
)

// XML soubor s obsahem Soc. encyklopedie
type WikiData struct {
	XMLName  xml.Name   `xml:"mediawiki"`
	XMLNS    string     `xml:"xmlns,attr"`
	XSI      string     `xml:"xsi,attr"`
	Location string     `xml:"schemaLocation,attr"`
	Version  string     `xml:"version,attr"`
	Lang     string     `xml:"lang,attr"`
	SiteInfo []SiteInfo `xml:"siteinfo"` // Metadata o stránce
	Page     []Page     `xml:"page"`     // Obsah pro danou HTML stránku v Soc. encyklopedii
}

// SiteInfo obsahuje informace o stránce a Namespace
type SiteInfo struct {
	XMLName    xml.Name     `xml:"siteinfo"`
	SiteName   string       `xml:"sitename"`
	DbName     string       `xml:"dbname"`
	Base       string       `xml:"base"`
	Generator  string       `xml:"generator"`
	Case       string       `xml:"case"`
	Namespaces []NameSpaces `xml:"namespaces"`
}

type NameSpaces struct {
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

// Načti XML DUMP soubor z SoC. encyklopedie
func unpackFile(v *WikiData, name string) (err error) {
	f, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(f, &v)
	if err != nil {
		return err
	}

	return nil
}
