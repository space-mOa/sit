// !! KÓD BYL PŘEVZAT NEBO UPRAVEN: https://d3js.org/d3-hierarchy/tree
//                                  https://observablehq.com/@d3/tree
//                                  https://observablehq.com/@d3/tree/2
//                                  https://observablehq.com/@d3/d3-hierarchy
//                                  https://observablehq.com/@d3/d3-stratify
//                                  https://d3js.org/d3-hierarchy/stratify
//                                  https://observablehq.com/@d3/radial-tree

import { setHTML } from './libraries/helpers.js';

// let edges = [
//     ["", 1],
//     [1, 2],
//     [1, 3]
// ].map(d => ({ "parent": d[0], "child": d[1] }))
// edges.columns = ["parent", "child"]
let edges = [
    ["", 1],
    [1, 2],
    [1, 3],
    [2, 4],

]
const labels = Array.from(new Set(edges.flat().filter(d => d != '')))
const [width, height] = [800, 450]
const r = 14
const colors = {
    "black": "black",
    "purple": "#e1d5f7",
    "grey": "#736f64",
    "white": "white"
}
const settings = {
    "stroke_width": "2px",
}
let root = d3.stratify()
        .id(d => d[1])
        .parentId(d => d[0])
        (edges)
d3.tree().size([width/4, height/2])(root) //[width/4, height/1.5])(root)
//console.log(edges, tree_layout, tree_layout.links())

// const s = d3.stratify()
//         .id(d => d[edges.columns[1]])
//         .parentId(d => d[edges.columns[0]])
//         (edges)
// const s = d3.stratify()
//         .id(d => d.child)
//         .parentId(d => d.parent)
//         (edges)
// console.log(s)

// DRAW
setHTML(colors.white);
const svg = d3.select("body")
        .append('svg')
            .attr('viewBox', [-width/2.5, -50, width, height])
            .style('border', `2.5px solid ${colors.black}`)

const links = svg.append("g")
        .attr("id", "edges")
        .attr("fill", "none")
        .attr("stroke",colors.black)
        .selectAll()
            .data(root.links())
            .join("path")
            .attr("d", d3.linkVertical()
                .x(d => d.x)
                .y(d => d.y));

const node = svg.append("g")
        .attr("stroke-linejoin", "round")
        .attr("stroke-width", 3)
        .selectAll()
            .data(root)
            .join("g")
                .attr("transform", d => `translate(${d.x},${d.y})`);
    
node.append("circle")
    .attr("fill", colors.black)
    .attr("r", r);

node.append("text")
    .style("font", 5)
    .attr("x", d => d.children ? 4.5 : -4.5)
    .attr("y", d => d.children ? 5 : 5)
    .attr("text-anchor", d => d.children ? "end" : "start")
    .attr("font-weight", "normal")
    .text(d => d.id)
        .attr("fill", colors.white)


console.log(root.descendants())
console.log(labels)