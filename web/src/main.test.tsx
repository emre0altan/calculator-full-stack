import { act, screen } from '@testing-library/react'
import { expect, it } from 'vitest'

it('mounts the calculator in the page root', async () => {
  document.body.innerHTML = '<div id="root"></div>'

  await act(async () => {
    await import('./main')
  })

  expect(
    screen.getByRole('heading', { name: 'Calculator.' }),
  ).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeInTheDocument()
})
