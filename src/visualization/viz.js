class Network {
    constructor(nodes, edges) {
        nodes.forEach(element => element.x = 0)
        nodes.forEach(element => element.y = 0)
        this.nodes = []
        this.nodes.push(nodes)
        this.edges = []
        this.edges.push(edges)
    }
    addEdges(edges) {
        this.edges.push(edges)
    }
    addNodes(nodes) {
        nodes.forEach(element => element.x = 0)
        nodes.forEach(element => element.y = 0)
        this.nodes.push(nodes)
    }
}

d3.csv("./data/living.csv").then(lived => {
    d3.json("./data/nodes.json").then((soc) => {
        let selection = dates(soc.nodes)
        .sort((a, b) => {
            return a.born.getFullYear() - b.born.getFullYear()
        })
        .slice(0, 20)
        let net = new Network(
            selection,
                lived)
        soc.nodes = null; lived = null; selection = null; 
        let paper = d3.select("body")
            .append("svg")
                .attr("width", innerWidth  - 25)
                .attr("height", innerHeight - 25)
            .append("g")
        let m = Math.min(...getRow(net.nodes[0], "died").map(element => element = element.getFullYear())) // Rok prvně narozeného
        console.log(m, net.nodes[0][1].died.getFullYear(), net, "☔️")
        years(net.nodes[0])
        
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
    let newArray = []
    array.forEach(element => {
        newArray.push(element[key])
    })
    return newArray
}

// nefunguje je pokud má tabulka prázdný rádek
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

// D3.js By Example - Autor: Heydt, Michael
// https://www.dashingd3js.com/svg-basic-shapes-and-d3js
// https://www.d3-graph-gallery.com/graph/arc_basic.html
// https://www.tutorialspoint.com/d3js/d3js_svg_transformation.htm
// https://www.npmjs.com/package/d3-transform
// https://www.dashingd3js.com/svg-paths-and-d3js
