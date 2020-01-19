class Network {
    constructor(nodes, edges) {
        nodes.forEach(element => element.x = 1)
        nodes.forEach(element => element.y = 1)
        this.nodes = nodes
        this.edges = edges 
    }
}
d3.csv("./data/living.csv").then(living => {
    d3.json("./data/nodes.json").then(sociologists => {
        dates(sociologists.nodes)
        
        
    })
})


// takeUnique vezme položky na základě daného klíče, které josu unikátní
// položky v novém poli jsou unikátní a jsou celé
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