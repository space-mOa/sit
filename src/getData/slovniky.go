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

// SIZCSg (Slovník institucionálního zázemí české sociologie)
func (w *WikiData) getSIZCSg() (inst Node) {
	rgx := getRegexp(
		`(\[\[Kategorie:SIZCSg.*\]\])`, // 0 Články spadající do Kategorie: SIZCSg (Slovník institucionálního zázemí české sociologie)
		`\[\[[^]]*\]\]`,                // 1 Odkazy
		`Soubor:`,                      // 2 Vybrání souboru pro jeho zahození
		`#REDIRECT`,                    // 3 REDIRECT odstraní duplikovaná hesla
	)
	inst.name = "Instituce"
	var index uint64 = 0
	for _, chunk := range w.Page {
		if len(rgx[0].FindAll(chunk.Revision.Text, -1)) == 0 || len(rgx[3].FindAll(chunk.Revision.Text, -1)) != 0 {
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

// VSgS (Velký sociologický slovník) - obsahuje i některá odkazová slovníky např. viz Adjustace
func (w *WikiData) getVSgS() (vsgs Node) {
	rgx := getRegexp(
		`\[\[Kategorie:VSgS.*\]\]`, // 0 najde hesla v kategorii: VSgS (Velký sociologický slovník)
		`\[\[[^]]*\]\]`,            // 1 Najde odkazy
		`#REDIRECT`,                // 2 Odstraní hesla, co mají #REDIRECT v textu
	)
	vsgs.name = "Velký sociologický slovník"
	var index uint64 = 0
	for _, chunk := range w.Page {
		if len(rgx[0].FindAll(chunk.Revision.Text, -1)) == 0 || len(rgx[2].FindAll(chunk.Revision.Text, -1)) != 0 {
			continue
		}
		index++
		fndLinks := rgx[1].FindAll(chunk.Revision.Text, -1)
		var links []string
		for _, link := range fndLinks {
			links = append(links, trimLink(string(link), string(chunk.Title)))
		}
		vsgs.addNode([]string{strconv.FormatUint(index, 10), string(chunk.Title)}, links)
	}
	return vsgs
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

// MSgS (Malý sociologický slovník)
func (w *WikiData) getMsgS() (msgs Node) {
	rgx := getRegexp(
		`\[\[Kategorie:MSgS.*\]\]`, // 0 najde hesla v kategorii: MSgS (Malý sociologický slovník)
		`\[\[[^]]*\]\]`,            // 1 Najde odkazy
		`#REDIRECT`,                // 2 Odstraní hesla, co mají #REDIRECT v textu
		`Soubor:`,                  // 3 Vybrání souboru pro jeho zahození
	)
	msgs.name = "Malý sociologický slovník"
	var index uint64 = 0
	for _, chunk := range w.Page {
		if len(rgx[0].FindAll(chunk.Revision.Text, -1)) == 0 || len(rgx[2].FindAll(chunk.Revision.Text, -1)) != 0 {
			continue
		}
		index++
		fndLinks := rgx[1].FindAll(chunk.Revision.Text, -1)
		var links []string
		for _, link := range fndLinks {
			if rgx[3].FindAllString(string(link), -1) == nil {
				links = append(links, trimLink(string(link), string(chunk.Title)))
			}
		}
		msgs.addNode([]string{strconv.FormatUint(index, 10), string(chunk.Title)}, links)
	}
	return msgs
}

// Pomáhá očistit odkazy
func trimLink(str string, title string) (newStr string) {
	new := strings.Split(string(str), "|")
	if len(new) > 2 {
		fmt.Println("Old:", str, "New:", new, "in:", title)
		panic("Více jak dva stringy ve splitu. Funkce trimLink()")
	}
	newStr = strings.TrimPrefix(new[0], `[[`)
	newStr = strings.TrimRight(newStr, "]]")
	return newStr
}

// SCSg (Slovník českých sociologů)
func (w *WikiData) getSCSg() (soc Node) {
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
