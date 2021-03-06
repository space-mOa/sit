class Network {
    constructor(nodes, edges) {
        nodes.forEach(element => element.x = 1)
        nodes.forEach(element => element.y = 1)
        this.nodes = nodes
        this.edges = edges 
    }
}

d3.csv("./data/SocJour.csv").then(socjour => {
    d3.json("./data/nodes.json").then(sociologists => {
            d3.csv("./data/casopisy.csv").then(casopisy => {
                dates(sociologists.nodes)
                sociologists.nodes.sort((a, b) => {
                    return a.born.getFullYear() - b.born.getFullYear()
                })
                socjour.sort((a, b) => {
                    return parseInt(a.Sociolog_ID) - parseInt(b.Sociolog_ID)
                })
                let filtredEdges = pluck(socjour, "Sociolog", 30) 
                console.log(filtredEdges)
                let nameOfJournals = shave(pick(filtredEdges, "Casopis"), "Casopis") 
                let nameOfSociologists = shave(pick(filtredEdges, "Sociolog"), "Sociolog")
                let soc = new Network(
                    sociologists.nodes.filter(element => {
                        if (nameOfSociologists.includes(element.name)) {
                            return element
                        }
                    }),
                    filtredEdges)
                let cas = new Network(
                    casopisy.filter(element => {
                        if (nameOfJournals.includes(element.Nazev)) {
                            return element
                        }
                    }),
                    filtredEdges
                );
                // nejvíce hran u Časopisů
                u = []
                nameOfJournals.forEach(n => {
                    let p = filtredEdges.filter(e => n == e.Casopis)
                    u.push([p.length -1, n])
                })
                let maxJour = u[0]
                u.forEach(e => {
                    if (maxJour[0] < e[0]) {
                        maxJour = e
                    }
                })
                // nejvíc hran Sociologů
                u = []
                nameOfSociologists.forEach(n => {
                    let p =  filtredEdges.filter(e => n == e.Sociolog)
                    u.push([p.length, n])
                })
                let maxSoc = u[0]
                u.forEach(e => {
                    if (maxSoc[0] < e[0]) {
                        maxSoc = e
                    }
                })
                sociologists = null
                socjour = null
                filtredEdges = null
                let paper = d3.select("body")
                    .append("svg")
                        .attr("width", innerWidth * 0.95)
                        .attr("height", 2300)
                        .append("g")
                            .attr("transform", "translate(0, 80)")
                            .attr("render-order", "1");
                paper.selectAll("sociologist")
                    .data(soc.nodes)
                    .enter()
                    .append("circle")
                        .attr("fill", "DarkSeaGreen")
                        .attr("cx", 350)
                        .attr("cy", yAxis)
                        .attr("r", 20)
                        .attr("transform", "translate(0, 110)")
                paper.selectAll("names")
                    .data(soc.nodes)
                    .enter()
                    .append("text")
                        .text((d) => {
                            return d.name
                        })
                        .attr("fill", "DarkSeaGreen")
                        .attr("text-anchor", "middle")
                        .attr("x", 350)
                        .attr("y", d => d.y)
                        .style("font", "8px")
                            .attr("transform", "translate(-125,110)");
                paper.selectAll("dates")
                        .data(soc.nodes)
                        .enter()
                        .append("text")
                            .text((d) => {
                                if (d.died.getFullYear() == 0) {
                                    return d.born.getFullYear() + "  -  " + "???"                                    
                                }
                                if (d.died.getFullYear() == 2030) {
                                    return d.born.getFullYear() + "  -  " + ""
                                }
                                return d.born.getFullYear() + "  -  " + d.died.getFullYear()
                            })
                            .attr("fill", "DarkSeaGreen")
                            .attr("text-anchor", "middle")
                            .attr("x", 350)
                            .attr("y", d => d.y)
                            .style("font-size", "12px")
                                .attr("transform", "translate(-125,129)");
                paper.selectAll("journal")
                    .data(cas.nodes)
                    .enter()
                    .append("circle")
                        .attr("fill", "LightSalmon")
                        .attr("cx", 1100)
                        .attr("cy", yAxis)
                        .attr("r", 20)
                paper.selectAll("titles")
                    .data(cas.nodes)
                    .enter()
                    .append("text")
                        .text((d) => {
                            return d.Nazev
                        })
                        .attr("fill", "LightSalmon")
                        .attr("text-anchor", "left")
                        .attr("x", 350)
                        .attr("y", d => d.y)
                        .style("font", "8px")
                            .attr("transform", "translate(800, 6)");         
                paper.selectAll("link")
                    .data(soc.edges)
                    .enter()
                    .append("line")
                        .attr("x1", 350)
                        .attr("y1", d => {
                            let source = soc.nodes.find(element => parseInt(element.id) == parseInt(d.Sociolog_ID))
                            return source.y + 110
                        })
                        .attr("x2", 1100)
                        .attr("y2", d => {
                            let target = cas.nodes.find(element => parseInt(element.index) == parseInt(d.Casopis_ID))
                            return target.y
                        })
                        .attr("stroke", d => {
                            let color = "lightgrey"
                            if (d.Casopis == maxJour[1]) {
                                color = "LightSalmon"
                            }
                            if (d.Sociolog == maxSoc[1]) {
                                color = "DarkSeaGreen"
                            }
                            return color
                        })
                        .attr("stroke-width", d => {
                            let width = "1.2"
                            if (d.Casopis == maxJour[1]) {
                                width = "2"
                            }
                            if (d.Sociolog == maxSoc[1]) {
                                width = "2"
                            }
                            return width})
        })  
    })   
})

// pluck vezme položky na základě daného klíče, které jsou unikátní + počet
// řádek v tabulce: KlíčA, KlíčC -> vrátí řádek např. 30 uníkátních autorů má 120 položek
function pluck(array, selection, quantity) {
    let newArray = []
    let unique = []
    unique.push(array[0][selection])
    for (element of array) {
        let match = unique.includes(element[selection])        
        if (unique.length <= quantity) {
            if (!match) {
                unique.push(element[selection])  
                newArray.push(element)
            } else {
                newArray.push(element)
            }
        }
    }
    return newArray
}

// takeUnique vezme položky na základě daného klíče, které josu unikátní
// řádek v tabulce: KlíčA, KlíčC -> vrátí řádek dle vybraného klíče
function pick(array, selection) {
    let newArray = []
    let unique = []
    unique.push(array[0][selection])
    newArray.push(array[0])
    for (element of array) {
        let match = unique.includes(element[selection])
        if (!match) {
            unique.push(element[selection])
            newArray.push(element)
        }
    }
    return newArray
}

// shave vezme položky na základě daného klíče
// vybyere jen ty položky na základě klíče
// řádek v tabulce: KlíčA, KlíčC -> vrátí jen položku dle vybraného klíče
function shave(array, selection) {
    let newArray = []
    array.forEach(element => {
        newArray.push(element[selection])
    })
    return newArray
}

function dates(dates) {
    dates.forEach(element => {
        element.born = new Date(element.born)
        element.died = new Date(element.died)
    })
}

function yAxis(d, i) {
    return d.y = d.y * i * 65
}

// D3.js By Example - Autor: Heydt, Michael
// https://www.dashingd3js.com/svg-basic-shapes-and-d3js
// https://www.d3-graph-gallery.com/graph/arc_basic.html
// https://www.tutorialspoint.com/d3js/d3js_svg_transformation.htm
// https://www.npmjs.com/package/d3-transform
// https://www.dashingd3js.com/svg-paths-and-d3js
// https://www.d3-graph-gallery.com/graph/shape.html