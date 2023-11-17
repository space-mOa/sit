// !! KÓD BYL PŘEVZAT A UPRAVEN: https://observablehq.com/@d3/tidy-tree?collection=@d3/d3-hierarchy

let f = [
    "../../data/VSgS/tree.csv",
]

let n = [
    "tree",
]

function load(fs) {
    const p = []
    for (f of fs) {
        p.push(d3.csv(f))
    }
    return p
}

function name(ns,v) {
    return ns.map((n, i) => {
        return {data: v[i], name: n}
    });
}


function makeHierarchy(table) {
    let nd = d3.stratify()
        .id(d => { return d.child})
        .parentId(d => {return d.parent})
        (table); // returns function that is immediately called
    if (nd instanceof d3.hierarchy) {
        return nd
    }
    throw "Function makeTree(table): couldn't make a tree."
}

function setCanvas(w,h) {
    return canvas = d3.select("body")
    .append("svg")
        .attr("width", w)
        .attr("height", h)
}

function makeTreeLayout(hierarchy, w, h) {
    return d3.tree()
        .size([w,h])
        (hierarchy);
}

function drawTree(canvas, l, a) {
    const g = canvas.append("g")
        .attr("font-family", "sans-serif")
        .attr("font-size", 10)
        .attr("transform", `translate(${100},${0})`);

    const links = g.append("g")
        .attr("fill", "none")
        .attr("stroke", a.link.stroke)
        .attr("width", a.link.width)
    .selectAll("path")
        .data(l.links())
        .join("path")
            .attr("d", 
                d3.linkHorizontal()
                    .x(d => d.y)
                    .y(d => d.x));
    
    const nodes = g.append("g")
        .attr("stroke-linejoin", a.link.stroke_linejoin)
        .attr("stroke-width", a.link.stroke_width)
    .selectAll("g")
        .data(l.descendants())
        .join("g")
            .attr("transform", d => `translate(${d.y},${d.x})`);
    
    nodes.append("circle")
            .attr("fill", d => a.node.choose_color(d.children))
            .attr("r", a.node.r);  

    nodes.append("text")
        .style("font", a.text.size)
        .attr("x", d => d.children ? -6 : 6)
        .attr("text-anchor", d => d.children ? "end" : "start")
        .text(d => d.id)
            .attr("fill", a.text.fill)
}

setHTML("white")
Promise.all(load(f)).then(v => draw(name(n, v)))

function draw(d) {
    let w = innerWidth
    let h = innerHeight + (innerHeight * 1/4)
    const c = {
        black:  "#222222",
        green:  "#227843",
        orange: "#F8A646",
        white:  "#FBF9F9",
    }
    const a = {
            link: {
                stroke:             c.black,
                width:              1.5,
                stroke_linejoin:    "round",
                stroke_width:       3.5,
            },
            node: {
                r:                  5,
                choose_color:       (t) => t ? c.orange : c.green,
            },
            text: {
                fill: c.black,
                size: 14,
            } 

    }
    let root = makeHierarchy(d[0].data)
    drawTree(
        setCanvas(w, h), 
        makeTreeLayout(root, w / 2, h), 
        a
    )
}

function setHTML(b) {
    document.body.style.backgroundColor = b
}

// !! KÓD BYL PŘEVZAT A UPRAVEN: https://observablehq.com/@d3/tidy-tree?collection=@d3/d3-hierarchy

// d3 .join:                https://observablehq.com/@d3/selection-join
// d3 layouts:              https://www.d3indepth.com/layouts/
// d3.stratify():           https://github.com/d3/d3-hierarchy/blob/v2.0.0/README.md#stratify
// d3.tree:                 https://observablehq.com/@d3/tidy-tree?collection=@d3/d3-hierarchy
// js, vysvětlení o.m()():  https://stackoverflow.com/questions/18234491/two-sets-of-parentheses-after-function-call
// js, body a stylování:    https://www.w3schools.com/jsref/prop_doc_body.asp