import Calculator from './Calculator'

export default function App() {
  return (
    <main className="page">
      <div className="layout">
        <header className="intro">
          <span className="eyebrow">A LITTLE MATH, MADE EASY</span>
          <h1>
            Calculator<span className="accent">.</span>
          </h1>
          <p>Pick an operation, enter your numbers, and get the answer.</p>
        </header>

        <Calculator />
      </div>
    </main>
  )
}
