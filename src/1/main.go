package main

import (
	"fmt"
)

func main() {
	wiki := &WikiData{}
	err := unpackFile(wiki, "./dump.xml")
	if err != nil {
		fmt.Println("Nepodařilo se nahrát data.\n", err)
	}
	v := wiki.getCategory("SCSg", []string{"ID", "NAME"}, `\[\[Kategorie:SCSg\]\]`)
	d := wiki.getCategory("VSgS", []string{"ID", "NAME"}, `\[\[Kategorie:VSgS\]\]`)
	fmt.Println(len(v.values), len(d.values))

	e := fromTwoNodes(d, v, "ab")
	// e.printEdges()
	e.save("./e")
	c := removeDuplicates(v.printNode("Fajfr František").links)
	for _, l := range c {
		fmt.Println(l)
	}
	k := removeDuplicates(d.printNode("Demografie").links)
	for _, l := range k {
		fmt.Println(l)
	}

	// f := ins.filterValues([]string{"Kategorie:Státní a veřejné výzkumné instituce"}, "links")
	// f := v.filterValues([]string{"Kategorie:Terminologie jednotlivých tematických okruhů sociologie (s přesahem do příbuzných disciplín)"}, "links")
	// e.save("./")
	// fmt.Println(containsListOfValues([]string{"ahoj", "  Čau"}, []string{"čau", "AHOj"}))
	// v.printValues()
	// linkCat(v, "", "./data/VSgS/pocet.csv", 0)
}
