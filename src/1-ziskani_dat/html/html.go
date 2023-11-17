package html

import (
//	"fmt"
	"regexp"
	"wikisit/wiki"
	"wikisit/entities"
)



var (
	Education		= wiki.MakeRegs(true, []string{`\s.{0,3}studo.{0,4}\s`,`\s.{0,2}kandidát.{0,5}věd.{0,2}\s`, 
`\s.{0,3}habito.{0,5}\s`, `\s.{0,4}absolv.{0,7}`, `\s.{0,4}dokto.{0,7}\s`, `\s.{0,5}PhDr.{0,5}\s`, 
`\s.{0,4}csc.{0,4}`, `\s.{0,4}docen.{0,5}\s`, `\sprofes.{0,5}\s`})
	Journals		= wiki.MakeRegs(true, []string{`\s.{0,3}vydáv.{0,4}\s`, `\s.{0,3}naplň.{0,4}\s`, `\s.{0,3}přisp.{0,4}\s`})
	Create 			= wiki.MakeRegs(true, []string{`\s.{0,3}tvoř.{0,4}\s`, `\s.{0,3}vzni.{0,4}\s`, 
`\s.{0,2}zalo.{0,4}\s`})
	Collaboration 	= wiki.MakeRegs(true, []string{`\s.{0,3}spolup.{0,4}\s`, `\s.{0,3}účast.{0,4}\s`})
	Page			= wiki.MakeRegs(false, []string{`Knihy:`, `Studie:`, `Sborníky:`, `Překlad:`, `Literatura:`})
)

func MakeDocHTML(t string) string {
	const docHeadr = `<!doctype html><html lang="cs"><head></head><body>`
	const docFooter = `</body></html>`

	return docHeadr + t + docFooter
}

func MarkTextHTML(start, end, text string, regs []*regexp.Regexp) string {
	for _, r := range regs {
		matches := r.FindAllStringIndex(text, -1)
		if len(matches) != 0 {
			//fmt.Println(r, matches)
			text = InsertEleHTML(text, start, end, matches)
		}		
	}

	return text
}

func InsertEleHTML(txt, start_tag, end_tag string, matches [][]int) string {
	newtxt := txt
	shift  := 0
	tag_length := len(start_tag) + len(end_tag)

	for _, positions := range matches {
		word := txt[positions[0]:positions[1]]
		newtxt = newtxt[:positions[0] + shift] + start_tag + word + end_tag + newtxt[positions[1] + shift:]	
		shift = shift + tag_length
	}

	return newtxt
}

type ElementHTML struct {
	Start string
	End string
	Length int
	Single bool
}

func MakeElementHTML(single bool, tag, attrs string) (e ElementHTML) {
	start, end := "<" + tag + " " + attrs + " />", "</" + tag + ">"
	if single {
		return ElementHTML{"<" + tag + ">", "", len(start), true}
	}
	return ElementHTML{start, end, len(start) + len(end), false}
}

func MakeButton(colour string) ElementHTML {
	return MakeElementHTML(false, "button", `style="background-color:` +  colour + `; font-weight:bold;"`)
}

func InsertElementToPageHTML(c entities.Result, on [][]*regexp.Regexp, tags []ElementHTML) (r entities.Result) {
	if len(on) != len(tags) {
		panic("Selections (on [][]*regexp.Regexp) and tags must be the same lentgth.")
	}
	for i, n := range c.Nodes {
		t := wiki.ReplaceCharacters(string(n.Text), ' ', []rune{' ', '\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085'})
		for i, _ := range on {
			t = MarkTextHTML(tags[i].Start, tags[i].End, t, on[i])
		}
		c.Nodes[i].HTML = t
	}
	return c
}

/*
func MakeEleHTML(tag, attrs string, single bool) (string, string) {
	if single {
		return "<" + tag + " " + attrs + " />", ""		
	}

	return "<" + tag + " " + attrs + ">",  "</" + tag + ">"
}

func TagPage(t, colour string, regs []*regexp.Regexp) string {
	start, end := MakeEleHTML("button", `style="background-color:` +  colour + `; font-weight:bold;"`, false)

	return MarkTextHTML(start, end, t, regs)
}


func MakePageHTML(c entities.Result, keywords [][]*regexp.Regexp, kcolours []string, categories []entities.Result, ccolours []string) (r entities.Result) {
	for i, n := range c.Nodes {
		t := wiki.ReplaceCharacters(string(n.Text), ' ', []rune{' ', '\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085'})
		// categories
		for i, c := range categories {
			t = TagPage(t, ccolours[i], c.MakeTitleRegs())
		}
		// keywords 
		for i, k := range keywords {
			t = TagPage(t, kcolours[i], k)
		}
		c.Nodes[i].HTML = t
	}
	return c
} 
*/




