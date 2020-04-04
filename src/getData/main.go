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

func getInstitutions(n Node) (institutions Node) {
	rgx := getRegexp(
		`Kategorie:Vědecká a odborná periodika.*`, // 0 vybere Vědecká a odborná periodika
	)
	institutions.name = "institutions"
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

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}

	soc := data.getSCSg()
	soc.save("sociologove.csv", []string{"index", "name", "born", "died"})

	MSgS := data.getMsgS()
	MSgS.save("ms.csv", []string{"index", "ms"})

	SIZCSg := data.getSIZCSg()
	journals := getJournals(SIZCSg)
	journals.save("casopisy.csv", []string{"index", "Nazev"})
	institutions := getInstitutions(SIZCSg)
	institutions.save("instituce.csv", []string{"index", "instituce"})

	VSgS := data.getVSgS()
	VSgS.save("vs.csv", []string{"index", "vs"})

	var edge Edge
	edge.socTime(soc)
	edge.save("living.csv", []string{"index", "Sociolog_1_ID", "Sociolog_2_ID", "Sociolog_1", "Sociolog_2"})

	var socJour Edge
	socJour.fromTwoNodes(soc, journals, "SociologistsJournals")
	socJour.save("SocJour.csv", []string{"index", "Sociolog_ID", "Casopis_ID", "Sociolog", "Casopis"})

	var insSoc Edge
	insSoc.fromTwoNodes(institutions, soc, "InsSoc")
	insSoc.save("InsSoc.csv", []string{"index", "Instituce_ID", "Sociolog_ID", "Instituce", "Sociolog"})

	var insJour Edge
	insJour.fromTwoNodes(institutions, journals, "insJour")
	insJour.save("InsJour.csv", []string{"index", "Instituce_ID", "Casopis_ID", "Instituce", "Casopis"})

	var msVs Edge
	msVs.fromTwoNodes(MSgS, VSgS, "msVs")
	msVs.save("msVs.csv", []string{"index", "ms_ID", "vs_ID", "ms", "vs"})

	var socVs Edge
	socVs.fromTwoNodes(VSgS, soc, "socVs")
	socVs.save("socVs.csv", []string{"index", "vs_ID", "Sociolog_ID", "vs", "Sociolog"})

	var sziVs Edge
	sziVs.fromTwoNodes(VSgS, SIZCSg, "inSlVS")
	sziVs.save("sziVs.csv", []string{"index", "vs_ID", "siz_ID", "vs", "siz"})

	var socMs Edge
	socMs.fromTwoNodes(MSgS, soc, "socMs")
	socMs.printEdges()
	socMs.save("socMs.csv", []string{"index", "ms_ID", "Sociolog_ID", "ms", "Sociolog"})

	var sziMs Edge
	sziMs.fromTwoNodes(MSgS, SIZCSg, "inSlVS")
	sziMs.save("msSZI.csv", []string{"index", "ms_ID", "siz_ID", "ms", "siz"})
}
