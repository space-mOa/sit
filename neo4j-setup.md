- [stáhni community edition](https://neo4j.com/download-center/) - pozor community edition je špatně vidět (malý kontrast, ale je tam)
- stáhni JDK - na stránce jsou odkazy: [2.1. System requirements](https://neo4j.com/docs/operations-manual/current/installation/requirements/) 
- nastav PATH pro javu - můžeš to ověřit, že všechno funguje jak má, když napíšeš java do command line. To co je v PATH můžeš zjistit v POWERSHELLu: `echo $Env:PATH` v CMD: `echo %PATH%`
- jdi do `neo4j-community-3...../bin/` v command line poté `neo4j console`
- v prohlížeči `http://localhost:7474/`
- pokus se přihlásit `username: neo4j` `password: neo4j` - pokud to z nějakého důvodu nefunguje, jdi do složky `neo4j-community-3...../data/dbms` a vymaž `auth`
- `.\neo4j-admin.bat set-initial-password novéHeslo`
- `.\cypher-shell.bat -u "neo4j" -p "neo4jj"` nebo `.\cypher-shell.bat -u "neo4j" -p "neo"`
- pro import potřebuješ apoc -> stáhni [tady](https://github.com/neo4j-contrib/neo4j-apoc-procedures), potřebuješ mít stejnou verzi apocu s neo4j
- poté přesuň stáhnutý `.jar`(apoc) do složky `tvoje_složka_s_neo4j\plugins` v neo4j 
- do `neo4j.conf` přidej toto: `dbms.security.procedures.unrestricted=apoc.*`, najdeš to ve složce `tvoje_složka_s_neo4j\conf`
- následně dej data do složky `tvoje_složka_s_neo4j\import` 

IMPORT
```
CREATE CONSTRAINT ON (instituce:Instituce) ASSERT instituce.id IS UNIQUE;

LOAD CSV WITH HEADERS FROM "file:///instituce.csv" as csvLine
CREATE (i:Instituce {id: toInteger(csvLine.index), name: csvLine.instituce});



CREATE CONSTRAINT ON (malySgS:MalySgS) ASSERT malySgS.id IS UNIQUE;

LOAD CSV WITH HEADERS FROM "file:///malySgS.csv" as csvLine
CREATE (m:MalySgS {id: toInteger(csvLine.index), title: csvLine.ms});




CREATE CONSTRAINT ON (velkySgS:VelkySgS) ASSERT velkySgS.id IS UNIQUE;

LOAD CSV WITH HEADERS FROM "file:///velkySgS.csv" as csvLine
CREATE (v:VelkySgS {id: toInteger(csvLine.index), title: csvLine.vs});




CREATE CONSTRAINT ON (casopis:Casopis) ASSERT casopis.id IS UNIQUE;

LOAD CSV WITH HEADERS FROM "file:///casopisy.csv" as csvLine
CREATE (c:Casopis {id: toInteger(csvLine.index), name: csvLine.Nazev});




CREATE CONSTRAINT ON (sociolog:Sociolog) ASSERT sociolog.id IS UNIQUE;

LOAD CSV WITH HEADERS FROM "file:///sociologove.csv" as csvLine
CREATE (s:Sociolog {id: toInteger(csvLine.index), name: csvLine.name, born: csvLine.born, died: csvLine.died});

```
