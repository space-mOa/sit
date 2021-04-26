package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

/*
type manualy struct {
	name string
}
*/

func getcsv(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	ls, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	return ls
}

func l(v Node) {
	ls := getcsv("./data/VSgS/6/pocet.csv")
	for _, l := range ls {
		p := l[0]
		c := v.filterValues([]string{"Kategorie:" + p}, "links")
		// c := wiki.getCategory(strconv.Itoa(i)+"-", []string{"ID", "NAME"}, `\[\[Kategorie:VSgS\]\]`)
		fmt.Println("\""+"Kategorie:"+p+"\",", len(c.values))
		// i += 1
		// i := strconv.Itoa(i)
	}
}
