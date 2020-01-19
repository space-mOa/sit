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

func (w *WikiData) getSIZCSg() (inst Node) {
	rgx := getRegexp(
		`(\[\[Kategorie:SIZCSg.*\]\])`, // 0 Články spadající do Kategorie: SIZCSg (Slovník institucionálního zázemí české sociologie)
		`\[\[[^]]*\]\]`,                // 1 Odkazy
		`Soubor:`)                      // 2 Vybrání souboru pro jeho zahození
	inst.name = "Instituce"
	var index uint64 = 0
	for _, chunk := range w.Page {
		fndArt := rgx[0].FindAll(chunk.Revision.Text, -1) // Najdi hesla, co patří do Slovníku institucí
		if len(fndArt) == 0 {
			continue
		}
		fndLinks := rgx[1].FindAll(chunk.Revision.Text, -1) // Najdi všechny odkazy v heslu
		var value Values                                    // value je value.line: ID INSTITUCE, value.links: ODKAZY V ČLÁNCÍCH
		index = index + 1
		value.line = append(value.line, strconv.FormatUint(index, 10))
		value.line = append(value.line, string(chunk.Title))
		for _, link := range fndLinks { // Uprav string v odkazech
			notWanted := rgx[2].FindAll(link, -1) // Vybere nechtěné linky
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
		inst.values = append(inst.values, value)
	}
	return inst
}

func (w *WikiData) getSociologists() (soc Node) {
	rgx := getRegexp(
		`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`, // 0 Vybere datum narození
		`datetime=".*"`,              // 1 Vyber časový údaj s datetime
		`(\[\[Kategorie:SCSg.*\]\])`, // 2 Články spadající do Kategorie: SCSg (Slovník českých sociologů)
		`\[\[[^]]*\]\]`,              // 3 Odkazy
		`Soubor:`,                    // 4 Vybrání souboru pro jeho zahození
		`[[:digit:]-]*"`,             // 5 Vyber číslice
		`<span class="PERSON_DIED"><time datetime=.*>.*</time>.*</span>`, // 6 vybere datum úmrtí
		`<span class="PERSON_DIED">\?\?\?</span>`)                        // 7 Vybere neznámé datum umrtí
	soc.name = "Sociologist"
	var index uint64 = 0
	for _, chunk := range w.Page {
		fndArt := rgx[2].FindAll(chunk.Revision.Text, -1) // Najdi heslo se sociology
		if len(fndArt) == 0 {
			continue
		}
		var born string
		fndBorn := rgx[0].FindAll(chunk.Revision.Text, -1) // Najdi datum narození
		for _, date := range fndBorn {
			datetime := rgx[1].FindAll(date, -1)
			for _, t := range datetime {
				timeCh := rgx[5].FindAll(t, -1) // Hledá číslice
				for _, t := range timeCh {      // Spojí číslice
					born = born + string(t)
				}
			}
		}
		// NAJDI JINÝ ZPŮSOB
		born = strings.ReplaceAll(born, "\"", "")
		var died string
		fndDied := rgx[6].FindAll(chunk.Revision.Text, -1) // Najdi datum úmrtí
		for _, date := range fndDied {
			datetime := rgx[1].FindAll(date, -1)
			for _, t := range datetime {
				timeCh := rgx[5].FindAll(t, -1) // Hledá číslice
				var count int
				for _, t := range timeCh { // Spojí číslice
					if string(t) == `"` { // Bere v potaz pouze první uvedené datum úmrtí ostatní zahodí
						count += 1
					}
					if count == 2 {
						break
					}
					died = died + string(t)
				}
			}
		}
		fndDied = rgx[7].FindAll(chunk.Revision.Text, -1) // Pokud není známo nebo ještě žije
		if fndDied != nil {
			died = "0000"
		}
		if len(died) == 0 {
			died = "2030"
		}
		// NAJDI JINÝ ZPŮSOB
		died = strings.ReplaceAll(died, "\"", "")
		fndLinks := rgx[3].FindAll(chunk.Revision.Text, -1) // Najdi linky
		if fndLinks == nil {
			fmt.Println("nic nenašlo u:", chunk.Title)
		}
		var value Values // value je: value.line: INDEX NAME BORN DIE, value.links: ODKAZY V ČLÁNCÍCH
		index = index + 1
		value.line = append(value.line, strconv.FormatUint(index, 10))
		value.line = append(value.line, string(chunk.Title))
		value.line = append(value.line, born)
		value.line = append(value.line, died)
		for _, link := range fndLinks { // Uprav string v odkazech
			notWanted := rgx[4].FindAll(link, -1) // Vybere nechtěné linky
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
	return soc
}

// Node v síti, hodnoty v sobě zahrnují název a atributy
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

func (n *Node) printNodes() {
	fmt.Println(n.name)
	for _, v := range n.values {
		fmt.Println("\nLINE:\n", v.line, "\nLINKS:")
		for _, link := range v.links {
			fmt.Println(" ", link)
		}
	}
}

func getJournals(n Node) (journals Node) {
	rgx := getRegexp(
		`Kategorie:Vědecká a odborná periodika.*`,            // 0 Vyber Vědecká a odborná periodika
		`Kategorie:Ostatní subjekty se vztahem k sociologii`) // 1 Vyber Ostatní subjekty se vztahem k sociologii pro vyhození
	journals.name = "journals"
	var index uint64 = 0
	for _, v := range n.values {
		for _, link := range v.links {
			fndJournals := rgx[0].FindAllString(link, -1) // Najdi články spadající do: Vědecká a odborná periodika
			if fndJournals != nil {
				NotWanted := rgx[1].MatchString(v.line[1]) // Vyhoď článek: Ostatní subjekty se vztahem k sociologii pro vyhození
				if NotWanted == false {
					index++
					var value Values // value je value.line: INDEX NÁZEV, value.links: ODKAZY V ČLÁNCÍCH
					value.line = append(value.line, strconv.FormatUint(index, 10))
					value.line = append(value.line, v.line[1])
					for _, l := range v.links { // Vybere nechtěné položky
						NotWanted := rgx[0].MatchString(v.line[1])
						if NotWanted == false {
							value.links = append(value.links, l)
						}
					}
					journals.values = append(journals.values, value)
				}
			}
		}
	}
	return journals
}

// Values jednotlivé Uzly
type Values struct {
	line  []string // !!! line: ŘADA MUSÍ BÝT NÁSLEDUJÍCÍ: ID NÁZEV ATRIBUTY...
	links []string
}

// Edge je vztah mezi dvěma uzly
type Edge struct {
	name string
	line [][]string
}

func (e *Edge) printEdges() {
	for _, line := range e.line {
		fmt.Println(line)
	}
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

// socTime: Udělá vztah pokud spolu sociologové žili
func (e *Edge) socTime(n Node) {
	var index uint64 = 0
	for _, soc1 := range n.values { // Vezmi Sociologa 1
		soc1Born, err := strconv.ParseInt(soc1.line[2][:4], 10, 64) // Narození Sociologa 1
		if err != nil {
			fmt.Println(err)
		}
		soc1Died, err := strconv.ParseInt(soc1.line[3][:4], 10, 64) // Úmrtí Sociologa 1
		if err != nil {
			fmt.Println(err)
		}
		for _, soc2 := range n.values { // Vezmi Sociologa 2
			if soc1.line[0] == soc2.line[0] {
				continue
			}
			soc2Born, err := strconv.ParseInt(soc2.line[2][:4], 10, 64) // Narození Sociologa 2
			if err != nil {
				fmt.Println(err)
			}
			soc2Died, err := strconv.ParseInt(soc2.line[3][:4], 10, 64) // Úmrtí Sociologa 2
			if err != nil {
				fmt.Println(err)
			}
			if soc1Died == 0000 || soc2Died == 0000 { // Ti u kterých se neví jejich doba úmrtí
				continue
			}
			var record []string
			skip := false
			if soc1Died >= soc2Born { // Zkontroluj zdali spolu žili Socilog 1 a Sociolog 2
				if soc2Died <= soc1Died {
					if soc1Born <= soc2Died {
						for _, line := range e.line {
							if soc1.line[0] == line[2] && soc2.line[0] == line[1] {
								skip = true
							}
						}
						if !(skip) {
							index++
							record = append(record, strconv.FormatUint(index, 10))
							record = append(record, soc1.line[0])
							record = append(record, soc2.line[0])
							record = append(record, soc1.line[1])
							record = append(record, soc2.line[1])
							e.line = append(e.line, record)
						}
					}
				} else {
					for _, line := range e.line {
						if soc1.line[0] == line[2] && soc2.line[0] == line[1] {
							skip = true
						}
					}
					if !(skip) {
						index++
						record = append(record, strconv.FormatUint(index, 10))
						record = append(record, soc1.line[0])
						record = append(record, soc2.line[0])
						record = append(record, soc1.line[1])
						record = append(record, soc2.line[1])
						e.line = append(e.line, record)
					}
				}
			}
		}
	}
}

// fromTwoNodes bere dva uzly a vytvoří pro ně hrany na záhladě odkázů a názvů
// názvy jsou totiž identické s první částí uvedenou v odkazech před znakem: |
func (e *Edge) fromTwoNodes(n1 Node, n2 Node) {
	var index uint64 = 0
	for _, n1V := range n1.values { // Vyber uzel z n1. vezme n1: název, n2: odkazy
		n1Title := n1V.line[1]          // Název pro n1
		for _, n2V := range n2.values { // Vyber uzel z n2
			n2V.links = removeDuplicates(n2V.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			var record []string                     // record: ID ID_N1 ID_N2 NÁZEV_N1 NÁZEV_N2
			for _, link := range n2V.links {        // Projdi všechny odkazy v uzlu
				if n1Title == link {
					index++
					record = append(record, strconv.FormatUint(index, 10))
					record = append(record, n1V.line[0])
					record = append(record, n2V.line[0])
					record = append(record, n1V.line[1])
					record = append(record, n2V.line[1])
					e.line = append(e.line, record)
				}
			}
		}
	}
	for _, n2V := range n2.values { // Vyber uzel z n2: vezme n2: název, n1: odkazy
		n2Title := n2V.line[1]          // Název pro n2
		for _, n1V := range n1.values { // Vyber uzel z nn1
			n1V.links = removeDuplicates(n1V.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			for _, link := range n1V.links {        // Projdi všechny odkazy v uzlu
				if n2Title == link {
					skip := false
					for _, line := range e.line { // Zkontroluj zdali už tam není stejný záznam
						if line[3] == n1V.line[1] && line[4] == n2Title {
							skip = true
						}
					}
					if !(skip) {
						var record []string // record: ID ID_N1 ID_N2 NÁZEV_N1 NÁZEV_N2
						index++
						record = append(record, strconv.FormatUint(index, 10))
						record = append(record, n1V.line[0])
						record = append(record, n2V.line[0])
						record = append(record, n1V.line[1])
						record = append(record, n2V.line[1])
						e.line = append(e.line, record)
					}
				}

			}
		}

	}
}

func removeDuplicates(slice []string) (newSlice []string) {
	if len(slice) == 0 {
		return newSlice
	}
	newSlice = append(newSlice, slice[0])
	for _, v1 := range slice {
		encountred := 0
		for _, v2 := range newSlice {
			if v1 == v2 {
				encountred++
			}
		}
		if encountred == 0 {
			newSlice = append(newSlice, v1)
		}
	}
	return newSlice
}

func getRegexp(rs ...string) []*regexp.Regexp {
	var listReg []*regexp.Regexp
	for _, s := range rs {
		reg := regexp.MustCompile(s)
		listReg = append(listReg, reg)
	}
	return listReg
}

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}
	soc := data.getSociologists()
	soc.save("soc.csv", []string{"index", "name", "born", "died"})
	var edge Edge
	edge.socTime(soc)
	edge.save("living.csv", []string{"index", "Sociolog_1_ID", "Sociolog_2_ID", "Sociolog_1", "Sociolog_2"})
	SIZCSg := data.getSIZCSg()
	journals := getJournals(SIZCSg)
	journals.save("casopisy.csv", []string{"index", "Nazev"})
	var edge2 Edge
	edge2.fromTwoNodes(soc, journals)
	edge2.save("souvisiSocJour.csv", []string{"index", "Sociolog_ID", "Casopis_ID", "Sociolog", "Casopis"})
}
