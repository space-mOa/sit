// !! KÓD BYL PŘEVZAT A UPRAVEN: https://github.com/d3/d3-scale
//                               https://observablehq.com/@d3/line-chart
//                               https://stackoverflow.com/questions/3746725/how-to-create-an-array-containing-1-n
//                               https://observablehq.com/@d3/styled-axes?collection=@d3/d3-axis

import { loadCSV, nameCSV, setHTML, parseDates, checkDateType } from './libraries/helpers.js';
import { SociologistCustomDates } from './libraries/custom_dates.js'


const CVSs = [
    './data/SCSg/Sociologove.csv'
]

const transformers = [
    (row) => { return {id: row.id, name: row.name, born: parseDates(row.born, SociologistCustomDates, row.name), died: parseDates(row.died, SociologistCustomDates, row.name)} }
]

const names = [
    'Sociologove'
]

const container =  {
        width: innerWidth * 2/3 + innerWidth,
        height: innerHeight * 2/3,
        createContainer() { return d3.select('body').append('svg').attr('width', this.width).attr('height', this.height) }//.attr("viewBox", [0, 0, this.width, this.height]) }
}

const scale = {
        from: checkDateType(new Date (1800)),
        to: checkDateType(new Date (2030)),
        margin: 20,
        createTimeScale(fromScreen, toScreen) {
            return d3.scaleTime()
                .domain([this.from, this.to])
                .range([fromScreen + this.margin, toScreen - this.margin])
        }
}

const axis = {
    orientation: d3.axisBottom,
    count: 20,
    style: {
            stroke_width: 5,
            stroke_linecap: 'round'
    },
    createAxis(scale, from, to) { return this.orientation().scale(scale).ticks(this.count) }
}

const rectangles = {
    height: 10,
    margin: 3,
    start: 10,
    step: 10,
    createRectangles(data, scale) {
        let y = this.start
        return {Sociologove: data.map((d,i) => {
            y = y + this.step
            const info = {id: d.id, name: d.name, born: d.born, died: d.died}
            const rectangle = {x: scale(d.born), y: y, width: scale(d.born) + scale(d.died), height: this.height}
            return {Sociolog: info, rectangle: rectangle}
        })}
    }
}

const plot = {
    container: container,
    scale: scale,
    axis: axis,
    rectangles: rectangles
};


setHTML('white')
Promise.all(loadCSV(CVSs, transformers)).then(f => prepare(nameCSV(names,f), plot))

function prepare(data, plot) {
    console.log(data)
    const container  = plot.container.createContainer()
    const scale      = plot.scale.createTimeScale(0, plot.container.width)
    const axis       = plot.axis.createAxis(scale, plot.scale.from, plot.scale.to)
    const rectangles = plot.rectangles.createRectangles(data.Sociologove, scale)
    console.log(scale)
    
    container.append('g')
        .call(axis)
        .attr('stroke-width', plot.axis.style.stroke_width)
        .attr('stroke-linecap', plot.axis.style.stroke_linecap)
}

function draw(svg, plot) {

}