d3.csv("./data/nodes/Sociologove.csv").then(soc => {
    d3.csv("./data/nodes/Instituce.csv").then(ins => {
        d3.csv("./data/nodes/Casopisy.csv").then(cas => {
            d3.csv("./data/edges/Sociologove_Casopisy.csv").then(soc_cas => {
                d3.csv("./data/edges/Sociologove_Instituce.csv").then(soc_ins => {
                    d3.csv("./data/edges/Casopisy_Instituce.csv").then(cas_ins => {
                        d3.csv("./data/edges/Casopisy_Casopisy.csv").then(cas_cas => {
                            d3.csv("./data/edges/Instituce_Instituce.csv").then(ins_ins => {
                                d3.csv("./data/edges/Sociologove_Sociologove.csv").then(soc_soc => {
                                    draw({
                                        nodes: giveIDs([soc, ins, cas], ["sociolog", "instituce", "casopis"]),
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

// Příliš nepřehledné
// Možná využít nápad dr. Dvořáka se sílou.
// kdy uzly, které se vyskytovali ve stejné době by byly spojovány větší silou, nežli k těm se kterými nemohli "žit"

function draw(data) {
    console.log(d3.scaleLinear([10, 100], ["brown", "steelblue"])(50), data)
    let w = (innerWidth * 6)
    let h = (innerHeight  * 6)
    let r = 35
    let simulation = d3.forceSimulation(data.nodes)
        .force("charge", d3.forceManyBody().strength(-5080))
        .force("colide", d3.forceCollide(r + 35))
        .force("link", d3.forceLink(data.edges).distance(350))
        .force("center", d3.forceCenter(w/2, h/2))
        let cvs = setCanvas(w,h)
        let edges = cvs.append("g").selectAll("links")
            .data(data.edges)
            .enter()
            .append("line")
                .style("stroke", d => {
                    switch (d.edgeName) {
                        case "sociolog_sociolog":
                            return "lightblue"
                        case "sociolog_instituce":
                            return "lightgreen"
                        case "sociolog_casopis":
                            return "violet"
                        case "casopis_casopis":
                            return "black"
                        case "instituce_instituce":
                            return "orange"
                        case "instituce_casopis":
                            return "cyan"                            
                        default:
                            return "salmon";
                    }
                })
                .style("stroke-width", 2)
        

        let gc = cvs.append("g").selectAll("nodes")
            .data(data.nodes)
            .enter()
            .append("g").attr("class", "g.nodes")
        
        let defs = gc.append("defs");
        let gradient = defs.append("linearGradient")
            .attr("id", "gradient")
            .attr("x1", "0%")
            .attr("x2", "100%")
            .attr("y1", "0%")
            .attr("y2", "100%");
        
        gradient.append("stop")
            .attr('class', 'start')
            .attr("offset", "0%")
            .attr("stop-color", "yellow")
            .attr("stop-opacity", 1);
        
        gradient.append("stop")
            .attr('class', 'start')
            .attr("offset", "50%")
            .attr("stop-color", "red")
            .attr("stop-opacity", 1);

        gradient.append("stop")
            .attr('class', 'end')
            .attr("offset", "100%")
            .attr("stop-color", "violet")
            .attr("stop-opacity", 1);    
        // použij místo circle symboly
        /*
        let nodes = gc.append("circle")
                .attr("r", r)
                .style("fill", "url(#svgGradient)")
        */
        let switcher = things([d3.symbol().size(1950).type(d3.symbolCircle), d3.symbol().size(1950).type(d3.symbolTriangle), d3.symbol().size(1950).type(d3.symbolSquare)], ["sociolog", "instituce", "casopis"])
        let l = d3.symbol().size(950).type(d3.symbolSquare)
        console.log(switcher("sociolog"), l)
        let symbols = gc.append("path")
            .attr("d", d => { 
                return switcher(d.nodeName)()
            }) 
            .style("fill", "url(#gradient)")
       
        let labels = gc
                .append("text")
                    .text(d => d.name)
                    .style("fill", "purple")
                    .attr("text-anchor", "middle")
                    .style("font", "8px")
                    .attr("x", 0)
                    .attr("y", 0)
        symbols.append("title")
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
            // use transform to move symbols
            symbols
               .attr('transform', d => { 
                   return 'translate('+d.x+','+d.y+')'; 
                })
            /*
            nodes
                .attr("cx", d => d.x)
                .attr("cy", d => d.y)
            */
        });
}

function colors(i) {
    // TYRKYSOVÁ: #47cde7 RŮŽOVÁ: #f97c8e FIALOVÁ: #ab3595 TMAVĚ MODRÁ: #295fab
    let c = ["#53c9e0", "#ddb4ed", "#53e089", "#295fab"]
    return c[i]
    
}

function things(things, nodenames) {
    let map = new Map()
    things.forEach((t, i) => {
        map.set(nodenames[i], t)
    })
    return function(nodename) {
        return map.get(nodename)
    }
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

function dates(dates) {
    dates.forEach(element => {
        element.born = new Date(element.born)
        element.died = new Date(element.died)
    })
}

// with help of: https://github.com/scotthmurray/d3-book/releases
// https://alignedleft.com/work/d3-book-2e chapter 13. layouts : force-layout
// colors: https://duo.alexpate.uk/
// Způsob přidání gradientu https://www.freshconsulting.com/d3-js-gradients-the-easy-way/
// with help of: https://bl.ocks.org/Andrew-Reid/24a5ddaab5e2756fbc029dccc1da3f8b
// with help of: http://using-d3js.com/05_10_symbols.html
// D3 for the Impatient: Interactive Graphics for Programmers and Scientists
