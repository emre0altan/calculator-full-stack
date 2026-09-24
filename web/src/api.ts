import type { Operation } from './operations'

export async function calculate(
  operation: Operation,
  a: number,
  b?: number,
): Promise<number> {
  const response = await fetch('/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(
      b === undefined ? { operation, a } : { operation, a, b },
    ),
  })

  const data: { result?: number; error?: string } = await response.json()
  if (!response.ok) {
    throw new Error(data.error || 'Calculation failed. Please try again.')
  }
  if (typeof data.result !== 'number') {
    throw new Error('The API returned an invalid result.')
  }
  return data.result
}
