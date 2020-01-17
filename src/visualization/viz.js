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
        network.nodes.sort((a, b) => {
            return a.born.getFullYear() - b.born.getFullYear()
            })
        nodes = null
        links = null
        console.log(network)
        let paper = d3.select("body")
            .append("svg")
                .attr("width", innerWidth * 0.95)
                .attr("height", innerHeight * 0.95)
                .append("g")
                    .attr("transform", "translate(135, 90)");
        paper.selectAll("circle")
            .data(network.nodes.filter(element => {
                return element.died.getFullYear() < 1940
            }))
            .enter()
            .append("circle")
                .attr("fill", "YellowGreen")
                .attr("cx", xAxis)
                .attr("cy", 400)
                .attr("r", 18);       
        paper.selectAll("text")
            .data(network.nodes.filter(element => {
                return element.died.getFullYear() < 1940
            }))
            .enter()
            .append("text")
                .text((d) => {
                    return d.id
                })
                .attr("fill", "Black")
                .attr("text-anchor", "middle")
                .attr("x", d => {return d.x})
                .attr("y", 405)
                .style("font", "8px")
        // vyfiltruj linky jen pro ty nodes, co máš
        let forties = []
        network.nodes.forEach(element => {
            if (element.died.getFullYear() < 1940) {
                forties.push(element.id)
                
            }
        })
        let nodesForties = network.nodes.filter(element => {
            if (element.died.getFullYear() < 1940) {
                return element 
            }
        })
        let linksForties = network.links.filter(element => {
            if (forties.includes(parseInt(element.Sociolog_1_ID)) && forties.includes(parseInt(element.Sociolog_2_ID))) {
                return element
            }
        })
        // https://vanseodesign.com/web-design/svg-paths-curve-commands/
        console.log(linksForties, nodesForties)
        paper.selectAll("link")
                .data(linksForties)
                .enter()
                    .append("path")
                    .attr("d", d => {
                        let source = nodesForties.find(element => {                        
                            return parseInt(element.id) == parseInt(d.Sociolog_1_ID)
                        })
                        let target = nodesForties.find(element => {                        
                            return parseInt(element.id) == parseInt(d.Sociolog_2_ID)
                        })
                        console.log(source.x, target.x)
                        return "M " + source.x + " "+ "405" + " Q1000,-350 " + target.x + " " + "405"
                        //return "M100,200 Q250,100 400,200"
                    })
                    .attr("stroke", "black")
                    .attr("fill-opacity", "0")
    })
})

function dates(dates) {
    dates.forEach(element => {
        element.born = new Date(element.born)
        element.died = new Date(element.died)
    })
}

function xAxis(d, i) {
    return d.x = d.x * i * 65
}

function modify(nodes) {
    nodes.sort((a, b) => {
        return a.born.getFullYear() - b.born.getFullYear()
    }).forEach(element => console.log(element.name, element.born.getFullYear()))
    console.log(nodes)
}

// Navrhované rozložení (died): 
// >= 1950 (136 položek)
// < 1950 (42 položek)
// >= 2000 (61 položek)

// D3.js By Example - Autor: Heydt, Michael
// https://www.dashingd3js.com/svg-basic-shapes-and-d3js
// https://www.d3-graph-gallery.com/graph/arc_basic.html
// https://www.tutorialspoint.com/d3js/d3js_svg_transformation.htm
// https://www.npmjs.com/package/d3-transform
// https://www.dashingd3js.com/svg-paths-and-d3js
