:) 

## Udělat
1. [] ulož autory 
2. [] ulož články
2. [] ulož autory jako sociology 


### misto
[data imp](https://neo4j.com/docs/getting-started/current/cypher-intro/load-csv/)
[dat imp2](https://neo4j.com/docs/cypher-manual/current/clauses/load-csv/)
getRegexp(`(\[\[([A-Za-zěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9|])*\]\])`, `:\s[A-Za-z\sěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9]*`)
func newNodes(name ...string) (nodes []Node) {
	for _, n := range name {
		node := Node{name: n}
		nodes = append(nodes, node)
	}
	return nodes
}
https://encyklopedie.soc.cas.cz/w/%C5%BDalud_Augustin -> [[Masaryk Tomáš Garrigue|Masarykovi]] odkaz se skloňováním
```<title>Masaryk Tomáš Garrigue</title>```

#### úkol na příště 
- v sešitu je vymyšleno jak vytvářet sociology a instituce, naprogramuj. 
- pamatuj, že v odkazu je první část přesné znění názvu hesla na který odkazuje - to před | [[Masaryk Tomáš Garrigue|

## Co mám
název článku 
autor
odkaz

## Co obsahuje
- je tvořena 5 publikacemi: 
    1. Velký sociologický slovník 
    2. Slovník českých sociologů
    3. Slovník čes. sg. institucí 
    5. Malý sociologický slovník 
    6. Knižní bibliografie 

### 1, Velký sociologický slovník 
- Obsah: 
    - Metodologie
    - Směry, školy, teorie a koncepce sociologického a sociálního myšlení 
    - Oblasti a disciplíny sociologie
    - Příbuzné společenskovědní oblasti a disciplíny a její základní směry
    - Sociologie některých národů a států
    - Terminologie jednotlivých tématických okruhů sociologie (a předem příbuzných disciplín)
- Dále pak:
    - 1996 vydání 
    - 1988 - 93 období vzniku
    - 2350 textových hesel
    - 1860 odkazových
    - k heslům byly přidány dodatky 
    - 260 autorů různých národností
    - seznam zkratek 

### 2, Slovník českých sociologů 
Obsah:
    - medailony českých sociologů 
2003 vydání
2010 - 13 období vzniku
178 hesel/medailonů
sociologové nejenom z ČSR (ti co měli podstatný vliv na Českou sociologii)
hesla o sociolozích od 19. století 
výběr osobností především na základě jejich publikační aktivity: knižní a periodické publikace (přiznaná subjektivita), tedy jako zásadní a napsali nějaké podstatné a diferenciované dílo
do přínosu nejsou zahrnovány jen ty pozitivní - dle dněšního pohledu 
uvedená bibliografie je výběrová, ale jsou tam i ty nějvíce významná a citovaná

Název hesla: VSgS:Původní předmluva
Autor: [[Kategorie:Aut: Petrusek Miloslav| Předmluva k původnímu knižnímu vydání Velkého sociologického slovníku]]
Autor: [[Kategorie:Aut: Vodáková Alena| Předmluva k původnímu knižnímu vydání Velkého sociologického slovníku]]
Link: [[teorie]]  

heslo život má v sobě odkazy např. čase -> tento odkaz, ale neodkazuje pouze na jedno heslo, ale hned na několik -> https://encyklopedie.soc.cas.cz/w/%C4%8Cas
https://encyklopedie.soc.cas.cz/w/%C5%BDivot -čase-> https://encyklopedie.soc.cas.cz/w/%C4%8Cas

autor - heslo
[[Kategorie:Aut: Petrusek Miloslav| Předmluva k původnímu knižnímu vydání Velkého sociologického slovníku]]

https://play.golang.org/p/JhvXm7cG1wt


## Užitečné příkazy
```Powershell
go run main.go | Out-File output.txt -encoding OEM
```
```Bash 
go run main.go > output.txt
```
