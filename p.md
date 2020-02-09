
pokud spolu žili a jsou spojeni se stejným časopisem

povede edge
unikátní časipisy -> čA: *soc ->  

[a,b,x,d,c,l,f]

(a, [b,c,x,d,l,f])
(a, [c,x,d,l,f])
(a, [l,f])
(a, [f])
(a, [])

(b, [x,b,c,x,d,l,f])
(b, [b,c,x,d,l,f])

(c, [b,c,x,d,l,f])
(c, [c,x,d,l,f])
(c, [x,d,l,f])

?([x-0], [x-1])
?([x-1], [x-2])


```javascript
function years(array) {
    let s = [...array]
    s[0].m = "k"
    let newArray = [...array, ...array]
    newArray
        .map(e => e.axis = 0)
    newArray.sort((a, b) => parseInt(a.axis) - parseInt(b.axis))
    console.log(newArray, s)
    return newArray
}


array3.find(e => {
    if ((e.Sociolog_1 === s1 && s2 === e.Sociolog_2) || (e.Sociolog_1 === s2 && s1 === e.Sociolog_2)) {
        newArray.push(e)
    }
})


let m = Math.min(...getRow(net.nodes[0], "died").map(element => element = element.getFullYear())) // Rok prvně narozeného
            console.log(m, net.nodes[0][1].died.getFullYear(), net, "☔️")
```
