import { expect, test, describe } from 'bun:test'
import { setByPath } from './path'

describe('set by path', () => {
  test('simple', () => {
    expect(setByPath({}, 'a', true)).toEqual({ a: true })
    expect(setByPath({}, 'f', { ok: true })).toEqual({ f: { ok: true } })
    expect(setByPath([], '0', true)).toEqual([true])
    expect(setByPath([], '2', false)).toEqual([undefined, undefined, false])
  })

  test('object', () => {
    expect(setByPath({}, 'a.b', true)).toEqual({ a: { b: true } })
    expect(setByPath({}, 'a.b.c', true)).toEqual({ a: { b: { c: true } } })
    expect(setByPath({ a: true }, 'b', false)).toEqual({ a: true, b: false })
    expect(setByPath({ a: { b: true } }, 'a.c', false)).toEqual({ a: { b: true, c: false } })

    expect(() => setByPath({ a: 'foo' }, 'a.b', true)).toThrow()
    expect(setByPath({ a: 'foo' }, 'a', true)).toEqual({ a: true })
  })

  test('array', () => {
    expect(() => setByPath([], 'a', true)).toThrow()
    expect(setByPath([], '0', true)).toEqual([true])
    expect(setByPath([], '0.0.0', true)).toEqual([[[true]]])
    expect(setByPath([], '0.1.2', true)).toEqual([[undefined, [undefined, undefined, true]]])
  })

  test('mixed', () => {
    expect(setByPath({ items: [] }, 'items.0.name', 'John')).toEqual({ items: [{ name: 'John' }] })
    expect(setByPath([], '0.name', 'John')).toEqual([{ name: 'John' }])
  })
})
