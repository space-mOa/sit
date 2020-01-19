class Network {
    constructor(nodes, edges) {
        nodes.forEach(element => element.x = 1)
        nodes.forEach(element => element.y = 1)
        this.nodes = nodes
        this.edges = edges 
    }
}
d3.csv("./data/living.csv").then(lived => {
    d3.json("./data/nodes.json").then(sociologists => {
        d3.csv("./data/casopisy.csv").then(casopisy => {
            d3.csv("./data/SocJour.csv").then(socjour => {
                dates(sociologists.nodes)
                sociologists.nodes.sort((a, b) => {
                    return a.born.getFullYear() - b.born.getFullYear()
                })
                let filtredEdges = pluck(socjour, "Sociolog", 10) // Vyber položky, kde je 10 unikátních socilogů
                // console.log(filtredEdges, "\n",pick(filtredEdges, "Sociolog"), "\n", shave(pick(filtredEdges, "Sociolog"), "Sociolog")) 
                let nameOfSociologists = shave(pick(filtredEdges, "Sociolog"), "Sociolog")
                let nameOfJournals = shave(pick(filtredEdges, "Casopis"), "Casopis") 
                let filtredLived = lived.filter(element => {
                    if (nameOfSociologists.includes(element.Sociolog_1) && nameOfSociologists.includes(element.Sociolog_2)) {
                        return element
                    }
                });
                // 10*9/2 = 45
                console.log(filtredLived, pick(filtredLived, "Sociolog_1"), pick(filtredLived, "Sociolog_2"))
                // console.log("hello", filtredEdges, nameOfJournals, nameOfSociologists)
                let soc = new Network(
                    sociologists.nodes.filter(element => {
                        if (nameOfSociologists.includes(element.name)) {
                            return element
                        }
                    }),
                    filtredEdges)
                let cas = new Network(
                    casopisy.filter(element => {
                        if (nameOfJournals.includes(element.Nazev)) {
                            return element
                        }
                    }),
                    filtredEdges
                );
                sociologists = null
                socjour = null
                filtredEdges = null


            })
        })
    })
})


// pluck vezme položky na základě daného klíče, které jsou unikátní
// položky v novém poli nejsou unikátní a jsou celé
function pluck(array, selection, quantity) {
    let newArray = []
    let unique = []
    unique.push(array[0][selection])
    for (element of array) {
        let match = unique.includes(element[selection])        
        if (unique.length <= quantity) {
            if (!match) {
                unique.push(element[selection])  
                newArray.push(element)
            } 
            newArray.push(element)
        }
    }
    return newArray
}

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