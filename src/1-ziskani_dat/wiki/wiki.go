package wiki

import (
	"fmt"
	"os"
	"regexp"

	//	"strings"
	//	"unicode"
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
)

/*
dump.xml - XML Tree:
- MediaWiki, attr: XLMNS, XSI, Location, Version, Lang
    - SiteInfo
        - SiteName
        - DbName
        - Base
        - Generator
        - Case
        - NameSpaces
            - NameSpace, attr: Key, Case, Name
    - Page
        - Title
        - Ns
        - Redirect
        - Revision
            - Contributor
                - UserName
            - Text
*/

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
func UnpackFile(name string) (w *WikiData, err error) {
	f, err := ioutil.ReadFile(name)
	if err != nil {
		return w, err
	}
	err = xml.Unmarshal(f, &w)
	if err != nil {
		return w, err
	}

	return w, nil
}

//
//  Helper funcs for extracting data from WikiData
//

// regex: ^ = není v Golnagu
func MakeRegs(ignoreCase bool, regs []string) (results []*regexp.Regexp) {
	for _, rs := range regs {
		if ignoreCase {
			rs = `(?i)` + rs
		}
		results = append(results, regexp.MustCompile(rs))
	}

	return results
}

// fmt.Println(ApplyRegsRec([]string{"< ahoj <ahoj <ahoj <ahoj> >> <ahoj>>"}, MakeRegs(true, []string{`<.*>`, `ahoj`, `hoj`, `oj`, `.*`})))
func ApplyRegsRec(p []string, regs []*regexp.Regexp) (results []string) {
	if len(regs) == 1 {
		for _, t := range p {
			results = FlatAppend(results, regs[0].FindAllString(t, -1))
		}
		return results
	}
	for _, t := range p {
		matches := regs[0].FindAllString(t, -1)
		results = FlatAppend(results, ApplyRegsRec(matches, regs[1:]))
	}

	return results
}

// [T comparable]
func FlatAppend[T any](to []T, slice []T) []T {
	for _, s := range slice {
		to = append(to, s)
	}

	return to
}

// Aplikuje na stejný text X regulárních výrazů na stále stejný text
// Vrací text pokud všechny regulární výrazy byly v textu nalezeny
func ApplyRegsCons(p []string, regs []*regexp.Regexp) (results []string) {
	for _, t := range p {
		found := true
		for _, r := range regs {
			if len(r.FindAllString(t, -1)) == 0 {
				found = false
				break
			}
		}
		if found {
			results = append(results, t)
		} else {
			results = append(results, "")
		}
	}

	return results
}

// Helper func for filtering text from wikidate based on regex
// WHERE (int): 1 = Text
func FilterPagesRegs(w *WikiData, regs []*regexp.Regexp) (nw *WikiData, rm []string) {
	nw = &WikiData{w.XMLName, w.XMLNS, w.XSI, w.Location, w.Version, w.Lang, w.SiteInfo, []Page{}}
	for i, p := range w.Page {
		t := string(p.Revision.Text)
		if len(ApplyRegsCons([]string{t}, regs)[0]) == 0 {
			rm = append(rm, p.Title)
		} else {
			nw.Page = append(nw.Page, w.Page[i])
		}
	}

	return nw, rm
}

func ContainsElement[T comparable](element T, slice []T) bool {
	if len(slice) == 0 {
		return false
	}
	for _, s := range slice {
		if s == element {
			return true
		}
	}

	return false
}

// NBSP '\u00A0' U+00A0 (NBSP), NEL '\u0085' U+0085 (NEL)
// []rune{' ', '\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085'}
func ReplaceCharacters(s string, sub rune, replace []rune) (r string) {
	for _, ch := range s {
		if ContainsElement(ch, replace) {
			ch = sub
		}
		r = r + string(ch)
	}
	return r
}

func RepeatSlice[T any](e T, times int) []T {
	s := make([]T, times)
	for i := 0; i < times; i++ {
		s[i] = e
	}
	return s
}

func CreateCSVFile(fullPath string) (file *os.File, writer *csv.Writer) {
	file, err := os.Create(fullPath + ".csv")
	if err != nil {
		fmt.Println("\nCheck if PATH is valid\nPATH:", fullPath, "\nERROR:", err)
	}
	writer = csv.NewWriter(file)
	return file, writer
	//writer.Write(n.head)
}

func WriteCSVFile(path string, file_name string, head []string, rows [][]string) {
	f, err := os.Create(path + "/" + file_name + ".csv")
	if err != nil {
		fmt.Println("\nCheck if PATH is valid\nPATH:", fullPath, "\nERROR:", err)
	}
	csv_writer := csv.NewWriter(f)
	csv_writer.Write(head)
	csv_writer.Flush()
}

/*
// Helper func for filtering text from wikidate based on regex

// Helper funcs for comparing strings
func ChangeChars(source []string, f func (r rune) string) (trimmed []string) {
  for _, s := range source {
    var n string
    for _, ch := range s {
      n = n + string(f(ch))
    }
    trimmed = append(trimmed, n)
  }

  return trimmed
}

func ToLowerCaseRemoveWhiteSpaces(r rune) string{
	if unicode.IsSpace(r) {
  		return ""
	}

	return strings.ToLower(string(r))
}

// Helper funcs for comparing strings

func RemoveDuplicatesTrimmed(source []string) (result []string, trimmedResults []string) {
	trimmed := ChangeChars(source, ToLowerCaseRemoveWhiteSpaces)
  	if len(trimmed) != len(source) {
  		panic("Wrongly trimed slice of strings.")
  	}
	for i, t := range trimmed {
		if !ContainsElement(t, trimmedResults) {
			trimmedResults = append(trimmedResults, t)
			result = append(result, source[i])
		}
	}
	if len(result) != len(trimmedResults) {
		panic("Diffrent length.")
	}
	//fmt.Println(result, trimmedResults)
	return result, trimmedResults
}




*/
