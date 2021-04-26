package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Hrany mezi dvěma vrcholy
type Edge struct {
	name string
	line [][]string // ID, ID1, ID2, NAME1, NAME2, ATRIBUTY...
}

// fromTwoNodes bere dva uzly a vytvoří pro ně hrany na základě odkazů a názvů
// názvy jsou totiž identické s první částí uvedenou v odkazech před znakem: |
func fromTwoNodes(n1 Node, n2 Node, edgeName string) (e Edge) {
	e.walk(n2, n1, 0)
	i, err := strconv.ParseUint(e.line[len(e.line)-1][0], 10, 64)
	if err != nil {
		panic("fromTwoNodes()")
	}
	e.walk(n1, n2, i)
	fmt.Println(len(e.line))
	return e
}

func (e *Edge) walk(x Node, y Node, index uint64) {
	for _, xv := range x.values {
		xt := cleanString(xv.line[1])
		for _, yv := range y.values {
			yls := removeDuplicates(yv.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			for _, yl := range yls {          // Projdi všechny odkazy
				if xt == cleanString(yl) { // Odstraní netisknutelné znaky a transformuje všechna písmenka na malá
					if !(e.isThere(xt, yv.line[1])) { // Zkontroluj zdali už tam není stejný záznam, A B = A B nebo A B = B A
						index++
						e.appendEdge(strconv.FormatUint(index, 10), xv.line[0], yv.line[0], xv.line[1], yv.line[1])
					}
				}
			}
		}
	}
}

// POMOCNÉ METODY

func (e *Edge) printEdges() {
	for _, line := range e.line {
		fmt.Println(line)
	}
}

func (e *Edge) isThere(T1, T2 string) (found bool) {
	found = false
	for _, L := range e.line {
		if cleanString(L[3]) == cleanString(T1) && cleanString(L[4]) == cleanString(T2) {
			found = true
		}
		if cleanString(L[3]) == cleanString(T2) && cleanString(L[4]) == cleanString(T1) {
			found = true
		}
	}
	return found
}

func (e *Edge) appendEdge(values ...string) {
	var record []string
	for _, v := range values {
		record = append(record, v)
	}
	e.line = append(e.line, record)
}

func (e *Edge) appendEdges(ne Edge) {
	for _, l := range ne.line {
		e.appendEdge(l...)
	}
}

// POMOCNÉ FUNKCE

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

// Odstraní netisknutelné znaky: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)
// + transformuje všechna VELKÁ písmenka na malá
func cleanString(s string) (n string) {
	return strings.ToLower(removeWhiteSpaces(s))
}

// Odstraní netisknutelné znaky: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)
func removeWhiteSpaces(s string) (n string) {
	for _, ch := range s {
		if !(unicode.IsSpace(ch)) {
			n += string(ch)
		}
	}
	return n
}

func linkCat(n Node, k string, p string, fd int) {
	ls := getcsv(p)
	for _, ln := range ls {
		for _, v := range n.values {
			for _, ls := range v.links {
				if ln[fd] == ls {
					fmt.Println(ln[fd], ls, "\n", v.line[1], "\n")
				}
			}
		}
	}
}

func (e *Edge) save(path string) {
	file, err := os.Create(path + e.name + ".csv")
	if err != nil {
		fmt.Println("\nZkontrolujte zda máte vytvořenou složku uvedenou v PATH.\nPATH:", path, "\nERROR:", err)
	}
	writer := csv.NewWriter(file)
	for _, v := range e.line {
		writer.Write(v)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("Edges \"%v\" saved to \"%v\"\n", e.name, path+e.name+".csv")
}
