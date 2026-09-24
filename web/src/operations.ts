export const operations = [
  { id: 'add', symbol: '+', label: 'Add', needsB: true },
  { id: 'subtract', symbol: '−', label: 'Subtract', needsB: true },
  { id: 'multiply', symbol: '×', label: 'Multiply', needsB: true },
  { id: 'divide', symbol: '÷', label: 'Divide', needsB: true },
  { id: 'modulus', symbol: 'mod', label: 'Modulus', needsB: true },
  { id: 'power', symbol: 'xʸ', label: 'Power', needsB: true },
  { id: 'sqrt', symbol: '√', label: 'Square root', needsB: false },
  { id: 'percentage', symbol: '%', label: 'Percentage', needsB: true },
] as const

export type Operation = (typeof operations)[number]['id']
