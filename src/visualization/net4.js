d3.csv("./data/nodes/Sociologove.csv").then(soc => {
    d3.csv("./data/nodes/Instituce.csv").then(ins => {
        d3.csv("./data/edges/Sociologove_Instituce-1800-1948.csv").then(time => {           
            draw({
                nodes: giveIDs([soc, ins], ["sociolog", "instituce"]), 
                edges: makeEdges(giveIDs([soc, ins], ["sociolog", "instituce"]), [time])
            })
        })
    })
})


function draw(data) {    
    console.log(data, innerWidth, data.edges[0].atr)
    let w = (innerWidth * 3)
    let h = (innerHeight  * 3)
    let r = 35
    let simulation = d3.forceSimulation(data.nodes)
        .force("charge", d3.forceManyBody().strength(-380))
        .force("colide", d3.forceCollide(r + 35))
        .force("link", d3.forceLink(data.edges).distance(350))
        .force("center", d3.forceCenter(w/2, h/2))
    let cvs = setCanvas(w,h)
    let edges = cvs.append("g").selectAll("links")
        .data(data.edges)
        .enter()
        .append("line")
            .style("stroke", d => colors(2))
            .style("stroke-width", 2)
    let gc = cvs.append("g").selectAll("nodes")
        .data(data.nodes)
        .enter()
        .append("g").attr("class", "g.nodes")
    let nodes = gc.append("circle")
            .attr("r", r)
            .style("fill", (d,i) => (d.nodeName === "sociolog") ? colors(1) : colors(0))
    let labels = gc
            .append("text")
                .text(d => d.name)
                .attr("fill", d => (d.nodeName === "sociolog") ? "#022b09" : "#42310a")
                .attr("text-anchor", "middle")
                .style("font", "8px")
                .attr("x", 0)
                .attr("y", 0)
    nodes.append("title")
            .text(d => d.name)
    simulation.on("tick", () => {
        edges
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y)
        labels
            .attr("x", d => d.x)
            .attr("y", d => d.y)
        nodes
            .attr("cx", d => d.x)
            .attr("cy", d => d.y)
    });
}


function colors(i) {
    // TYRKYSOVÁ: #47cde7 RŮŽOVÁ: #f97c8e FIALOVÁ: #ab3595 TMAVĚ MODRÁ: #295fab
    let c = ["#47cde7", "#f97c8e", "#ab3595", "#295fab"]
    return c[i]
}

// Přiřadí nové ID pro uzly, které spojí dohromady 
// + přiřadí jméno např. všechny UZLY reprezenrující časopisy budou mít atribut nodeName: časopisy
function giveIDs(nodes, ids) {
    let newID = 0
    nodes.forEach((e, i) => {
        e.forEach(n => {
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

// with help of: https://github.com/scotthmurray/d3-book/releases
// https://alignedleft.com/work/d3-book-2e chapter 13. layouts : force-layout
// colors: https://duo.alexpate.uk/