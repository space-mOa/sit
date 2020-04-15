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
	head   []string
	values []Values
	rgx    []string
}

// Values jednotlivé Uzly
type Values struct {
	line  []string // !!! line: ŘADA MUSÍ BÝT NÁSLEDUJÍCÍ: ID NÁZEV ATRIBUTY...
	links []string
}

func (n *Node) save(path string) {
	file, err := os.Create(path + n.name + ".csv")
	if err != nil {
		fmt.Println("\nZkontrolujte zda máte vytvořenou složku uvedenou v PATH.\nPATH:", path, "\nERROR:", err)
	}
	writer := csv.NewWriter(file)
	writer.Write(n.head)
	for _, v := range n.values {
		writer.Write(v.line)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("Nodes \"%v\" saved to \"%v\"\n", n.name, path+n.name+".csv")
}

func saveNodes(path string, nodes ...Node) {
	for _, n := range nodes {
		n.save(path)
	}
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

// makeAndSaveEdges round-robin
// This function makes edges from nodes and saves them. All nodes with All nodes.
func makeAndSaveEdges(path string, nodes ...Node) {
	fmt.Printf("This function makes edges from nodes and saves them. All nodes with All nodes.")
	terminate := len(nodes) - 1
	var namesN string
	for _, n := range nodes {
		namesN = namesN + " " + n.name
	}
	copyN := nodes[:]
	fmt.Printf("\nNODES: %v\nlength: %v\nterminate at: %v\n", namesN, len(nodes), terminate)
	for i, n := range nodes {
		fmt.Printf("\n%v. %v\n", i, n.name)
		for _, withN := range copyN[i:] {
			var e Edge
			fmt.Printf(" - %v\n", withN.name)
			e.fromTwoNodes(n, withN, n.name+"_"+withN.name)
			e.save(path)
		}
	}
}

// fromTwoNodes bere dva uzly a vytvoří pro ně hrany na základě odkazů a názvů
// názvy jsou totiž identické s první částí uvedenou v odkazech před znakem: |
func (e *Edge) fromTwoNodes(n1 Node, n2 Node, edgeName string) {
	e.name = edgeName
	e.line = append(e.line, []string{"ID", n1.name + "ID1", n2.name + "ID2", n1.name, n2.name})
	var index uint64 = 0
	fmt.Println(n1.name, n2.name)
	for _, n1V := range n1.values { // Vyber uzel z n1. vezme n1: název, n2: odkazy
		n1Title := n1V.line[1]          // Název pro n1
		for _, n2V := range n2.values { // Vyber uzel z n2
			n2V.links = removeDuplicates(n2V.links) // Někdy jsou v článku uvedené stejné odkazy vícekrát, proto je odstraníme
			for _, link := range n2V.links {        // Projdi všechny odkazy v uzlu
				if strings.ToLower(n1Title) == strings.ToLower(link) { // strings.ToLower je volaná, protože == rozlišuje mezi velkými a malými písmeny
					if !(e.isThere(n2V.line[1], n1Title)) {
						index++
						e.appendEdge(strconv.FormatUint(index, 10), n1V.line[0], n2V.line[0], n1V.line[1], n2V.line[1])
					}
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
					if !(e.isThere(n1V.line[1], n2Title)) { // Zkontroluj zdali už tam není stejný záznam, akorát obráceně A B = B A
						index++
						e.appendEdge(strconv.FormatUint(index, 10), n1V.line[0], n2V.line[0], n1V.line[1], n2V.line[1])
					}
				}

			}
		}

	}
}

// Zkontroluje zdali už tam není stejný záznam A B = B A nebo A B = A B
func (e *Edge) isThere(T1, T2 string) (found bool) {
	found = false
	for _, L := range e.line {
		if strings.ToLower(L[3]) == strings.ToLower(T1) && strings.ToLower(L[4]) == strings.ToLower(T2) {
			found = true
		}
		if strings.ToLower(L[3]) == strings.ToLower(T2) && strings.ToLower(L[4]) == strings.ToLower(T1) {
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

// socTime: Udělá vztah pokud spolu sociologové žili
func (e *Edge) socTime(n Node) {
	e.name = "lived"
	var index uint64 = 0
	e.line = append(e.line, []string{"index", "Sociolog_1_ID", "Sociolog_2_ID", "Sociolog_1", "Sociolog_2"})
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
