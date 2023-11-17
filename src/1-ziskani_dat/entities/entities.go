package entities

import (
	"fmt"
	"regexp"
	"strings"
	"wikisit/wiki"
)

type Result struct {
	Nodes     []Node
	WikiData  *wiki.WikiData
	Discarded []string
}

func (r *Result) MakeTitleRegs() (result []*regexp.Regexp) {
	for _, n := range r.Nodes {
		title := wiki.ReplaceCharacters(n.Title, ' ', []rune{'\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085'})
		result = wiki.FlatAppend(result, wiki.MakeRegs(false, []string{`(\[\[` + title + `.{0,70}\]\])` + `|(` + title + `)`}))
	}
	return result
}

type Node struct {
	NodeType   string
	Title      string
	Text       []byte
	Attributes []string
	Links      []string
	HTML       string
}

type Edge struct {
	Node1    Node
	Node2    Node
	EdgeType string
}

func (e *Edge) IsEqual(c Edge) bool {
	replace := []rune{'\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085'}
	eNT1 := strings.ToLower(wiki.ReplaceCharacters(e.Node1.Title, ' ', replace))
	eNT2 := strings.ToLower(wiki.ReplaceCharacters(e.Node2.Title, ' ', replace))
	cNT1 := strings.ToLower(wiki.ReplaceCharacters(c.Node1.Title, ' ', replace))
	cNT2 := strings.ToLower(wiki.ReplaceCharacters(c.Node2.Title, ' ', replace))
	// fmt.Println(eNT1, eNT2, cNT1, cNT2)
	if (eNT1 == cNT1 && eNT2 == cNT2 || eNT1 == cNT2 && eNT2 == cNT1) && (e.EdgeType == c.EdgeType) {
		return true
	}
	return false
}

func ContainsEdge(element Edge, slice []Edge) bool {
	if len(slice) == 0 {
		return false
	}
	for _, s := range slice {
		if element.IsEqual(s) {
			return true
		}
	}

	return false
}

func RemoveDuplicateEdges(edges []Edge) (newEdges []Edge) {
	for _, e := range edges {
		if !ContainsEdge(e, newEdges) {
			newEdges = append(newEdges, e)
		}
	}
	return newEdges
}

func MakeEdgesLinks(nodes1 []Node, nodes2 []Node) (edges []Edge) {
	for _, n1 := range nodes1 {
		for _, n2 := range nodes2 {

			// Matches: N1 Title - in -> N2 Links
			matchesN1TitleN2Links := MatchLinkTitle(n1.Title, n2.Links)
			if len(matchesN1TitleN2Links) != 0 {
				newEdges := wiki.RepeatSlice(Edge{n2, n1, "odkazuje"}, len(matchesN1TitleN2Links))
				edges = wiki.FlatAppend(edges, newEdges)
			}

			// Matches: N2 Title - in -> N1 Links
			matchesN2TitleN1Links := MatchLinkTitle(n2.Title, n1.Links)
			if len(matchesN2TitleN1Links) != 0 {
				newEdges := wiki.RepeatSlice(Edge{n1, n2, "odkazuje"}, len(matchesN2TitleN1Links))
				edges = wiki.FlatAppend(edges, newEdges)
			}
		}
	}
	return edges
}

func MatchLinkTitle(title string, links []string) (matches []string) {
	replace := []rune{'\t', '\n', '\v', '\f', '\r', '\u00A0', '\u0085', '[', ']'}
	titleTrimmed := strings.ToLower(wiki.ReplaceCharacters(title, ' ', replace))
	for _, l := range links {
		linkTrimmed := strings.ToLower(wiki.ReplaceCharacters(l, ' ', replace))
		r := wiki.MakeRegs(true, []string{titleTrimmed})[0]
		if len(r.FindAllString(linkTrimmed, 1)) != 0 {
			matches = append(matches, title)
		}
	}
	return matches
}

func EdgesSaveToCSV(pathToFolder string, edges map[string][]Edge) (empty []string) {
	for k, v := range edges {
		fullFilePath := pathToFolder + k
		file, csvWriter := wiki.CreateCSVFile(fullFilePath)
		for _, e := range v {
			row := []string{e.Node1.Title, e.Node2.Title, e.EdgeType}
			csvWriter.Write(row)
		}
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}
	return empty
}
