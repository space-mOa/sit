// ! KÓD BYL PŘEVZAT A UPRAVEN: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

export { loadCSV, nameCSV, setHTML, parseDates, checkDateType, makeObjectFrom2Arrays, trimString };

function makeObjectFrom2Arrays(keys, values) {
    if (keys.length !== values.length) {
        throw new Error ("keys and values must have same length")
    }
    const object = {}
    keys.forEach((k, i) => object[k] = values[i])
    return object
}

// Odstraní mezery a převede string na lowercase 
function trimString(s) {
    return s.replaceAll(' ', '').toLowerCase()
}

function setHTML(b) {
    document.body.style.backgroundColor = b
}

// Nahraje data z CSV souborů
// Vrací Promise
function loadCSV(files, translate) {
    return files.map((f,i) => d3.csv(f, translate[i]))
}

function nameCSV(names,v) {
    return makeObjectFrom2Arrays(names, v)
}

// Překlopí data ze stringu do Date() objektu
// 0000 = neznámé datum, vybere nové datum dle nabídky customDates
// 2030 = stále žije
const parseDates = (date, customDates, name) => {
    let newDate = ''
    switch (date) {
        case "0000":
            newDate = customDates[trimString(name)]
            break;
        case "2030":
            newDate = new Date (date);
            break;
        default:
            newDate = new Date (date);
    }
    return checkDateType(newDate)
}

// Zkontroluje, zdali je datum validní
function checkDateType(newDate){
    if ('Invalid Date' === newDate.toString() || !(newDate instanceof Date)) {
        console.log('Variable:', newDate, typeof(newDate))
        throw new Error ('Unexpected input.')
    }
    return newDate
}