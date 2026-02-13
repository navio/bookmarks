import { describe, expect, it } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

function read(p: string) {
  return fs.readFileSync(p, 'utf8')
}

describe('docs site', () => {
  it('has the expected pages', () => {
    const root = path.resolve(__dirname, '..')
    const pages = ['index.md', 'getting-started.md', 'commands.md', 'workflows.md', 'store.md']

    for (const page of pages) {
      expect(fs.existsSync(path.join(root, page))).toBe(true)
    }
  })

  it('documents key commands', () => {
    const root = path.resolve(__dirname, '..')
    const commands = read(path.join(root, 'commands.md'))
    expect(commands).toContain('bm add')
    expect(commands).toContain('bm ls')
    expect(commands).toContain('bm find')
    expect(commands).toContain('bm table')
    expect(commands).toContain('bm path')
    expect(commands).toContain('bm update')
    expect(commands).toContain('bm rm')
  })
})
