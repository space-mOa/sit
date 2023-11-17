// ! KÓD BYL PŘEVZAT A UPRAVEN: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

import { makeObjectFrom2Arrays, trimString } from './helpers.js';

export { SociologistCustomDates }

const SociologistsNames = [
    'Dědek František',
    'Slaminka Vladimír',
    'Marušiak Martin'
]

const SociologistsDates = [
    new Date (1970),
    new Date (1970),
    new Date (1970)
]

const SociologistCustomDates = makeObjectFrom2Arrays(SociologistsNames.map(s => trimString(s)), SociologistsDates)


