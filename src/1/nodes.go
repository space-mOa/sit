package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Vrchol pro vybrané extrahované entity = kategorie
type Node struct {
	name   string
	head   []string
	values []Values
}

// Extrahované hodnoty pro danou entitu = kategorie
type Values struct {
	line  []string // !!! ŘADA MUSÍ BÝT NÁSLEDUJÍCÍ: ID NÁZEV ATRIBUTY...
	links []string // Odkazy pro stránku
	text  []byte   // Text stránky
}

// Přiřadí pro daný vrchol další atributy
// POZOR: mutuje state
func (n *Node) addNode(val []string, lin []string, text []byte) {
	var v Values
	for _, i := range val {
		v.line = append(v.line, i)
	}
	for _, i := range lin {
		v.links = append(v.links, i)
	}
	v.text = text
	n.values = append(n.values, v)
}

// Vytvoří ze stringu regulární výrazy
// Vrací slice s regulárními výrazy
func makeRegExp(rs ...string) []*regexp.Regexp {
	var rgx []*regexp.Regexp
	for _, s := range rs {
		r := regexp.MustCompile(s)
		rgx = append(rgx, r)
	}
	return rgx
}

// Upraví odkazy získané z dané stránky
func modifyLinks(discard []string, links []string, title string) (modifiedLinks []string) {
	rgx := makeRegExp(discard...) // Regurální výrazy pro odkazy, které chceme zahodit
	for _, l := range links {
		fndLink := false
		for _, r := range rgx {
			if r.FindAllString(string(l), 1) != nil {
				fndLink = true
				break
			}
		}
		if !(fndLink) {
			modifiedLinks = append(modifiedLinks, trimLink(string(l), string(title)))
		}
	}
	return modifiedLinks
}

// Očistí odkaz
func trimLink(str string, title string) (newStr string) {
	new := strings.Split(string(str), "|")
	if len(new) > 2 {
		fmt.Println("\nERROR:\nOriginální string:", str, "\nNový string:", new, "\nV textu:", title)
		panic("Více jak dva stringy ve splitu. Funkce trimLink()")
	}
	newStr = strings.TrimPrefix(new[0], `[[`)
	newStr = strings.TrimRight(newStr, "]]")
	return newStr
}

// Extrahuje data na základě kategorie: k string
// Higher-order function
// Bere jako argument funkcí, která dovoluje PŘIDAT další atributy mimo Název, Odkazy
// Vrací vrchol: ctg Node s extrahovanými daty
func (w *WikiData) getCategoryModify(name string, head []string, k string, f func(p Page) []string) (ctg Node) {
	rgx := makeRegExp(
		k,               // 0 Kategorie
		`#REDIRECT`,     // 1 REDIRECT, odstraní duplikovaná hesla
		`\[\[[^]]*\]\]`, // 2 Odkazy
	)
	ctg.name = name            // Název pro skupinu vrcholů
	ctg.head = head            // Seznam atributů
	var index uint64 = 0       // Index pro jednotlivé vrchol
	for _, p := range w.Page { // Prochází stránku po stránce v "dump.xml"
		if len(rgx[0].FindAll(p.Revision.Text, -1)) == 0 || len(rgx[1].FindAll(p.Revision.Text, -1)) != 0 { // Vybere text s danou kategorií
			continue
		}
		if attributes := f(p); attributes != nil {
			index++
			fndLinks := rgx[2].FindAllString(string(p.Revision.Text), -1)
			vals := []string{strconv.FormatUint(index, 10)}
			for _, a := range attributes {
				vals = append(vals, a)
			}
			ctg.addNode(
				vals, // ID
				modifyLinks([]string{`Soubor:`}, fndLinks, p.Title), // Title HTML stránky
				p.Revision.Text) // Text pro danou HTML stránku
		}

	}
	return ctg
}

// Extrahuje data na základě kategorie: k string
func (w *WikiData) getCategory(name string, head []string, k string) (ctg Node) {
	rgx := makeRegExp(
		k,               // 0 Kategorie
		`#REDIRECT`,     // 1 REDIRECT, odstraní duplikovaná hesla
		`\[\[[^]]*\]\]`, // 2 Odkazy
	)
	ctg.name = name            // Název pro skupinu vrcholů
	ctg.head = head            // Seznam atributů
	var index uint64 = 0       // Index pro jednotlivé vrchol
	for _, p := range w.Page { // Prochází stránku po stránce v "dump.xml"
		text := p.Revision.Text
		if len(rgx[0].FindAll(text, -1)) == 0 || len(rgx[1].FindAll(text, 1)) != 0 {
			continue
		}
		index++
		fndLinks := rgx[2].FindAllString(string(text), -1)
		ctg.addNode([]string{
			strconv.FormatUint(index, 10), // ID
			string(p.Title)},              // Title HTML stránky
			modifyLinks([]string{`Soubor:`}, fndLinks, p.Title), // Odkazy pro danou HTML stránku
			p.Revision.Text) // Text pro danou HTML stránku
	}
	return ctg
}

// POMOCENÉ METODY

// Vytiskno všechny linky
func (v *Values) printLinks() {
	for _, l := range v.links {
		fmt.Printf("\t %v\n", l)
	}
}

// Vytisne uložené "Values" pro daný Vrchol
func (n *Node) printValues() {
	fmt.Println("\n", n.name, "\n")
	for _, v := range n.values {
		fmt.Printf("\n\nTITLE:\t %v, %v\n", v.line[0], v.line[1])
		fmt.Println("LINE:\t", v.line)
		fmt.Println("LINKS:")
		v.printLinks()
	}
}

// Vytiskne Vales na základě zadaného titlu
func (n *Node) printNode(title string) {
	fmt.Println("\n", n.name)
	for _, v := range n.values {
		if title == v.line[1] {
			fmt.Println("\nLINE:\n", v.line, "\nLINKS:")
			for _, link := range v.links {
				fmt.Println(" ", link)
			}
		} else if title == "" {
			fmt.Println("\nLINE:\n", v.line, "\nLINKS:")
			for _, link := range v.links {
				fmt.Println(" ", link)
			}
		}
	}
}

// Uloží skupinu "Vrcholů" do ".CSV" souboru
func (n *Node) save(path string) {
	file, err := os.Create(path + n.name + ".csv")
	if err != nil {
		fmt.Println("\nZkontrolujte zda máte vytvořenou složku uvedenou v PATH.\nPATH:", path, "\nERROR:", err)
	}
	writer := csv.NewWriter(file)
	writer.Write(n.head)
	for _, v := range n.values {
		writer.Write(v.line)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("Nodes \"%v\" saved to \"%v\"\n", n.name, path+n.name+".csv")
}

// Uloží slice skupiny vrcholů
func saveNodes(path string, nodes ...Node) {
	for _, n := range nodes {
		n.save(path)
	}
}

// Ignoruje velká i malá písmena a mezery
func contains(s string, in []string) (there bool) {
	there = false
	for _, i := range in {
		if modifyString(s) == modifyString(i) {
			there = true
		}
	}
	return there
}

// Ignoruje velká i malá písmena a mezery
func containsListOfValues(equals []string, to []string) bool {
	for _, e := range equals {
		if !contains(e, to) {
			return false
		}
	}
	return true
}

// Vyfiltruje Values na základě dle zadanách stringů
// Porovnává, buď to v value.line nebo value.links
func (n *Node) filterValues(want []string, where string) (filtered Node) {
	for _, v := range n.values {
		switch where {
		case "line":
			if containsListOfValues(want, v.line) {
				filtered.values = append(filtered.values, v)
			}
		case "links":
			if containsListOfValues(want, v.links) {
				filtered.values = append(filtered.values, v)
			}

		}
	}
	return filtered
}
