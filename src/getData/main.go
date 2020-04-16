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

func sociologists(page Page) []string {
	rgx := getRegexp(
		`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`, // 0 Vybere datum narození
		`<span class="PERSON_DIED"><time datetime=.*>.*</time>.*</span>`, // 1 vybere datum úmrtí
		`<span class="PERSON_DIED">\?\?\?</span>`,                        // 2 Vybere neznámé datum umrtí
		`datetime=".*"`, // 3 Vyber časový údaj s datetime
		`[[:digit:]-]*`, // 4 Vyber číslice
		`<span class="PERSON_DIED"><time datetime=.*>.*</time>.*<tim`, // 5 vybere datum úmrtí
	)
	text := string(page.Revision.Text)
	var born string
	for _, ch := range rgx[4].FindAllString(rgx[3].FindAllString(rgx[0].FindAllString(text, -1)[0], -1)[0], -1) {
		born = born + ch
	}
	var died string
	if s := rgx[1].FindAllString(text, -1); s == nil {
		died = "2030"
		if rgx[2].FindAllString(text, -1) != nil {
			died = "0000"
		}
	} else {
		d := rgx[1].FindAllString(text, -1)[0]
		if d1 := rgx[5].FindAllString(d, -1); d1 != nil {
			for _, ch := range rgx[4].FindAllString(rgx[3].FindAllString(d1[0], -1)[0], -1) {
				died = died + ch
			}
		} else {
			for _, ch := range rgx[4].FindAllString(rgx[3].FindAllString(d, -1)[0], -1) {
				died = died + ch
			}
		}
	}
	return []string{string(page.Title), born, died}
}

func main() {
	data := &WikiData{}
	err := unpackFile(data, "./dump.xml")
	if err != nil {
		fmt.Println(err)
	}
	sociologists := data.getCategoryModify("Sociologove", []string{"ID", "jmeno", "narozeni", "umrti"}, `\[\[Kategorie:SCSg.*\]\]`, sociologists)
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
	saveNodes("./data/nodes/", sociologists, MalySgS, Sizc, journals, institutions, VelkySgS)
	makeAndSaveEdges("./data/edges/", sociologists, MalySgS, Sizc, journals, institutions, VelkySgS)
}
