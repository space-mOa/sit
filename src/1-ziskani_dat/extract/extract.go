package extract

import (
	"regexp"
	"strings"
	"wikisit/entities"
	"wikisit/wiki"
)

//
// Links
//

func GetLinks(t string) []string {
	rl := wiki.MakeRegs(true, []string{`\[\[[^]]*\]\]`})
	return rl[0].FindAllString(t, -1)
}

//
// Sociologists
//

// // dates
func checkDiedLen(d []string) string {
	if len(d) > 1 {
		panic("Could not parse dates.")
	} else if len(d) == 0 {
		return ""
	}
	return d[0]
}

func BornDiedSocDates(t string) (born, died string) {
	rb := wiki.MakeRegs(true, []string{`<span class="PERSON_BORN"><time datetime=.*>.*</time>.*</span>`, `"[0-9-].*"`})
	rd := wiki.MakeRegs(true, []string{`<span class="PERSON_DIED">.*</span>`, `"[0-9-]+"`, `\?\?\?`, `[^0-9-]`})
	const alive = "2030"
	const unkownDate = "???"

	born = checkDiedLen(wiki.ApplyRegsRec([]string{t}, rb))
	died = checkDiedLen(wiki.ApplyRegsRec([]string{t}, rd[:1]))
	if len(born) == 0 {
		born = unkownDate
	}
	born = strings.ReplaceAll(born, `"`, "")
	if len(died) == 0 {
		return born, alive
	}
	if strings.Contains(died, unkownDate) {
		return born, unkownDate
	}

	return born, strings.ReplaceAll(wiki.ApplyRegsRec([]string{died}, rd[:2])[0], `"`, "")
}

//// dates

func GetSociologist(w *wiki.WikiData) (r entities.Result) {
	r.WikiData, r.Discarded = wiki.FilterPagesRegs(w, wiki.MakeRegs(false, []string{`\[\[Kategorie:SCSg\]\]`}))
	for _, p := range r.WikiData.Page {
		t := string(p.Revision.Text)
		b, d := BornDiedSocDates(t)
		//fmt.Printf("%-28s\t%-7s\t%s\n", p.Title, b, d)
		r.Nodes = append(r.Nodes, entities.Node{"Sociolog", p.Title, p.Revision.Text, []string{b, d}, GetLinks(t), ""})
	}
	//l, tl := wiki.RemoveDuplicatesTrimmed(r.Nodes[1].Links)
	//fmt.Println(len(r.Nodes),r.Nodes[1].Title, r.Nodes[1].Attributes, "\nLinks:\n", r.Nodes[1].Links, "\n\n", l, "\n\n", tl, "\n")
	return r
}

//
// Journals
//

func GetJournals(w *wiki.WikiData) (r entities.Result) {
	r.WikiData, _ = wiki.FilterPagesRegs(w, wiki.MakeRegs(false, []string{`\[\[Kategorie:SIZCSg.*\]\]`, `\[\[Kategorie:Vědecká a odborná periodika.*\]\]`}))
	for _, p := range r.WikiData.Page {
		r.Nodes = append(r.Nodes, entities.Node{"Časopis", p.Title, p.Revision.Text, []string{}, GetLinks(string(p.Revision.Text)), ""})
	}

	return r
}

//
// Institutions
//

func GetInstitutions(w *wiki.WikiData) (r entities.Result) {
	regs := wiki.MakeRegs(false, []string{`\[\[Kategorie:SIZCSg.*\]\]`,
		`\[\[Kategorie:Orgány řízení vědy.{0,10}\]\]`,
		`\[\[Kategorie:Státní a veřejné výzkumné instituce.{0,10}\]\]`,
		`\[\[Kategorie:Vysoké, případně vyšší školy.{0,10}\]\]`,
		`\[\[Kategorie:Vědecké společnosti a spolky.{0,10}\]\]`,
		`\[\[Kategorie:Ostatní subjekty se vztahem k sociologii.{0,10}\]\]`})
	for i := 1; i < len(regs)-1; i++ {
		filtredWiki, _ := wiki.FilterPagesRegs(w, []*regexp.Regexp{regs[0], regs[i]})
		if i == 1 {
			r.WikiData = filtredWiki
		} else {
			// Black magic is happening here
			r.WikiData.Page = wiki.FlatAppend(r.WikiData.Page, filtredWiki.Page)
		}
	}
	for _, p := range r.WikiData.Page {
		r.Nodes = append(r.Nodes, entities.Node{"Instituce", p.Title, p.Revision.Text, []string{}, GetLinks(string(p.Revision.Text)), ""})
	}

	return r
}

//
// Branches, methods, concepts and theories
//

func GetBranchesMethods(w *wiki.WikiData) (r entities.Result) {
	r.WikiData, _ = wiki.FilterPagesRegs(w, wiki.MakeRegs(false, []string{`\[\[Kategorie:VSgS.*\]\]`, `\[\[Kategorie:Směry, školy, teorie a koncepce sociologického a sociálního myšlení\]\]`}))
	for _, p := range r.WikiData.Page {
		r.Nodes = append(r.Nodes, entities.Node{"Metody", p.Title, p.Revision.Text, []string{}, GetLinks(string(p.Revision.Text)), ""})
	}
	return r
}

//
// Edges
//

// OBASUHUJE i linky autorství
func MakeEdges(scg, jrl, inst, brn entities.Result) map[string][]entities.Edge {
	edges := make(map[string][]entities.Edge)
	// Sociologists
	edges["scg_scg"] = entities.MakeEdgesLinks(scg.Nodes, scg.Nodes)
	edges["scg_jrl"] = entities.MakeEdgesLinks(scg.Nodes, jrl.Nodes)
	edges["scg_inst"] = entities.MakeEdgesLinks(scg.Nodes, inst.Nodes)
	edges["scg_brn"] = entities.MakeEdgesLinks(scg.Nodes, brn.Nodes)

	// Journals
	edges["jrl_jrl"] = entities.MakeEdgesLinks(jrl.Nodes, jrl.Nodes)
	edges["jrl_inst"] = entities.MakeEdgesLinks(jrl.Nodes, inst.Nodes)
	edges["jrl_brn"] = entities.MakeEdgesLinks(jrl.Nodes, brn.Nodes)

	// Institutions
	edges["inst_inst"] = entities.MakeEdgesLinks(inst.Nodes, inst.Nodes)
	edges["inst_brn"] = entities.MakeEdgesLinks(inst.Nodes, brn.Nodes)

	// Branches, Methods
	edges["brn_brn"] = entities.MakeEdgesLinks(brn.Nodes, brn.Nodes)
	return edges
}

func RemoveDuplicates(edges map[string][]entities.Edge) map[string][]entities.Edge {
	newEdges := make(map[string][]entities.Edge)
	for k, v := range edges {
		newEdges[k] = entities.RemoveDuplicateEdges(v)
	}
	return newEdges
}

// Tree web
