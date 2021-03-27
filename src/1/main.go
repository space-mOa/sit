package main

import (
	"fmt"
)

func main() {
	wiki := &WikiData{}
	err := unpackFile(wiki, "../dump.xml")
	if err != nil {
		fmt.Println("Nepodařilo se nahrát data.\n", err)
	}

	soc := wiki.getCategory("sociologove", []string{}, `\[\[Kategorie:SCSg.*\]\]`)
	ins := wiki.getCategory("Instituce", []string{"id", "name", "t1", "t2", "t3"}, `\[\[Kategorie:SIZCSg.*\]\]`)
	var e Edge
	e.fromTwoNodes(ins, soc, "try")
	// ins.printValues()
	// e.printEdges()
	f := ins.filterValues([]string{"Kategorie:Státní a veřejné výzkumné instituce"}, "links")
	f.printValues()
	fmt.Println(len(f.values))
	// e.save("./")
	fmt.Println(containsListOfValues([]string{"ahoj", "  Čau"}, []string{"čau", "AHOj"}))
}
