- [stáhni community edition](https://neo4j.com/download-center/) - pozor community edition je špatně vidět (malý kontrast, ale je tam)
- stáhni JDK - na stránce jsou odkazy: [2.1. System requirements](https://neo4j.com/docs/operations-manual/current/installation/requirements/) 
- nastav PATH pro javu - můžeš to ověřit, že všechno funguje jak má, když napíšeš java do command line
_S tímhle jsem měl velký problém. Nakonec to funguje, ale řešení asi není úplně správné. Přidal jsem další složku, aby seděla PATH. To co je v PATH můžeš zjistit v POWERSHELLu: `echo $Env:PATH` v CMD: `echo %PATH%`_
- jdi do `neo4j-community-3...../bin/` v command line poté `neo4j console`
- v prohlížeči `http://localhost:7474/`
- pokus se přihlásit `username: neo4j` `password: neo4j` - pokud to z nějakého důvodu nefunguje, jdi do složky `neo4j-community-3...../data/dbms` a vymaž `auth`
