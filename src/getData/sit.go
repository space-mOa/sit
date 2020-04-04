package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Node v síti, hodnoty v sobě zahrnují název a atributy
type Node struct {
	name   string
	values []Values
	rgx    []string
}

// Values jednotlivé Uzly
type Values struct {
	line  []string // !!! line: ŘADA MUSÍ BÝT NÁSLEDUJÍCÍ: ID NÁZEV ATRIBUTY...
	links []string
}

func (n *Node) save(name string, head []string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(file)
	writer.Write(head)
	for _, v := range n.values {
		writer.Write(v.line)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	//	fmt.Printf("Nodes \"%v\" saved to \"%v\"\n", n.name, name)
}

func (n *Node) printNodes() {
	fmt.Println(n.name)
	for _, v := range n.values {
		fmt.Println("\nLINE:\n", v.line, "\nLINKS:")
		for _, link := range v.links {
			fmt.Println(" ", link)
		}
	}
}

func (n *Node) addNode(val []string, lin []string) {
	var v Values
	for _, i := range val {
		v.line = append(v.line, i)
	}
	for _, i := range lin {
		v.links = append(v.links, i)
	}
	n.values = append(n.values, v)
}

// Edge je vztah mezi dvěma uzly
type Edge struct {
	name string
	line [][]string
}

func (e *Edge) printEdges() {
	for _, line := range e.line {
		fmt.Println(line)
	}
}

func (e *Edge) save(name string, head []string) {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(file)
	writer.Write(head)
	for _, v := range e.line {
		writer.Write(v)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	//	fmt.Printf("Edges \"%v\" saved to \"%v\"\n", e.name, name)
}

// fromTwoNodes bere dva uzly a vytvoří pro ně hrany na základě odkazů a názvů
// názvy jsou totiž identické s první částí uvedenou v odkazech před znakem: |
func (e *Edge) fromTwoNodes(n1 Node, n2 Node, edgeName string) {
	e.name = edgeName
	var index uint64 = 0
	for _, n1V := range n1.values { // Vyber uzel z n1. vezme n1: název, n2: odkazy
		n1Title := n1V.line[1]          // Název pro n1
		for _, n2V := range n2.values { // Vyber uzel z n2
			n2V.links = removeDuplicates(n2V.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			var record []string                     // record: ID ID_N1 ID_N2 NÁZEV_N1 NÁZEV_N2
			for _, link := range n2V.links {        // Projdi všechny odkazy v uzlu
				if strings.ToLower(n1Title) == strings.ToLower(link) { // strings.ToLower je volaná, protože == rozlišuje mezi velkými a malými písmeny
					index++
					record = append(record, strconv.FormatUint(index, 10))
					record = append(record, n1V.line[0])
					record = append(record, n2V.line[0])
					record = append(record, n1V.line[1])
					record = append(record, n2V.line[1])
					e.line = append(e.line, record)
				}
			}
		}
	}
	for _, n2V := range n2.values { // Vyber uzel z n2: vezme n2: název, n1: odkazy
		n2Title := n2V.line[1]          // Název pro n2
		for _, n1V := range n1.values { // Vyber uzel z nn1
			n1V.links = removeDuplicates(n1V.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			for _, link := range n1V.links {        // Projdi všechny odkazy v uzlu
				if strings.ToLower(n2Title) == strings.ToLower(link) { // strings.ToLower je volaná, protože == rozlišuje mezi velkými a malými písmeny
					index++
					skip := false
					for _, line := range e.line { // Zkontroluj zdali už tam není stejný záznam
						if line[3] == n1V.line[1] && line[4] == n2Title {
							skip = true
						}
					}
					if !(skip) {
						var record []string // record: ID ID_N1 ID_N2 NÁZEV_N1 NÁZEV_N2
						index++
						record = append(record, strconv.FormatUint(index, 10))
						record = append(record, n1V.line[0])
						record = append(record, n2V.line[0])
						record = append(record, n1V.line[1])
						record = append(record, n2V.line[1])
						e.line = append(e.line, record)
					}
				}

			}
		}

	}
}

// socTime: Udělá vztah pokud spolu sociologové žili
func (e *Edge) socTime(n Node) {
	e.name = "lived"
	var index uint64 = 0
	for _, soc1 := range n.values { // Vezmi Sociologa 1
		soc1Born, err := strconv.ParseInt(soc1.line[2][:4], 10, 64) // Narození Sociologa 1
		if err != nil {
			fmt.Println(err)
		}
		soc1Died, err := strconv.ParseInt(soc1.line[3][:4], 10, 64) // Úmrtí Sociologa 1
		if err != nil {
			fmt.Println(err)
		}
		for _, soc2 := range n.values { // Vezmi Sociologa 2
			if soc1.line[0] == soc2.line[0] {
				continue
			}
			soc2Born, err := strconv.ParseInt(soc2.line[2][:4], 10, 64) // Narození Sociologa 2
			if err != nil {
				fmt.Println(err)
			}
			soc2Died, err := strconv.ParseInt(soc2.line[3][:4], 10, 64) // Úmrtí Sociologa 2
			if err != nil {
				fmt.Println(err)
			}
			if soc1Died == 0000 || soc2Died == 0000 { // Ti u kterých se neví jejich doba úmrtí
				continue
			}
			var record []string
			skip := false
			if soc1Died >= soc2Born { // Zkontroluj zdali spolu žili Socilog 1 a Sociolog 2
				if soc2Died <= soc1Died {
					if soc1Born <= soc2Died {
						for _, line := range e.line {
							if soc1.line[0] == line[2] && soc2.line[0] == line[1] {
								skip = true
							}
						}
						if !(skip) {
							index++
							record = append(record, strconv.FormatUint(index, 10))
							record = append(record, soc1.line[0])
							record = append(record, soc2.line[0])
							record = append(record, soc1.line[1])
							record = append(record, soc2.line[1])
							e.line = append(e.line, record)
						}
					}
				} else {
					for _, line := range e.line {
						if soc1.line[0] == line[2] && soc2.line[0] == line[1] {
							skip = true
						}
					}
					if !(skip) {
						index++
						record = append(record, strconv.FormatUint(index, 10))
						record = append(record, soc1.line[0])
						record = append(record, soc2.line[0])
						record = append(record, soc1.line[1])
						record = append(record, soc2.line[1])
						e.line = append(e.line, record)
					}
				}
			}
		}
	}
}
