Manual Setup

To manually set up Vitest, install vitest and the following packages as dev dependencies:
Terminal

# Using TypeScript
npm install -D vitest @vitejs/plugin-react jsdom @testing-library/react @testing-library/dom vite-tsconfig-paths
# Using JavaScript
npm install -D vitest @vitejs/plugin-react jsdom @testing-library/react @testing-library/dom

Create a vitest.config.mts|js file in the root of your project, and add the following options:
vitest.config.mts
TypeScript

import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'
 
export default defineConfig({
  plugins: [tsconfigPaths(), react()],
  test: {
    environment: 'jsdom',
  },
})

For more information on configuring Vitest, please refer to the Vitest Configuration

docs.

Then, add a test script to your package.json:
package.json

{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "test": "vitest"
  }
}

When you run npm run test, Vitest will watch for changes in your project by default.
Creating your first Vitest Unit Test

Check that everything is working by creating a test to check if the <Page /> component successfully renders a heading:
app/page.tsx
TypeScript

import Link from 'next/link'
 
export default function Page() {
  return (
    <div>
      <h1>Home</h1>
      <Link href="/about">About</Link>
    </div>
  )
}

__tests__/page.test.tsx
TypeScript

import { expect, test } from 'vitest'
import { render, screen } from '@testing-library/react'
import Page from '../app/page'
 
test('Page', () => {
  render(<Page />)
  expect(screen.getByRole('heading', { level: 1, name: 'Home' })).toBeDefined()
})

    Good to know: The example above uses the common __tests__ convention, but test files can also be colocated inside the app router.

Running your tests

Then, run the following command to run your tests:
Terminal

npm run test
# or
yarn test
# or
pnpm test
# or
bun test
