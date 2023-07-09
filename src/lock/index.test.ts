import { describe, expect, mock, test, beforeEach } from 'bun:test'
import { lockRepo } from '.'
import { Context } from '../models/context'
import { mkdir, rm } from 'node:fs/promises'

const mockPath = './test/'
const mockContext: Context = { config: { meta: { path: mockPath } } } as any

describe('lock', () => {
  beforeEach(async () => {
    // Cleanup lock file
    await rm(mockPath, { recursive: true, force: true })
    await mkdir(mockPath, { recursive: true })
  })

  test('lock', () => {
    lockRepo(mockContext, 'foo')
    // lockRepo(mockContext, 'foo')
  })
})
