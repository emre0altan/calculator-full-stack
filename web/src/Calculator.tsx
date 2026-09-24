import { useState, type FormEvent } from 'react'
import { calculate } from './api'
import { operations, type Operation } from './operations'

export default function Calculator() {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [operation, setOperation] = useState<Operation>('add')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const selectedOperation = operations.find((item) => item.id === operation)!
  const isUnary = !selectedOperation.needsB
  const inputOperator =
    operation === 'percentage'
      ? '% of'
      : operation === 'power'
        ? '^'
        : selectedOperation.symbol

  function resetFeedback() {
    setResult(null)
    setError(null)
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const first = Number(a)
    const second = Number(b)

    if (
      !a.trim() ||
      (!isUnary && !b.trim()) ||
      !Number.isFinite(first) ||
      (!isUnary && !Number.isFinite(second))
    ) {
      setResult(null)
      setError('Enter a valid number in each required field.')
      return
    }

    setLoading(true)
    resetFeedback()
    try {
      setResult(await calculate(operation, first, isUnary ? undefined : second))
    } catch (caught) {
      setError(
        caught instanceof TypeError
          ? 'Could not reach the API. Check that the Go server is running.'
          : caught instanceof Error
            ? caught.message
            : 'Calculation failed. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="calculator" onSubmit={handleSubmit}>
      <div className="section-heading">
        <span className="step">01</span>
        <h2>Your numbers</h2>
      </div>
      <div className={`input-grid ${isUnary ? 'single-input' : ''}`}>
        <label className="field">
          <span>First number</span>
          <input
            type="number"
            step="any"
            required
            value={a}
            onChange={(event) => {
              setA(event.target.value)
              resetFeedback()
            }}
            placeholder="e.g. 12"
            disabled={loading}
          />
        </label>
        {!isUnary && (
          <>
            <span className="input-operator">{inputOperator}</span>
            <label className="field">
              <span>Second number</span>
              <input
                type="number"
                step="any"
                required
                value={b}
                onChange={(event) => {
                  setB(event.target.value)
                  resetFeedback()
                }}
                placeholder="e.g. 4"
                disabled={loading}
              />
            </label>
          </>
        )}
      </div>

      <div className="section-heading operation-heading">
        <span className="step">02</span>
        <h2>Choose an operation</h2>
      </div>
      <div className="operations" role="group" aria-label="Operation">
        {operations.map((item) => (
          <button
            key={item.id}
            className={`operation ${operation === item.id ? 'selected' : ''}`}
            type="button"
            aria-pressed={operation === item.id}
            disabled={loading}
            onClick={() => {
              setOperation(item.id)
              resetFeedback()
            }}
          >
            <span className="symbol" aria-hidden="true">
              {item.symbol}
            </span>
            <span className="operation-label">{item.label}</span>
          </button>
        ))}
      </div>
      {operation === 'percentage' && (
        <p className="hint">
          Calculates the first number as a percentage of the second.
        </p>
      )}

      <button className="calculate-button" type="submit" disabled={loading}>
        {loading ? 'Calculating…' : 'Calculate'}
        <span aria-hidden="true">→</span>
      </button>

      <div
        className={`answer ${error ? 'answer-error' : ''}`}
        aria-live="polite"
        aria-atomic="true"
      >
        <span className="answer-label">RESULT</span>
        {error ? (
          <p className="error-message">{error}</p>
        ) : (
          <output>{result === null ? '—' : String(result)}</output>
        )}
      </div>
    </form>
  )
}
