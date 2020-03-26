class Network {
    constructor(nodes, edges) {
        nodes.forEach(element => element.x = 1)
        nodes.forEach(element => element.y = 1)
        this.nodes = []
        this.nodes.push(nodes)
        this.edges = []
        this.edges.push(edges)
    }
    addEdges(edges) {
        this.edges.push(edges)  
    }
    addNodes(nodes) {
        nodes.forEach(element => element.x = 1)
        nodes.forEach(element => element.y = 1)
        this.nodes.push(nodes)
    }
    setX(number, i) {
        this.nodes[i].forEach((e,i) => {
            e.x = e.x * i * number
        })
    }
    setY(number, i) {
        this.nodes[i].forEach(e => {
            e.y = number
        })
    }
    getX(input, key, i) {
        let n = this.nodes[i].find(f => {
            if (f[key] === input) {
                return f
            }
        })
        return n.x
    }
    getY(input, key, i) {
        let n = this.nodes[i].find(f => {
            if (f[key] === input) {
                return f
            }
        })
        return n.y
    }
}

d3.csv("./data/SocJour.csv").then(SocJour => {
    d3.csv("./data/living.csv").then(lived => {
        d3.json("./data/nodes.json").then((soc) => {
            // PREPARE
            let selection = dates(soc.nodes)
                .sort((a, b) => {
                    return a.born.getFullYear() - b.born.getFullYear()
                })
                .slice(0, 70)
            let net = new Network(
                selection,
                    lived)
            net.addEdges(SocJour)
            soc.nodes = null; lived = null; selection = null; SocJour = null;
            let journals = [] // Array s objekty, kdy 1 objekt = jeden časopis. Objekt obsahuje všechny soc. se kterými souvisí 
            getRow(pick(net.edges[1], "Casopis"), "Casopis")    // Najdi unikátní časopisy
                .forEach(e => {
                    journals.push(net.edges[1].filter(f => {  // Získej všechny sociology pro Časipis A
                        if (f.Casopis === e) {
                            return f
                        }
                    }))
                })
            let eds = []; // Array s objekty {casopisy, sociologové{kteří souvisí s casopisem a zároveň spolu žili}}
            let nds = []; // Array se všemi sociology dle eds
            journals.forEach(j => {
                // console.log("|", j[0].Casopis, "|")
                let socs = []
                getEdges(j, j, net.edges[0], socs) // Pokud vrátí pro časopis A prázdny array znamená, to že je s časopisem spojen jen jeden Sociolog B, socs se mění
                if(!isEmpty(socs)) {
                    eds.push({journal: j[0].Casopis, soc: socs})
                    if (isEmpty(nds)) {
                        nds.push(socs[0].Sociolog_1)
                        nds.push(socs[0].Sociolog_2)
                    }
                    socs.forEach(e => {
                        let m = nds.find(n => {
                            if (n === e.Sociolog_1) {
                                return n
                            }
                        })
                        let n = nds.find(n => {
                            if (n === e.Sociolog_2) {
                                return n
                            }
                        })
                        if (m === undefined) {
                            nds.push(e.Sociolog_1)
                        }
                        if (n === undefined) {
                            nds.push(e.Sociolog_2)
                        }
                        
                    })
                }
            })
            net.addEdges(eds)
            nds.sort()
            net.addNodes(net.nodes[0].filter(e => { // jen ty sociologové, co jsou v nds
                if (nds.includes(e.name)) {
                    return e
                }
            }))
            net.setY(700, 1)
            net.setX(45, 1)
            console.log(net,"🥽")
            console.log(net.edges[2],"Vypsaní sociologové:", net.nodes[1], net.nodes[1][0].name, getRow(net.nodes[1], "name").includes("Durdík Josef"))
            net.addEdges(net.edges[2].map(e => { // vyber jen ty hrany, pro soc. kteří jsou vykreslení, net.nodes[1]
                let a = e.soc.filter(element => {
                    if (getRow(net.nodes[1], "name").includes(element.Sociolog_1) && getRow(net.nodes[1], "name").includes(element.Sociolog_2)) {
                        return element
                    }
                })
                return {journal: e.journal, soc: a}
            }).filter(f => {
                if (f.soc.length != 0) {
                    return f
                }
            }))
            // DRAW
            let paper = d3.select("body")
                .append("svg")
                    .attr("width", innerWidth - 25)
                    .attr("height", innerHeight - 25)
                .append("g")
                    .attr("transform", "translate(80,0)")
            paper.selectAll("Sociologist-name")
                .data(net.nodes[1])
                .enter()
                .append("text")
                    .text(d => {
                        return d.name
                    })
                    .attr("fill", "darkblue")
                    .attr("text-anchor", "middle")
                    .attr("x", d => d.x)
                    .attr("y", d => d.y)
                    .style("font", "12px")
                    .attr("transform", d => {                        
                        return "rotate("+ "290" + "," + d.x + "," + d.y +")" // šikmé jména
                    })
            /* jen ukazka vyběru */
            let y = net.nodes[1][0].y - 200
            let colors = ["asdf", "red", "sadf", "green", "asdf", "asdf", "sdfsdf", "sdfs", "magenta", "brown"]
            net.edges[3].forEach((e, i) => {
                console.log(y)
                if (i == 1 || i == 3 || i == 9 || i == 8) {
                y = y - 250
                paper.selectAll("link")
                    .data(e.soc)
                    .enter()
                    .append("path")
                        .attr("d", d => {
                            let x1 = net.getX(d.Sociolog_1, "name", 1) + 10
                            let x2 = net.getX(d.Sociolog_2, "name", 1) + 10
                            let y1 = net.getY(d.Sociolog_1, "name", 1) - 60
                            let y2 = net.getY(d.Sociolog_1, "name", 1) - 60
                            let mid = 0
                            if (x1 > x2) {
                                mid = (x1 - x2) / 2
                                mid = x2 + mid
                            } else {
                                mid = (x2 - x1) / 2
                                mid = x1 + mid
                            }
                            //console.log(x1, x2, mid, y)
                            let m = "M " + x1 + " " + y1 + " Q " + mid + " " + y + " " + x2 + " " + y2;
                            return m
                        })
                        .attr("stroke", colors[i])
                        .attr("fill-opacity", "0")
                }
            })
            let sc = ["red", "green", "magenta", "brown"]
            let xsc = [80, 300, 640, 950]
            paper.selectAll("journals")
                .data(["Statistický obzor",  "Sociální pracovnice", "Časopis Svobodné školy politických nauk v Praze", "Athenaeum"])
                .enter()
                .append("text")
                    .text(d => {return d})
                    .attr("fill", (d, i) => {
                        return sc[i]
                    })
                    .attr("text-anchor", "middle")
                    .attr("x", (d,i) => { return xsc[i]})
                    .attr("y", 850)
                    .style("font", "12px")
                    .attr("transform", "translate(80,0)")
        })
    })
})

function dates(dates) {
    dates.forEach(element => {
        element.born = new Date(element.born)
        element.died = new Date(element.died)
    })
    return dates
}

// bere pole s objekty a vybere položky na základě klíče -> vrací jen sloupec newArray
function getRow(array, key) {
    let newArray = [];
    array.forEach(element => {
        newArray.push(element[key])
    })
    return newArray
}

function pick(array, selection) {
    let newArray = [];
    let unique = [];
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

function isEmpty(array) {
    if (array.length === 0) {
        return true
    } else {
        return false
    }
}

// Vráti pole se sociology, kteří spolu žili a zároveň jsou spojeni se stejným časopisem A
function getEdges(array1, array2, array3, newArray) {
    if (!(array1.length ===  1)) {
        let s1;
        (array1.length >= 1) ? s1 = array1[array1.length - 1].Sociolog : s1 = array1[array1.length].Sociolog; // Sociolog
        array2 = array2.slice(0, array2.length - 1)
        // console.log(s1)
        findMatch(s1, array2, array3, newArray)
        getEdges(array1.slice(0, array1.length - 1), array2, array3, newArray)
    }
}

function findMatch(s1, array, searchIn, newArray) {
    let s2 = array[array.length - 1].Sociolog
    // console.log("------", s2)
    let match = searchIn.find(e => {
        if ((e.Sociolog_1 === s1 && s2 === e.Sociolog_2) || (e.Sociolog_1 === s2 && s1 === e.Sociolog_2)) {
            return e
        }
    })
    if (!(match === undefined)) {
        newArray.push(match)
    }
    if (array.length === 1 ) { 
        return; 
    } else {
        findMatch(s1, array.slice(0, array.length -1), searchIn, newArray)
    }
}
// D3.js By Example - Autor: Heydt, Michael
// https://www.dashingd3js.com/svg-basic-shapes-and-d3js
// https://www.d3-graph-gallery.com/graph/arc_basic.html
// https://www.tutorialspoint.com/d3js/d3js_svg_transformation.htm
// https://www.npmjs.com/package/d3-transform
// https://www.dashingd3js.com/svg-paths-and-d3js
// http://apex.infogridpacific.com/SVG/svg-tutorial-lesson18-path-quadratics.html
// https://www.w3schools.com/js/js_random.asp