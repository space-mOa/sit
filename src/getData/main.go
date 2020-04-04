package main

import (
	"fmt"
	"regexp"
	"strings"
)

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
	soc := data.getSCSg()
	Sizc := data.getCategory("Sizc", []string{"ID", "heslo"}, `(\[\[Kategorie:SIZCSg.*\]\])`)
	MalySgS := data.getCategory("MalySgS", []string{"ID", "heslo"}, `\[\[Kategorie:MSgS.*\]\]`)
	VelkySgS := data.getCategory("VelkySgS", []string{"ID", "heslo"}, `\[\[Kategorie:VSgS.*\]\]`)
	journals := data.getCategoryModify("Casopisy", []string{"ID", "nazev"}, `(\[\[Kategorie:SIZCSg.*\]\])`, func(page Page) []string {
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
	institutions := data.getCategoryModify("Instituce", []string{"ID", "nazev"}, `\[\[Kategorie:SIZCSg.*\]\]`, func(page Page) []string {
		rgx := getRegexp(
			`\[\[Kategorie:Vědecká a odborná periodika.*\]\]`, // 0 vybere Vědecká a odborná periodika)
		)
		text := string(page.Revision.Text)
		if rgx[0].FindAllString(text, -1) == nil && strings.Contains(string(page.Title), "Kategorie:") == false {
			return []string{string(page.Title)}
		}
		return nil
	})
	saveNodes("./data/nodes/", soc, MalySgS, Sizc, journals, institutions, VelkySgS)
	makeAndSaveEdges("./data/edges/", soc, MalySgS, Sizc, journals, institutions, VelkySgS)

}
