import { act, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import { calculate } from './api'

vi.mock('./api', () => ({ calculate: vi.fn() }))

const mockedCalculate = vi.mocked(calculate)

beforeEach(() => {
  mockedCalculate.mockReset()
})

describe('calculator form', () => {
  it('submits the chosen binary operation and clears the answer after an edit', async () => {
    const user = userEvent.setup()
    mockedCalculate.mockResolvedValue(4)
    render(<App />)

    await user.click(screen.getByRole('button', { name: 'Divide' }))
    expect(screen.getByRole('button', { name: 'Divide' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '8',
    )
    await user.type(
      screen.getByRole('spinbutton', { name: 'Second number' }),
      '2',
    )
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(mockedCalculate).toHaveBeenCalledWith('divide', 8, 2)
    expect(await screen.findByText('4')).toBeInTheDocument()

    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '0',
    )
    expect(screen.getByText('—')).toBeInTheDocument()
  })

  it('uses one input for square root and omits the second argument', async () => {
    const user = userEvent.setup()
    mockedCalculate.mockResolvedValue(3)
    render(<App />)

    await user.click(screen.getByRole('button', { name: 'Square root' }))
    expect(
      screen.queryByRole('spinbutton', { name: 'Second number' }),
    ).not.toBeInTheDocument()
    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '9',
    )
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(mockedCalculate).toHaveBeenCalledWith('sqrt', 9, undefined)
    expect(await screen.findByText('3')).toBeInTheDocument()
  })

  it('shows the percentage explanation for that operation', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(screen.getByRole('button', { name: 'Percentage' }))

    expect(
      screen.getByText(
        'Calculates the first number as a percentage of the second.',
      ),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Percentage' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
  })

  it('shows the exponent operator for power', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(screen.getByRole('button', { name: 'Power' }))

    expect(screen.getByText('^')).toBeInTheDocument()
  })

  it('does not submit when a required number is missing', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '8',
    )
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(mockedCalculate).not.toHaveBeenCalled()
  })

  it('shows its own validation error when the form is submitted without a number', async () => {
    const user = userEvent.setup()
    render(<App />)
    const form = screen
      .getByRole('button', { name: 'Calculate' })
      .closest('form')!
    form.noValidate = true

    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(
      screen.getByText('Enter a valid number in each required field.'),
    ).toBeInTheDocument()
    expect(mockedCalculate).not.toHaveBeenCalled()
  })

  it('disables controls while calculating and restores them when done', async () => {
    const user = userEvent.setup()
    let finishCalculation: (value: number) => void = () => {}
    mockedCalculate.mockImplementation(
      () =>
        new Promise<number>((resolve) => {
          finishCalculation = resolve
        }),
    )
    render(<App />)

    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '1',
    )
    await user.type(
      screen.getByRole('spinbutton', { name: 'Second number' }),
      '2',
    )
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(screen.getByRole('button', { name: 'Calculating…' })).toBeDisabled()
    expect(
      screen.getByRole('spinbutton', { name: 'First number' }),
    ).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Add' })).toBeDisabled()

    await act(async () => finishCalculation(3))

    expect(screen.getByText('3')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  })

  it.each([
    [new Error('division by zero'), 'division by zero'],
    [
      new TypeError('Failed to fetch'),
      'Could not reach the API. Check that the Go server is running.',
    ],
    ['unexpected failure', 'Calculation failed. Please try again.'],
  ])('displays the appropriate error for %s', async (failure, message) => {
    const user = userEvent.setup()
    mockedCalculate.mockRejectedValue(failure)
    render(<App />)

    await user.type(
      screen.getByRole('spinbutton', { name: 'First number' }),
      '8',
    )
    await user.type(
      screen.getByRole('spinbutton', { name: 'Second number' }),
      '0',
    )
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText(message)).toBeInTheDocument()
    expect(screen.queryByText('—')).not.toBeInTheDocument()
  })
})
