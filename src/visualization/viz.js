d3.csv("./data/living.csv").then(links => {
    d3.json("./data/nodes.json").then((nodes) => {
        let network = {
            nodes: [],
            links: []
        };
        nodes.nodes.forEach(element => element.x = 1)
        nodes.nodes.forEach(element => element.y = 1)
        network.nodes = nodes.nodes;
        network.links = links;
        dates(network.nodes)
        nodes = null
        links = null
        let paper = d3.select("body")
            .append("svg")
                .attr("width", innerWidth - 200)
                .attr("height", innerHeight - 200)
            .append("g");
        paper.selectAll("circle")
            .data(network.nodes.filter(element => {
                return element.died.getFullYear() < 1950
            }))
            .enter()
            .append("circle")
                .attr("fill", "coral")
                .attr("cx", translate)
                .attr("cy", 400)
                .attr("r", 12);
    })
})

function dates(dates) {
    dates.forEach(element => {
        element.born = new Date(element.born)
        element.died = new Date(element.died)
    })
}

function translate(d, i) {
    return d.x * i * 35
}

// Navrhované rozložení (died): 
// >= 1950 (136 položek)
// < 1950 (42 položek)
// >= 2000 (61 položek)

// https://www.dashingd3js.com/svg-basic-shapes-and-d3js
// https://www.d3-graph-gallery.com/graph/arc_basic.html
// D3.js By Example - Autor: Heydt, Michael
