// !! KÓD BYL PŘEVZAT NEBO UPRAVEN: https://observablehq.com/@d3/modifying-a-force-directed-graph?collection=@d3/d3-force
//                                  https://github.com/d3/d3-force#simulation_nodes
//                                  https://stackoverflow.com/questions/34589488/es6-immediately-invoked-arrow-function
//                                  https://observablehq.com/@d3/force-directed-graph
//                                  https://observablehq.com/@d3/sticky-force-layout?collection=@d3/d3-force

import { setHTML } from './libraries/helpers.js';

setHTML('pink')

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

const edgesSimulation = ((edges) => {   // vztahy musí být objekty obsahující atributy: index, source, target
    return edges.map(e => { return {index: 0, source: e[0], target: e[1]} })
})(edges)

let nodes = ((edges) => {   // pro uzly si můžeme stanovit nějaký atribut, který bude sloužit jako id
    return Array.from(new Set(edges.flat())).map(n => { return {id: n} })
})(edges)

const viewbox = d3.select("body")
    .append('svg')
        .attr('viewBox', [-width/2, -height/2, width, height])
        .style('border', '2.5px solid #243357')
       
const simulation = d3.forceSimulation(nodes)
    .force('charge', d3.forceManyBody())
    .force("colide", d3.forceCollide(50))
    .force('link', d3.forceLink(edgesSimulation).id(e => e.id)) // .id() určí identifikátor pro uzly, viz dokumentace links.links() a links.id()
    .force("center",  d3.forceCenter())

const svg_edges = viewbox.append('g').attr("id", "svg.edges")   // svg_edges budou za svg_nodes, protože byly vykresleny jako první
    .selectAll('edges')
    .data(edgesSimulation)
    .enter()
    .append('line')
        .attr('stroke', 'black')
        .attr('stroke-width', 0.8)

const svg_nodes = viewbox.append('g').attr("id", "svg.nodes")
    .selectAll('nodes')
    .data(nodes)
    .enter()
    .append('circle')
        .attr('cx', 0)  // na cx nezáleží, bude přepsána v simulation.on()
        .attr('cy', 0)  // na cy nezáleží, bude přepsána v simulation.on()
        .attr('r', r)
        .style('fill', 'black')

const svg_labels = viewbox.append('g').attr("id", "svg.labels")
    .selectAll('text')
    .data(nodes)
    .enter()
    .append('text')
        .text(d => d.id)
        .attr('fill', 'white')
        .attr('text-anchor', 'middle')
        .attr("font-weight", "normal")
        .attr("x", 0)   // na x nezáleží, bude přepsána v simulation.on()
        .attr("y", 0)   // na y nezáleží, bude přepsána v simulation.on()

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

console.log(
      nodes
    , simulation
)


// !! KÓD BYL PŘEVZAT NEBO UPRAVEN: https://observablehq.com/@d3/modifying-a-force-directed-graph?collection=@d3/d3-force
//                                  https://github.com/d3/d3-force#simulation_nodes
//                                  https://stackoverflow.com/questions/34589488/es6-immediately-invoked-arrow-function
//                                  https://observablehq.com/@d3/force-directed-graph
//                                  https://observablehq.com/@d3/sticky-force-layout?collection=@d3/d3-force
