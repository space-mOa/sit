package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func getJournals(n Node) (journals Node) {
	rgx := getRegexp(
		`Kategorie:Vědecká a odborná periodika.*`,            // 0 Vyber Vědecká a odborná periodika
		`Kategorie:Ostatní subjekty se vztahem k sociologii`) // 1 Vyber Ostatní subjekty se vztahem k sociologii pro vyhození
	journals.name = "Casopisy"
	journals.head = []string{"ID", "nazev"}
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

func getInstitutions(n Node) (institutions Node) {
	rgx := getRegexp(
		`Kategorie:Vědecká a odborná periodika.*`, // 0 vybere Vědecká a odborná periodika
	)
	institutions.name = "Instituce"
	institutions.head = []string{"ID", "nazev"}
	var index uint64 = 0
	for _, v := range n.values {
		if strings.Contains(v.line[1], "Kategorie:") { // Přeskočí ty hesla, co mají v názvu (<title><title/>) Kategorie:...
			continue
		}
		skip := false
		for _, link := range v.links {
			if rgx[0].FindAllString(link, -1) != nil { // Najde všechny články, které neobsahují kategorii: Vědecká a odborná periodika
				skip = true
				break
			}
		}
		if !(skip) { // Skočí dovnitř pokud, heslo nemá odkaz na kategorii: Vědecká a odborná periodika
			index++
			institutions.addNode([]string{strconv.FormatUint(index, 10), v.line[1]}, v.links)
		}
	}
	return institutions
}

func removeDuplicates(slice []string) (newSlice []string) {
	if len(slice) == 0 {
		return newSlice
	}
	newSlice = append(newSlice, slice[0])
	for _, v1 := range slice {
		encountred := false
		for _, v2 := range newSlice {
			if v1 == v2 {
				encountred = true
			}
		}
		if encountred == false {
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

func checkCorrectness(n Node, checkN Node) {
	fmt.Println(n.name, checkN.name, "length is same:", len(checkN.values) == len(n.values))
	for i := 0; i < len(checkN.values)-1; i++ {
		for iV, v := range checkN.values[i].line {
			if v != n.values[i].line[iV] {
				fmt.Println("error")
				panic("PANIC")
			}
		}
		for iL, l := range checkN.values[i].links {
			if l != n.values[i].links[iL] {
				fmt.Println("ERORR", n.values[i].links[iL], l)
				fmt.Println("\n\ncheck:\n", checkN.values[i].links, "nodes\n", n.values[i].links)
				panic("PANIC")
			}
		}
	}
}

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}

	MSgS := data.getMsgS()
	SIZCSg := data.getSIZCSg()
	journals := getJournals(SIZCSg)
	institutions := getInstitutions(SIZCSg)
	/*
		VSgS := data.getVSgS()
		soc := data.getSCSg()
		saveNodes("./data/nodes/", soc, MSgS, SIZCSg, journals, institutions, VSgS)
		makeAndSaveEdges("./data/edges/", soc, MSgS, SIZCSg, journals, institutions, VSgS)
	*/
	SizcN := data.getCategory("SizcN", []string{"ID", "heslo"}, `(\[\[Kategorie:SIZCSg.*\]\])`)
	MalySgSN := data.getCategory("MalySgSN", []string{"ID", "heslo"}, `\[\[Kategorie:MSgS.*\]\]`)
	checkCorrectness(SIZCSg, SizcN)
	checkCorrectness(MSgS, MalySgSN)
	VelkySgSN := data.getCategory("VelkySgSN", []string{"ID", "heslo"}, `\[\[Kategorie:VSgS.*\]\]`)
	journalsN := data.getCategoryModify("CasopisyN", []string{"ID", "nazev"}, `(\[\[Kategorie:SIZCSg.*\]\])`, func(page Page) []string {
		rgx := getRegexp(
			`\[\[Kategorie:Vědecká a odborná periodika.*\]\]`,            // 0 Vyber Vědecká a odborná periodika
			`\[\[Kategorie:Ostatní subjekty se vztahem k sociologii\]\]`, // 1 Vyber Ostatní subjekty se vztahem k sociologii pro vyhození
		)
		text := string(page.Revision.Text)
		if rgx[0].FindAllString(text, -1) != nil {
			if nil == rgx[1].FindAllString(text, -1) {
				return []string{string(page.Title)}
			}
		}
		return nil
	})
	checkCorrectness(journals, journalsN)
	institutionsN := data.getCategoryModify("InstituceN", []string{"ID", "nazev"}, `\[\[Kategorie:SIZCSg.*\]\]`, func(page Page) []string {
		rgx := getRegexp(
			`\[\[Kategorie:Vědecká a odborná periodika.*\]\]`, // 0 vybere Vědecká a odborná periodika)
		)
		text := string(page.Revision.Text)
		if rgx[0].FindAllString(text, -1) == nil && strings.Contains(string(page.Title), "Kategorie:") == false {
			return []string{string(page.Title)}
		}
		return nil
	})
	checkCorrectness(institutions, institutionsN)
	makeAndSaveEdges("./data/", VelkySgSN, MalySgSN, SizcN, journalsN, institutionsN)
}
