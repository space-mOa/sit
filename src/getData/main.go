package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"
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
	sociologists := data.getCategoryModify("Sociologove", []string{"id", "name", "born", "died"}, `\[\[Kategorie:SCSg.*\]\]`, sociologists)
	/*
		Sizc := data.getCategory("Sizc", []string{"id", "name"}, `(\[\[Kategorie:SIZCSg.*\]\])`)
		MalySgS := data.getCategory("MalySgS", []string{"id", "name"}, `\[\[Kategorie:MSgS.*\]\]`)
		VelkySgS := data.getCategory("VelkySgS", []string{"ID", "NAME"}, `\[\[Kategorie:VSgS.*\]\]`)
	*/
	journals := data.getCategoryModify("Casopisy", []string{"id", "name", "t1", "t2", "t3"}, `(\[\[Kategorie:SIZCSg.*\]\])`, func(page Page) []string {
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

	institutions := data.getCategoryModify("Instituce", []string{"id", "name", "t1", "t2", "t3"}, `\[\[Kategorie:SIZCSg.*\]\]`, func(page Page) []string {
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
	/*
		institutions.printNodes("Karlova univerzita v Praze")
		sociologists.printNodes("Sedláček Jan")
		/*
		// makeEdges(institutions, sociologists)
		// checkStrings(institutions, sociologists, VelkySgS, journals, MalySgS, Sizc)

		saveNodes("./data/nodes/", sociologists, MalySgS, Sizc, journals, institutions)
		makeAndSaveEdges("./data/edges/", sociologists, MalySgS, Sizc, journals, institutions)
		timeNodes := makeTimeRangeNodes([][]int{
			{1969, 1989},
		},
			[]Node{sociologists, journals, institutions})
		timeNodes.save("./data/nodes/")
		makeTimeEdgeEdgeAndSave("./data/edges/", timeNodes, makeEdges(sociologists, journals, institutions))
		output([]Node{sociologists, institutions, VelkySgS})
	*/
	saveText(highlight(sociologists, []Node{sociologists, institutions, journals}), "out.txt")
	//highlight(sociologists)
}

func highlight(node Node, col []Node) string {
	var newText string
	for _, val := range node.values {
		txt := string(val.text)
		for fndlinks(txt) != nil {
			txt = insertSpan(txt, val.line[1], node.name, fndlinks(txt), col)
		}
		newText = newText + txt
	}
	return newText
}

func saveText(text, name string) {
	err := ioutil.WriteFile(name, []byte(text), 0644)
	if err != nil {
		fmt.Println("\nERROR:", err)
	}

}

func insertSpan(text, title, nodeType string, positions []int, col []Node) string {
	// sta := ` <span style="color:red">`
	end := `</span> `
	trimed := trimLink(text[positions[0]:positions[1]], title)
	newlink := colorizeSpan(trimed, col) + trimLink(text[positions[0]:positions[1]], title) + end
	return text[:positions[0]] + newlink + text[positions[1]:]
}

func fndlinks(text string) []int {
	rgx := getRegexp(
		`\[\[[^]]*\]\]`,
	)
	return rgx[0].FindStringIndex(text)
}

func colorizeSpan(link string, nodes []Node) string {
	for _, n := range nodes {
		for _, v := range n.values {
			if strings.ToLower(checkforWs(v.line[1])) == strings.ToLower(checkforWs(link)) {
				switch n.name {
				case "Sociologove":
					return ` <span style="background-color: #66ccff">`
				case "Instituce":
					return ` <span style="background-color: #99ff33">`
				case "Casopisy":
					return ` <span style="background-color: #ff6699">`
				}
			}
		}
	}
	return ` <span style="background-color: #ff9966">`
}

func checkforWs(link string) string {
	var runes []rune
	for _, r := range link {
		if unicode.IsSpace(r) {
			runes = append(runes, []rune(" ")[0])
		} else {
			runes = append(runes, r)
		}
	}
	return string(runes)
}
