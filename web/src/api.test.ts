import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'

afterEach(() => {
  vi.unstubAllGlobals()
})

function mockFetch(status: number, body: unknown) {
  const fetchMock = vi.fn().mockImplementation(() =>
    Promise.resolve(
      new Response(JSON.stringify(body), {
        status,
        headers: { 'Content-Type': 'application/json' },
      }),
    ),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('calculate', () => {
  it('posts a binary operation with both numbers and returns zero', async () => {
    const fetchMock = mockFetch(200, { result: 0 })

    await expect(calculate('subtract', 4, 4)).resolves.toBe(0)
    expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation: 'subtract', a: 4, b: 4 }),
    })
  })

  it('omits the second number for square root but includes a zero second number', async () => {
    const fetchMock = mockFetch(200, { result: 3 })

    await calculate('sqrt', 9)
    await calculate('multiply', 9, 0)

    expect(JSON.parse(fetchMock.mock.calls[0][1].body)).toEqual({
      operation: 'sqrt',
      a: 9,
    })
    expect(JSON.parse(fetchMock.mock.calls[1][1].body)).toEqual({
      operation: 'multiply',
      a: 9,
      b: 0,
    })
  })

  it('uses an API error message when the request fails', async () => {
    mockFetch(400, { error: 'division by zero' })

    await expect(calculate('divide', 4, 0)).rejects.toThrow('division by zero')
  })

  it('uses a fallback message when the API error has no message', async () => {
    mockFetch(500, {})

    await expect(calculate('add', 1, 2)).rejects.toThrow(
      'Calculation failed. Please try again.',
    )
  })

  it('rejects a successful response without a numeric result', async () => {
    mockFetch(200, { result: '3' })

    await expect(calculate('add', 1, 2)).rejects.toThrow(
      'The API returned an invalid result.',
    )
  })
})
