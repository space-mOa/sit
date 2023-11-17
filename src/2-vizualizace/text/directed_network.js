// !! KÓD BYL PŘEVZAT NEBO UPRAVEN: https://observablehq.com/@d3/modifying-a-force-directed-graph?collection=@d3/d3-force
//                                  https://github.com/d3/d3-force#simulation_nodes
//                                  https://stackoverflow.com/questions/34589488/es6-immediately-invoked-arrow-function
//                                  https://observablehq.com/@d3/force-directed-graph
//                                  https://observablehq.com/@d3/sticky-force-layout?collection=@d3/d3-force
//                                  https://observablehq.com/@d3/d3-path
//                                  https://webdesign.tutsplus.com/svg-viewport-and-viewbox-for-beginners--cms-30844t

import { setHTML } from './libraries/helpers.js';

const edges = [
    [1,2],
    [3,4],
    [2,4],
    [1,5],
    [3,2],
    [1,4]
]
const [width, height] = [800, 450]
const r = 14
const colors = {
    "black": "black",
    "cyan": "#c3f5f7",
    "grey": "#736f64",
    "white": "white"
}
let nodes = Array.from(new Set(edges.flat())).map(n => { return {id: n} })
const edgesSimulation = edges.map(e => { return {index: 0, source: e[0], target: e[1]} })



// DRAW
setHTML(colors.white);

const viewbox = d3.select("body")
    .append('svg')
        .attr('viewBox', [-width/2, -height/2, width, height])
        .style('border', `2.5px solid ${colors.black}`)

const svg_edges = viewbox.append('g').attr("id", "svg.edges")   // svg_edges budou za svg_nodes, protože byly vykresleny jako první
.selectAll('edges')
.data(edgesSimulation)
.enter()
.append('line')
    .attr('stroke', colors.black)
    .attr('stroke-width', 0.8)
    .attr('marker-end', 'url(#arrow)')

const svg_nodes = viewbox.append('g').attr("id", "svg.nodes")
.selectAll('nodes')
.data(nodes)
.enter()
.append('circle')
    .attr('cx', 0)  // na cx nezáleží, bude přepsána v simulation.on()
    .attr('cy', 0)  // na cy nezáleží, bude přepsána v simulation.on()
    .attr('r', r)
    .style('fill', colors.black)

const svg_labels = viewbox.append('g').attr("id", "svg.labels")
.selectAll('text')
.data(nodes)
.enter()
.append('text')
    .text(d => d.id)
    .attr('fill', colors.white)
    .attr('text-anchor', 'middle')
    .attr("font-weight", "normal")
    .attr("x", 0)   // na x nezáleží, bude přepsána v simulation.on()
    .attr("y", 0)   // na y nezáleží, bude přepsána v simulation.on()

const simulation = d3.forceSimulation(nodes)
    .force('charge', d3.forceManyBody())
    .force("colide", d3.forceCollide(50))
    .force('link', d3.forceLink(edgesSimulation).id(e => e.id)) // .id() určí identifikátor pro uzly, viz dokumentace links.links() a links.id()
    .force("center",  d3.forceCenter())

simulation.on("tick", () => {   // updatuje souřadnice každý tick 
svg_edges
    .attr("x1", d => d.source.x)
    .attr("y1", d => d.source.y)
    .attr("x2", d => d.target.x) 
    .attr("y2", d => d.target.y)
svg_nodes
    .attr("cx", d => d.x)
    .attr("cy", d => d.y)
svg_labels
    .attr('x', d => d.x - 0.26)
    .attr('y', d => d.y + 5)
});    

const svg_markers = viewbox.append('defs').append('marker')
    .attr('id', 'arrow')
    .attr('orient', 'auto')
    .attr('refX', r + r -1)
    .attr('refY', 5)
    .attr('markerWidth', 30)
    .attr('markerHeight', 38)
    .attr('viewBox', '0 0 30 30')
    .attr('fill', colors.black)
    .append('path')
        .attr('d', 'M 0 0 L 10 5 L 0 10 z')

console.log(
  nodes
, simulation
, svg_markers
)


// !! KÓD BYL PŘEVZAT NEBO UPRAVEN: https://observablehq.com/@d3/modifying-a-force-directed-graph?collection=@d3/d3-force
//                                  https://github.com/d3/d3-force#simulation_nodes
//                                  https://stackoverflow.com/questions/34589488/es6-immediately-invoked-arrow-function
//                                  https://observablehq.com/@d3/force-directed-graph
//                                  https://observablehq.com/@d3/sticky-force-layout?collection=@d3/d3-force
