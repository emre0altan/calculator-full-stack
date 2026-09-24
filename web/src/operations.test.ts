import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { cwd } from 'node:process'
import { describe, expect, it } from 'vitest'
import { operations } from './operations'

// Read the real API contract so a backend change cannot silently drift from the UI.
const schema = JSON.parse(
  readFileSync(
    resolve(cwd(), '../backend/internal/httpapi/calculation.schema.json'),
    'utf8',
  ),
)

describe('operation contract', () => {
  it('offers exactly the operations accepted by the API', () => {
    expect(operations.map(({ id }) => id).sort()).toEqual(
      [...schema.properties.operation.enum].sort(),
    )
  })

  it('requires the second operand for the same operations as the API', () => {
    expect(schema.required).toContain('a')
    expect(schema.required).not.toContain('b')
    expect(schema.then.required).toEqual(['b'])
    expect(
      operations.filter(({ needsB }) => !needsB).map(({ id }) => id),
    ).toEqual([schema.if.properties.operation.not.const])
  })
})
