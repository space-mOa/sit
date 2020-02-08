:) 

## Udělat
1. [x] ulož sociology
2. [x] ulož časopisy
3. [x] diagram


### misto
[data imp](https://neo4j.com/docs/getting-started/current/cypher-intro/load-csv/)
[dat imp2](https://neo4j.com/docs/cypher-manual/current/clauses/load-csv/)
[dat imp3](https://neo4j.com/graphgist/importing-csv-files-with-cypher)
https://encyklopedie.soc.cas.cz/w/%C5%BDalud_Augustin -> [[Masaryk Tomáš Garrigue|Masarykovi]] odkaz se skloňováním
```<title>Masaryk Tomáš Garrigue</title>```

možnost omezit na aktivní dobu věku

původně se směry
26 554
/2 bez
13 277

sociolog bude spolu svázány pokud žili ve stejný čassy

IMPORT
```
CREATE CONSTRAINT ON (sociolog:Sociolog) ASSERT sociolog.id IS UNIQUE

LOAD CSV WITH HEADERS FROM "file:///soc.csv" as csvLine
CREATE (s:Sociolog {id: toInteger(csvLine.index), name: csvLine.name, born: csvLine.born, died: csvLine.died});

LOAD CSV WITH HEADERS FROM "file:///living.csv" AS csvLine
MATCH (s1:Sociolog {id: toInteger(csvLine.Sociolog_1_ID)}),(s2:Sociolog {id: toInteger(csvLine.Sociolog_2_ID)})
CREATE (s1)-[:LIVED_WITH]->(s2);


Sociolog_1_ID,Sociolog_2_ID

MATCH (n)
DETACH DELETE n;

MATCH ()-[r:LIVED_WITH]-() 
DELETE r;

USING PERIODIC COMMIT 500
LOAD CSV WITH HEADERS FROM "file:///roles.csv" AS csvLine
MATCH (person:Person {id: toInteger(csvLine.personId)}),(movie:Movie {id: toInteger(csvLine.movieId)})
CREATE (person)-[:PLAYED {role: csvLine.role}]->(movie)
```

https://neo4j.com/docs/getting-started/current/cypher-intro/load-csv/
https://www.quackit.com/neo4j/tutorial/neo4j_delete_a_node_using_cypher.cfm

EXPORT 
```
/// UZLY
CALL apoc.export.json.query("MATCH (n:Sociolog) RETURN collect(n{.id, .name, .born, .died}) as nodes","nodes.json");

/// HRANY 
/// jako csv
CALL apoc.export.csv.query("MATCH (n)-[r:LIVED_WITH]-(l) 
RETURN n.id as v1_id, l.id as v2_id, n.name as v1_name, l.name as v2_name", "query.csv", {});


CALL apoc.export.json.query("MATCH (n)-[r:LIVED_WITH]-(l) RETURN collect(n{.name, .id}) as nodes","query.json");
CALL apoc.export.json.query("MATCH (n)-[r:LIVED_WITH]-(l) RETURN distinct(n{.id, .name, .born, .died}) as nodes","query.json");

/// příklady
CALL apoc.export.json.query("MATCH (u:User)-[r:KNOWS]->(d:User) RETURN u {.*}, d {.*}, r {.*}","/tmp/map.json",{params:{age:10}})

MATCH (nod:User)
MATCH ()-[rels:KNOWS]->()
WITH collect(nod) as a, collect(rels) as b
CALL apoc.export.json.data(a, b, "tmp/data.json", null)
YIELD nodes, relationships, properties, file, source,format, time
RETURN *

```

#### úkol na příště 
- v sešitu je vymyšleno jak vytvářet sociology a instituce, naprogramuj. 
- pamatuj, že v odkazu je první část přesného znění názvu hesla na který odkazuje - to před | [[Masaryk Tomáš Garrigue|
- máš skoro vytvořený uzel. Musíš upravit odkazy - tak aby se neopakovaly, nevynechávaly a zformátovat 
- U Von Wieser Friedrich neneajde odkazy, ikdyž tam jsou - Německý uni. v Praze 
- U Von Wieser Friedric - pokud prohledávám se speciálními znaky XML - nic to nenajde, když dám UNESCAPE najde to! \[\[[A-Za-zěščřžýáíéůúťňďĚŠČŘŽÝÁÍÉÚŮŤĎŇ0-9|\s]*\]\]

## Co mám
- vrcholy: socilogy a časopisy
- vztahy: žili spolu, souvisí s časopisem
- první vizualizace

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
- Obsah:
    - medailony českých sociologů 
- Dále pak:    
    - 2003 vydání
    - 2010 - 13 období vzniku
    - 178 hesel/medailonů
    - sociologové nejenom z ČSR (ti co měli podstatný vliv na Českou sociologii)
    - hesla o sociolozích od 19. století 
    - výběr osobností především na základě jejich publikační aktivity: knižní a periodické publikace (přiznaná subjektivita), tedy jako zásadní a napsali nějaké podstatné a diferenciované dílo
    - do přínosu nejsou zahrnovány jen ty pozitivní - dle dněšního pohledu 
    - uvedená bibliografie je výběrová, ale jsou tam i ty nějvíce významná a citovaná

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
