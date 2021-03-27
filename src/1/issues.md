## 1 Netisknutelné znaky, při vytváření hran (CLOSED)
"Karlova univerzita v Praze" "Sedláček Jan"
nemají link, i když by měli mít, "Sedláček Jan má v odkazech "Karlova univerzita v Praze" 
- ovlivněno: `fromTwoNodes()`

### Řešení
- přidány dvě funkce: `modifyString` a `removeWhiteSpaces`, druhá zmmíněná pro odstranění používá `unicode.IsSpace()` - odstraní následující znaky: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)