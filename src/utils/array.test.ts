import { describe, expect, test } from 'bun:test'
import { isSubset } from './array'

describe('set theory', () => {
  test('subset', () => {
    expect(isSubset([1], [1, 2])).toBe(true)
    expect(isSubset([1], [2])).toBe(false)
    expect(isSubset([], [])).toBe(true)
    expect(isSubset([1, 2, 3], [1, 2])).toBe(false)
    expect(isSubset([1, 2], [1, 2])).toBe(true)
  })
})
