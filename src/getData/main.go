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

func getTimeRanges(text string) (ranges []string) {
	rgx := getRegexp(
		`<\/span>.*\(\d*.*\)`, // 0 Najde období
		`[\d–]*`,
	)
	if time := rgx[0].FindAllString(text, -1); len(time) == 1 {
		for _, rng := range strings.SplitAfterN(time[0], ",", -1) {
			var timeRange string
			for _, ch := range rgx[1].FindAllString(rng, -1) {
				timeRange += ch
			}
			ranges = append(ranges, timeRange)
			// fmt.Println(time, rng, "=>", timeRange)
		}
		// fmt.Println(ranges)
		return ranges
	}
	panic("fn getTime() našla vícero období v heslu, nežli pouze u nadpisu")
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
	Sizc := data.getCategory("Sizc", []string{"ID", "NAME"}, `(\[\[Kategorie:SIZCSg.*\]\])`)
	MalySgS := data.getCategory("MalySgS", []string{"ID", "NAME"}, `\[\[Kategorie:MSgS.*\]\]`)
	VelkySgS := data.getCategory("VelkySgS", []string{"ID", "NAME"}, `\[\[Kategorie:VSgS.*\]\]`)
	sociologists := data.getCategoryModify("Sociologove", []string{"ID", "NAME", "BORN", "DIED"}, `\[\[Kategorie:SCSg.*\]\]`, sociologists)
	journals := data.getCategoryModify("Casopisy", []string{"ID", "NAME", "T1", "T2", "T3"}, `(\[\[Kategorie:SIZCSg.*\]\])`, func(page Page) []string {
		rgx := getRegexp(
			`\[\[Kategorie:Vědecká a odborná periodika.*\]\]`,            // 0 Vyber Vědecká a odborná periodika
			`\[\[Kategorie:Ostatní subjekty se vztahem k sociologii\]\]`, // 1 Vyber Ostatní subjekty se vztahem k sociologii pro vyhození
		)
		text := string(page.Revision.Text)
		if rgx[0].FindAllString(text, -1) != nil {
			if nil == rgx[1].FindAllString(text, -1) {
				atr := []string{string(page.Title)}
				for _, time := range getTimeRanges(text) {
					atr = append(atr, time)
				}
				return atr
			}
		}
		return nil
	})
	institutions := data.getCategoryModify("Instituce", []string{"ID", "NAME", "T1", "T2", "T3"}, `\[\[Kategorie:SIZCSg.*\]\]`, func(page Page) []string {
		rgx := getRegexp(
			`\[\[Kategorie:Vědecká a odborná periodika.*\]\]`, // 0 vybere Vědecká a odborná periodika)
		)
		text := string(page.Revision.Text)
		if rgx[0].FindAllString(text, -1) == nil && strings.Contains(string(page.Title), "Kategorie:") == false {
			atr := []string{string(page.Title)}
			for _, time := range getTimeRanges(text) {
				atr = append(atr, time)
			}
			return atr
		}
		return nil
	})
	saveNodes("./data/nodes/", sociologists, MalySgS, Sizc, journals, institutions, VelkySgS)
	makeAndSaveEdges("./data/edges/", sociologists, MalySgS, Sizc, journals, institutions, VelkySgS)
	timeNodes := makeTimeRangeNodes([][]int{
		{1800, 1900},
		{1901, 1930},
		{1931, 1980},
	},
		[]Node{sociologists, journals, institutions})
	timeEdges := makeTimeEdges(timeNodes, makeEdges(sociologists, journals, institutions))
	timeNodes.save("./data/nodes/")
	timeEdges.save("./data/edges/")
}
