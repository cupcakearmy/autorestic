import { beforeEach, describe, expect, test } from 'bun:test'
import { mkdir, rm } from 'node:fs/promises'
import { lockRepo, unlockRepo, waitForRepo } from '.'
import { Context } from '../models/context'

const mockPath = './test/'
const mockContext: Context = { config: { meta: { path: mockPath } } } as any

describe('lock', () => {
  beforeEach(async () => {
    // Cleanup lock file
    await rm(mockPath, { recursive: true, force: true })
    await mkdir(mockPath, { recursive: true })
  })

  test('simple lock and unlock', () => {
    const repo = 'foo'
    lockRepo(mockContext, repo)
    unlockRepo(mockContext, repo)
  })

  test('should not be able to lock twice', () => {
    const repo = 'foo'
    lockRepo(mockContext, repo)
    expect(() => {
      lockRepo(mockContext, repo)
    }).toThrow()
    unlockRepo(mockContext, repo)
    lockRepo(mockContext, repo)
  })

  test('should be able to eventually acquire lock', async () => {
    const repo = 'foo'
    lockRepo(mockContext, repo)
    setTimeout(() => unlockRepo(mockContext, repo), 50)
    await waitForRepo(mockContext, repo, 1)
  })

  test('unlock', () => {
    unlockRepo(mockContext, 'foo')
  })

  test('multiple', () => {
    const a = 'foo'
    const b = 'bar'
    lockRepo(mockContext, a)
    lockRepo(mockContext, b)
    unlockRepo(mockContext, b)
    unlockRepo(mockContext, a)
  })
})
