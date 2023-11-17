package main

import (
	"fmt"

	//	"unicode"
	"wikisit/entities"
	"wikisit/extract"
	"wikisit/html"
	"wikisit/wiki"
)

func test2() {
	t := `Časopis s podtitulem „Revue pro výzkum populačního vývoje“ začal vycházet v roce 1959 péčí [[Státní úřad statistický|Státního úřadu statistického]] (v letech 1962–1966 Ústřední úřad státní kontroly a statistiky, později Federální statistický úřad), od rozdělení Československa (1993) je vydavatelem Český statistický úřad. Revue přitom navazuje na demograficky orientované studie dřívějšího ''[[Statistický obzor|Statistického obzoru]]'', zatímco vlastní statistika (především ekonomická statistika) přešla do časopisu ''Statistika a kontrola'' / ''[[Statistika]]'' (od 1962, resp. 1964). Od svého založení ''Demografie'' vychází čtyřikrát ročně, prvním šéfredaktorem byl František [[Fajfr František|Fajfr]] (1959–1961), následovaný Vladimírem [[Srb Vladimír|Srbem]] (1962–1976), Ladislavem Mikušem, Jánem Kurillou a Jiřinou Růžkovou. Současnou předsedkyní redakční rady je Terezie Štyglerová.
''Demografie'' je jediným českým demografickým časopisem, který přináší aktuální přehledy populačního vývoje (každoroční přehledy sňatečnosti, rozvodovosti, úmrtnosti, potratovosti, migračních pohybů a hlavní výsledky sčítání lidu) a odborné stati, recenze a zprávy. Většina příspěvků má ovšem tradičně spíše demografickopopisný než sociologickoanalytický charakter a svojí délkou nedosahuje zvyklostí běžných u odborných časopiseckých statí ve společenských vědách. Přestože metodické ukotvení časopisu zůstává v oblasti demografie, řada příspěvků má přinejmenším implicitní sociologickou relevanci, samotní sociologové ovšem do ''Demografie'' přispívají jen sporadicky (výrazněji v období před rokem 1989, kdy mnozí z nich měli omezené publikační možnosti v sociologických časopisech; do ''Demografie'' psali například Václav [[Lamser Václav|Lamser]], Ivo [[Možný Ivo|Možný]], Jiří [[Musil Jiří|Musil]] nebo Jiří [[Večerník Jiří|Večerník]], samozřejmě vedle představitelů vlastní demografické obce).
V letech 2007–2010 vycházel také anglickojazyčný výběr studií z příslušného ročníku jako elektronický časopis pod názvem ''Czech Demography''.
V roce 2012 vydal Český statistický úřad elektronickou edici naskenovaných prvních 52 ročníků časopisu (1959–2010).
''[[:Kategorie:Aut: Nešpor Zdeněk R.|Zdeněk R. Nešpor]]''&lt;br /&gt;
[[Kategorie:Aut: Nešpor Zdeněk R.]]
[[Kategorie:Vědecká a odborná periodika ]]
[[Kategorie:SIZCSg]]`
	r := wiki.ApplyRegsCons([]string{"1", "ahoj", "ahoj 1", t}, wiki.MakeRegs(false, []string{"A", "1"}))
	fmt.Println("\n\n", r)
}
func testNBSP(s string) {
	// NBSP '\u00A0' U+00A0 (NBSP), NEL '\u0085' U+0085 (NEL)
	for _, ch := range s {
		if '\u00A0' == ch {
			fmt.Println("nbsp")
		} else if ' ' == ch {
			fmt.Println("sp")
		}
	}
}
func testPrintHTML(entities entities.Result) (htmlDoc string) {
	t := ""
	// for _, n := range entities.Nodes[:60] {
	// 	t = t + n.HTML + "<br /> <br /> <hr /> <br /> <br /> \n\n"
	// }
	t = t + entities.Nodes[12].HTML + "<br /> <br /> <hr /> <br /> <br /> \n\n"
	htmlDoc = html.MakeDocHTML(t)
	fmt.Println(htmlDoc)
	return htmlDoc
}

/*
How does this work?

Dump.xml - contents of wiki

1) File is loaded to memory, data type in wiki.go
2) Things are extracted from data (dump.xml)
	- For extraction of desired pices of data regexes are used
		- any hyper-Links. 							`\[\[[^]]*\]\]`
		- specific categories on wikipedia e.g.:	`\[\[Kategorie:SIZCSg.*\]\]`
		- born, died dates: 						`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`, `"[0-9-].*"`; `<span class="PERSON_DIED">.*</span>`, `"[0-9-]+"`, `\?\?\?`, `[^0-9-]`
	- Regexes must be first compiled, for titles (in wiki page) regex helper function exists, MakeTitleRegs (entities.go)
	- Regexes can be applied in two diffrent ways:
		- Consecutively: every applied regex must be found in schearhed text, ApplyRegsCons (wiki.go)
		- Recursively: regexes are applied one by one on each previous match, ApplyRegsReq (wiki.go)
	- Special functions is available for finding hyper-links between two pages/texts, MakeEdgesLinks (entities.go)
3) Four entites are extracted from wiki dump.xml, (extract.go):
	- Sociologists, GetSociologist
	- Journals, GetJournals
	- Institutions, GetInstitutions
	- Branches and Methods, GetBranchesMethods
	- Example, GetSociologist:
		1. regexes are compiled
		2. results are saved into Result type, (entities.go)
			- found matches are saved to nodes type, (entities.go)
		3. we go through every page, WikiData.Page (wiki.go)
			- matches are .

*/

func PrintSlice[T comparable](s []T) {
	for _, p := range s {
		fmt.Println(p)
	}
}
func PrintNode(ns []entities.Node) {
	for _, n := range ns {
		fmt.Println(n.Title, n.Attributes)
		//		fmt.Println(n.Links)
	}
}
func PrintEdges(edges []entities.Edge) {
	for _, e := range edges {
		//fmt.Println(e.Node1.Title, "- " + e.EdgeType + " ->" ,e.Node2.Title)
		fmt.Println(e.Node1.NodeType + ", " + e.Node2.NodeType + ", " + e.Node1.Title + ", " + e.Node2.Title + ", ")
	}
}

func main() {
	save_files_to := "../2-vizualizace/data"
	_ = save_files_to
	w, err := wiki.UnpackFile("./data/dump.xml")
	if err != nil {
		fmt.Println(err)
	}

	// sociologists
	scg := extract.GetSociologist(w)
	// PrintNode(scg.Nodes)
	_ = scg

	// journals
	jrl := extract.GetJournals(w)
	// PrintNode(jrl.Nodes)
	_ = jrl

	// institutions
	inst := extract.GetInstitutions(w)
	//PrintNode(inst.Nodes)
	_ = inst

	// branches methods
	brn := extract.GetBranchesMethods(w)
	_ = brn
	//PrintNode(brn.Nodes)

	wiki.WriteCSVFile("./testing_folder", "test_csv", []string{"a", "b"}, [][]string{})

	/*
		// HTML example
		regsPage := [][]*regexp.Regexp{scg.MakeTitleRegs(), inst.MakeTitleRegs(),
			jrl.MakeTitleRegs(), brn.MakeTitleRegs(),
			html.Education, html.Journals, html.Create, html.Collaboration, html.Page}
		tagsPage := []html.ElementHTML{html.MakeButton("LightGreen"), html.MakeButton("LightBlue"),
			html.MakeButton("Pink"), html.MakeButton("Wheat"),
			html.MakeButton("LightCyan"),
			html.MakeButton("LightSalmon"), html.MakeButton("Lavender"),
			html.MakeButton("Thistle"), html.ElementHTML{"<br>", "<br>", len("<br><br>"), false}} // Do not look

		n := html.InsertElementToPageHTML(scg, regsPage, tagsPage)
		testPrintHTML(n)
	*/
	// Make edges from: title and links
	// edges := extract.MakeEdges(scg, jrl, inst, brn)
	// edges = extract.RemoveDuplicates(edges)
	//PrintEdges(edges["inst_brn"])

	//for _,v := range edges {
	//	PrintEdges(v)
	//}

	// NĚCO JE TU ŠPATNĚ
	//entities.EdgesSaveToCSV("./data/edges/", edges)
}

///////////////////////////////////////////////////////////////////////////////
/*







func test(w *wiki.WikiData) {
		txt := `<span class="PERSON_BORN"><time datetime="1914-09-09">9. září 1914</time> v Klech (okr. Mělník)</span>

<br /><span class="PERSON_DIED">???</span>`

		regs := wiki.MakeRegs(true, []string{"af", "sf"})
		PrintSlice(regs)
		extract.E()
		PrintSlice(wiki.FlatAppend([]int{00, 00}, []int{11, 11}))
		fmt.Println(extract.BornDiedSocDates(txt))
		wiki.RemoveDuplicatesTrimmed([]string{"AA","  a  "," a", " a ", "A"})
}
func testMatch(p entities.Result, m string) {
	fmt.Println("\n\ntestMatch()\n\n")

	alan_josef := string(p.Nodes[1].Text)
	t,_ := wiki.RemoveDuplicatesTrimmed([]string{alan_josef})
	r := wiki.MakeRegs(false, []string{`(?i)\[\[` + m + `.{0,30}\]\]`})
	ms := wiki.ApplyRegsCons([]string{alan_josef}, r)
	fmt.Println(ms, r)
	fmt.Println("\n\ntestMatch()\n\n")
	testNBSP(t[0])
}
*/
