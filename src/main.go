package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
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

func (w *WikiData) getAuthors(regex string, trim string) {
	var authors []string
	rgx := getRegexp(regex, trim)
	for _, chunk := range w.Page {
		fndAut := rgx[0].FindAll(chunk.Revision.Text, -1)
		for _, trimAut := range fndAut {
			aut := rgx[1].FindAll(trimAut, -1)
			for _, a := range aut {
				a := strings.TrimPrefix(string(a), `: `)
				fmt.Println(string(a))
				authors = append(authors, a)
			}
		}
	}
}

func (w *WikiData) getSociologists() (soc Node) {
	rgx := getRegexp(`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`,
		`datetime=".*"`,
		`(\[\[Kategorie:SCSg.*\]\])`,
		`\[\[[^]]*\]\]`,
		`Soubor:`,
		`[[:digit:]]-*`)
	soc.name = "Sociologist"
	var index uint64 = 0
	for _, chunk := range w.Page {
		// Najdi heslo se sociology
		fndArt := rgx[2].FindAll(chunk.Revision.Text, -1)
		if len(fndArt) == 0 {
			continue
		}
		// Najdi datum narozeni, úmrtí, másto narození
		var time string
		fndBorn := rgx[0].FindAll(chunk.Revision.Text, -1)
		for _, date := range fndBorn {
			datetime := rgx[1].FindAll(date, -1)
			fmt.Println(string(date))
			for _, t := range datetime {
				timeCh := rgx[5].FindAll(t, -1)
				for _, t := range timeCh {
					time = time + string(t)
				}
			}
		}
		fmt.Println(time, "end")
		// Najdi linky
		fndLinks := rgx[3].FindAll(chunk.Revision.Text, -1)
		if fndLinks == nil {
			fmt.Println("nic nenašlo u:", chunk.Title)
		}
		// Index Name Born Died BornIn
		var value Values
		index = index + 1
		value.line = append(value.line, strconv.FormatUint(index, 10))
		value.line = append(value.line, string(chunk.Title))
		value.line = append(value.line, time)
		// Vyhoď nepotřebné linky např. Soubor:...
		for _, link := range fndLinks {
			notWanted := rgx[4].FindAll(link, -1)
			if notWanted == nil {
				str := strings.Split(string(link), "|")
				if len(str) > 2 {
					panic("Více jak dva stringy ve splitu. funkce trimString(s)")
				}
				s := strings.TrimPrefix(str[0], `[[`)
				s = strings.TrimRight(s, "]]")
				value.links = append(value.links, s)
			}
		}
		soc.values = append(soc.values, value)
	}
	for _, sociologist := range soc.values {
		fmt.Println(sociologist)
	}
	return soc
}

// udělej více obecné
func (w *WikiData) getData(terms map[string][]string) {
	reAut := regexp.MustCompile(terms["authors"][0])
	reLinks := regexp.MustCompile(terms["links"][0])
	reAut2 := regexp.MustCompile(terms["authors"][1])
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

// Node v síti, hodnoty v sobě zahtnují název a atributy
type Node struct {
	name   string
	values []Values
	rgx    []string
}

type Values struct {
	line  []string
	links []string
}

// Edge je vztah mezi dvěma uzly
type Edge struct {
	name string
	line [][]string
}

// Save uloží do souboru CSV
type Save interface {
	save()
}

func trimString(s string) (newString string) {

	return newString
}

func getRegexp(rs ...string) []*regexp.Regexp {
	var listReg []*regexp.Regexp
	for _, s := range rs {
		reg := regexp.MustCompile(s)
		listReg = append(listReg, reg)
	}
	return listReg
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
	/*
		fmt.Println("DÉLKA", len(data.Page)-1, "\n")
		fmt.Println(string(data.Page[300].Revision.Text))
	*/
	/*
		searchTerms := map[string][]string{
			"authors": []string{`(\[\[Kategorie:Aut:.*\]\])`, `:\s[A-Za-z\sěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9]*`},
			"lang":    []string{`(<span lang=.*)(<\/span>)`},
			"links":   []string{`(\[\[([A-Za-zěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9|])*\]\])`}}
		data.getData(searchTerms)
	*/
	// data.getAuthors(`(\[\[Kategorie:Aut:.*\]\])`, `:\s[A-Za-z\sěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9]*`)
	data.getSociologists()
}
