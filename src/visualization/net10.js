d3.csv("./data/nodes/Sociologove.csv").then(soc => {
    d3.csv("./data/nodes/Instituce.csv").then(ins => {
        d3.csv("./data/nodes/Casopisy.csv").then(cas => {
            d3.csv("./data/edges/Sociologove_Casopisy.csv").then(soc_cas => {
                d3.csv("./data/edges/Sociologove_Instituce.csv").then(soc_ins => {
                    d3.csv("./data/edges/Casopisy_Instituce.csv").then(cas_ins => {
                        d3.csv("./data/edges/Casopisy_Casopisy.csv").then(cas_cas => {
                            d3.csv("./data/edges/Instituce_Instituce.csv").then(ins_ins => {
                                d3.csv("./data/edges/Sociologove_Sociologove.csv").then(soc_soc => {
                                    console.log(ins)
                                    draw({
                                        nodes: giveYs(giveIDs([soc, ins, cas], ["sociolog", "instituce", "casopis"]), ["sociolog", "sociolog", "sociolog", "sociolog", "sociolog", "instituce", "casopis"]),
                                        edges: makeEdges(
                                            giveIDs([soc, ins, cas], ["sociolog", "instituce", "casopis"]), 
                                            [soc_cas, soc_ins, cas_ins, cas_cas, ins_ins, soc_soc]),
                                    })
                                })
                            })
                        })
                    })
                })
            })
        })
    })
})

// Expert Data Visualization - kniha

// Potřeba opravit: čvut - chybí, čtvrtý čas; opravit data končící např. -46
// Rozdělit časovou aktivitu - na nečinost v sociologii dle skeče
// Opravit 2030

// upravit na to uzly na line tak, aby šli po sobě časově a zároveň uzly vedle sebe by měli co nejvíce vztahů, rozdělit do období např. ?


function draw(data) {
    let w = (innerWidth * 2)
    let h = (innerHeight  * 50)
    let rH = 70
    let cvs = setCanvas(w,h)
    console.log(data)
    let tS = d3.scaleTime()
        .domain([new Date("1781"), new Date("2030")])
        .range([80, w])
        .nice()
    let xA = d3.axisBottom()
        .scale(tS)
    cvs.append("g")
        .attr("class", "xAxis")
        .call(xA);
        // .attr("transform", `translate( ${0}, ${h-180})`)
    /*
    cvs.append("g").selectAll("links")
    .data(data.edges)
    .enter()
    .append("line")
        .attr("x1", d => tS(getNode(data.nodes, d.startName, "name")[0].t1s))
        .attr("y1", d => getNode(data.nodes, d.startName, "name")[0].y)
        .attr("x2", d => tS(highestTimeForANode(getNode(data.nodes, d.endName, "name")[0])))
        .attr("y2", d => getNode(data.nodes, d.endName, "name")[0].y)
        .style("stroke",  d => { if (d.startName === "Masaryk Tomáš Garrigue" || d.endName === "Masaryk Tomáš Garrigue") {return "black"} else {return"#FF70A6"}})
        .style("stroke-width", 2)   
    */
    let grl = cvs.append("g").selectAll("rectangles")
        .data(data.nodes)
        .enter()
        .append("g").attr("class", "g.rectangles");

    grl.append("rect")
        .attr("x", d => rectangle(tS(d.t1s), d.y, tS(highestTimeForANode(d)), rH).x)
        .attr("y", d => rectangle(tS(d.t1s), d.y, tS(highestTimeForANode(d)), rH).y)
        .attr("width", d => rectangle(tS(d.t1s), d.y, tS(highestTimeForANode(d)), rH).width)
        .attr("height", d=> rectangle(tS(d.t1s), d.y, tS(highestTimeForANode(d)), rH).height)
        .attr("fill", d => colors(d.nodeName))
        /*
            .attr("stroke", "red")
            .attr("stroke-width", 5)
        */
    grl.append("text")
        .text(d => d.name)
        .attr("fill", "white")
        .attr("text-anchor", "left")
        .style("font", "15px")
        .attr("x", d => rectangle(tS(d.t1s), d.y, tS(d.t1e), rH).x + 10)
        .attr("y", d => rectangle(tS(d.t1s), d.y, tS(d.t1e), rH).y + 40)
        
    console.log(rectangle(50, 20, 150, 40), getNode(data.nodes, data.edges[0].startName, "name"))
}

function getNode(nodes, value, key) {
    let node = nodes.filter(n => value === n[key])
    if (node.length > 1 ) { 
        console.log("ERROR, fn getNode()")
    } 
    return node
}

// Přiřadí nové ID pro uzly, které spojí dohromady 
// + přiřadí jméno např. všechny UZLY reprezenrující časopisy budou mít atribut nodeName: časopisy
function giveIDs(nodes, ids) {
    let newID = 0
    nodes.forEach((e, i) => {
        e.forEach(n => {
            n = dates(n, ids[i])
            n.nodeName = ids[i]
            n.newID = newID
            newID += 1
        })
    })
    let newNodes = nodes.flat()
    isUnique(newNodes, "newID") ? console.log("nodes have unique ids") : console.log("ERROR, nodes dont have unique ids");
    return newNodes
}

// Vytvoří source + target na základě dodaných vztahů
function makeEdges(nodes, edges) {
    let newEdges = []
    edges
        .flat()
        .forEach(e => {
            newEdges.push(makeEdge(e, nodes))
        })
    return newEdges
}

// Vytvoří novou hranu 
// !!! předpokládá u hran toto pořadí: ID, ID1, ID2, NÁZEV1, NÁZEV2 v hran -> ID: 1, NÁZEV1: 3, ID2: 2...
function makeEdge(edge, nodes) {
    let start = lookUp(edge[Object.keys(edge)[1]], edge[Object.keys(edge)[3]], nodes) // start node
    let end = lookUp(edge[Object.keys(edge)[2]], edge[Object.keys(edge)[4]], nodes)   // end node
    return {
        source: start["newID"], 
        target: end["newID"],
        startName: start[Object.keys(start)[1]],
        endName: end[Object.keys(end)[1]],
        edgeName: start["nodeName"] + "_" + end["nodeName"],
        atr: edge["atr"]
    }
}

// Vyhledá newID pro id a value uvedenou ve vztahu 
// !!! předpokládá u uzlů toto pořadí: ID, NÁZEV, ATRIBUTY -> ID: 0, NÁZEV: 1, ATRIBUT1: 2...
function lookUp(id, value, nodes) {
    let newNode = nodes
        .filter(e => {
            return (parseInt(e[Object.keys(e)[0]]) === parseInt(id) && e[Object.keys(e)[1]] === value); // Pokud se rovná v nodes ID z HranyX a Název X -> vrať Node
        })
    if (newNode.length >= 2) {
        console.log("ERROR, fn lookUp() našla více jak jeden uzel", newNode)
    } else {
        return newNode[0]
    }    
}

// Kontrolní funkce zdali jsou hodnoty unikátní
function isUnique(array, key) {
    let encountred = []
    let unique = true
    array.forEach((e,i) => {
        if (i === 0) {
            encountred.push(e[key])
        } else {
            if (encountred.includes(e[key])) {
                console.log("ERROR", e[key], e, encountred)
                encountred.push(e[key])
                unique = false
            } else {
                encountred.push(e[key])
            }
        }
    })
    return unique
}

function setCanvas(w,h) {
    return paper = d3.select("body")
    .append("svg")
        .attr("width", w)
        .attr("height", h)
}

function dates(node, nodeName) {
    switch (nodeName) {
        case "sociolog":
            node.t1s = parseTime(node.born)
            node.t1e = parseTime(node.died)
            if (node.t1e.getTime() === new Date ("0000").getTime()) {
                ny = node.t1s.getFullYear() + 2
                // console.log(node.name, node.t1s.getFullYear(), "died:",node.t1e.getFullYear(), "| 0 means unkown death, date of death is set to born + 2 years:", new Date(ny.toString(10)).getFullYear(), new Date(ny.toString(10)))
                node.t1e = new Date(ny.toString(10))
                console.log("for", node.name, "0000 ->", node.t1e.getFullYear())
            }
            node.t2 = ""
            node.t3 = ""
            return node
        case "instituce":
            node = individualTimes(node)
            return node
        case "casopis":
            node = individualTimes(node)
            return node
        default:
            console.log("ERROR, něco je špatně s nodeNames fn dates().", nodeName, node)
            return node
    }
}

function parseTime(string) {
    return new Date(string)
}

function individualTimes(node) {
    ["t1", "t2", "t3"].forEach(t => {
        let nt = checkLengthAndParse(node[t].split("–"))
        if (nt.length !== 0) {
            node[t + "s"] = nt[0]
            node[t + "e"] = nt[1]
        }
    })
    return node
}

function checkLengthAndParse(strings) {
    if (strings.length === 0) {
        console.log("ERROR, array má velikost nula fn checkLengthAndParse().")
    } else if (strings.length === 1) {
        if (strings[0].length === 0) {
                return []
        }
        return [parseTime(strings[0]), parseTime("2030")]
    } else if (strings.length === 2) {
        if (strings[1].length === 0) {
            return [parseTime(strings[0]), parseTime("2030")]
        } 
        return [parseTime(strings[0]), parseTime(strings[1])]
    } else {
        console.log("ERROR, array má velikost >2 fn checkLengthAndParse().")
    }
}

function giveYs(nodes, nodeNames) {
    let [y, p, newNs] = [40, 180, []]
    while (true) {
        nodeNames.forEach(name => {
            let named = nodes
                .filter(n => n.nodeName === name)
                .filter(n => !(newNs.includes(n)))
            if (named.length != 0) {
                y = y + p
                newNs.push(findLine(named, y))
                newNs = newNs.flat()
            }
        })
        if (nodes.length === newNs.length) {
            break
        }

    }
    console.log("after fn giveYs() nodes have unique newIDs and node names: ",isUnique(newNs, "newID"), isUnique(newNs, "name"))
    return newNs
}

function findLine(nodes, y) {
    if (nodes.length === 1 ) {
        nodes[0].y = y 
        return nodes[0]
    }
    let line = []
    min = findLowestTime(nodes)
    nodes = nodes.filter(n => n.newID != min.newID)
    min.y = y
    line.push(min)
    if (nodes.length === 0) {
        return line
    }
    while(true) {
        min = findLowestBasedOnPrevious(nodes, line[line.length - 1])
        if (min != null) {
            min.y = y
            line.push(min)
        } else { break }
        nodes = nodes.filter(n => n.newID != min.newID)
        if (nodes.length === 0) {
            return line
        }
    }
    return line
}

function findLowestTime(nodes) {
    let min = nodes[0]
    nodes.forEach(n => {
        if (n.t1s === undefined) { console.log("ERROR, node nemá přidělený čas fn findLowestTime().", n) }
        if (n.t1s.getTime() < min.t1s.getTime()) { min = n }
    })
    return min
}

function findLowestBasedOnPrevious(nodes, previousMin) {
    let min = []
    nodes.forEach(n => {
        if (n.t1s === undefined) { console.log("ERROR, node nemá přidělený čas fn findLowestBasedOnPrevious().", n) } 
        if (highestTimeForANode(previousMin).getTime() < n.t1s.getTime()) {
            if (min.length === 0) {
                min[0] = n
            } else {
                if (n.t1s.getTime() < min[0].t1s.getTime()) {
                    min[0] = n
                }
            }
        }
    })
    if (min.length === 0 ){
        return null
    }
    return min[0]
}

function highestTimeForANode(node) {
    if (node.t1e.getTime() === new Date ("0000").getTime()) { 
        console.log(node.name, "has time of death: 0000")
    }
    if (node.t1s === undefined) { console.log("ERROR, node nemá přidělený čas fn highestTimeForNode().", node) }
    if (node.t3 !== "") { return node.t3e }
    if (node.t2 !== "") { return node.t2e }
    return node.t1e
}

/* 
<svg width="800" height="400">
  <rect x="50" y="20" width="150" height="150" style="fill:blue;opacity:0.5" />
  <path d = "M50 180 L200 180 Z" style = "stroke:pink; stroke-width:5" />
</svg>
    width = x1 - x2
*/

function rectangle(x1, y, x2, height) {
    return {
        x: x1, 
        y: y, 
        width: x2 - x1, 
        height: height, 
    }
}

function colors(nodeName) {
    switch (nodeName) {
        case "sociolog":
            return "#5C5D8D"    
        case "instituce":
            return "#654F6F"
        case "casopis":
            return "#0FA3B1"
        default:
            console.log(`Some node does not have correct nodeName, ${nodeName}`)
            return "red"
    }
}
