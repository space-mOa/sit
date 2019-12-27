package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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

func (w *WikiData) getSociologists() (soc Node) {
	rgx := getRegexp(`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`,
		`datetime=".*"`,
		`(\[\[Kategorie:SCSg.*\]\])`,
		`\[\[[^]]*\]\]`,
		`Soubor:`,
		`[[:digit:]-]*"`,
		`<span class="PERSON_DIED"><time datetime=.*>.*</time>.*</span>`)
	soc.name = "Sociologist"
	var index uint64 = 0
	for _, chunk := range w.Page {
		// Najdi heslo se sociology
		fndArt := rgx[2].FindAll(chunk.Revision.Text, -1)
		if len(fndArt) == 0 {
			continue
		}
		var born string
		fndBorn := rgx[0].FindAll(chunk.Revision.Text, -1)
		for _, date := range fndBorn {
			datetime := rgx[1].FindAll(date, -1)
			for _, t := range datetime {
				timeCh := rgx[5].FindAll(t, -1)
				for _, t := range timeCh {
					born = born + string(t)
				}
			}
		}
		// NAJDI JINÝ ZPŮSOB
		born = strings.ReplaceAll(born, "\"", "")
		var died string
		fndDied := rgx[6].FindAll(chunk.Revision.Text, -1)
		for _, date := range fndDied {
			datetime := rgx[1].FindAll(date, -1)
			for _, t := range datetime {
				timeCh := rgx[5].FindAll(t, -1)
				var count int
				for _, t := range timeCh {
					if string(t) == `"` {
						count += 1
					}
					if count == 2 {
						break
					}
					died = died + string(t)
				}
			}
		}
		if len(died) == 0 {
			died = "2030"
		}
		// NAJDI JINÝ ZPŮSOB
		died = strings.ReplaceAll(died, "\"", "")
		// Najdi linky
		fndLinks := rgx[3].FindAll(chunk.Revision.Text, -1)
		if fndLinks == nil {
			fmt.Println("nic nenašlo u:", chunk.Title)
		}
		// Index Name Born Died
		var value Values
		index = index + 1
		value.line = append(value.line, strconv.FormatUint(index, 10))
		value.line = append(value.line, string(chunk.Title))
		value.line = append(value.line, born)
		value.line = append(value.line, died)
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
	/*
		for _, sociologist := range soc.values {
			fmt.Println(sociologist)
		}
	*/
	return soc
}

// Node v síti, hodnoty v sobě zahtnují název a atributy
type Node struct {
	name   string
	values []Values
	rgx    []string
}

func (n *Node) save(name string, head []string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(file)
	writer.Write(head)
	for _, v := range n.values {
		writer.Write(v.line)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
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

func (e *Edge) save(name string, head []string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(file)
	writer.Write(head)
	for _, v := range e.line {
		writer.Write(v)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (e *Edge) makeFromOne(n Node) {
	var index uint64 = 0
	for _, soc1 := range n.values {
		soc1Died, err := strconv.ParseInt(soc1.line[3][:4], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		for _, soc2 := range n.values {
			if soc1.line[1] == soc2.line[1] {
				continue
			}
			soc2Born, err := strconv.ParseInt(soc2.line[2][:4], 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			soc2Died, err := strconv.ParseInt(soc2.line[3][:4], 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			if soc1Died > soc2Born || soc1Died >= soc2Died {
				var record []string
				index += 1
				record = append(record, strconv.FormatUint(index, 10))
				record = append(record, soc1.line[1])
				record = append(record, soc2.line[1])
				e.line = append(e.line, record)
				//	fmt.Println("---------------\n", soc1.line[1], soc1Born, soc1Died, soc2.line[1], soc2Born, soc2Died)
			}
		}
	}
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
	soc := data.getSociologists()
	soc.save("soc.csv", []string{"index", "name", "born", "died"})
	var edge Edge
	edge.makeFromOne(soc)
	edge.save("living.csv", []string{"index", "Sociolog_1", "Sociolog_2"})
}
