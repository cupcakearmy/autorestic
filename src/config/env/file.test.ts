import { describe, expect, test } from 'bun:test'
import { InvalidEnvFileLine } from '../../errors'
import { parseFile } from './file'

describe('env file', () => {
  test('simple', () => {
    expect(parseFile(`test_foo=ok`)).toEqual({ test_foo: 'ok' })
  })

  test('multiple values', () => {
    expect(parseFile(`test_foo=ok\n \n spacing = foo \n`)).toEqual({ test_foo: 'ok', spacing: 'foo' })
  })

  test('invalid: key', () => {
    expect(() => parseFile(`a=123\na f=ok`)).toThrow(new InvalidEnvFileLine('a f=ok'))
  })

  test('invalid: missing =', () => {
    expect(() => parseFile(`a=123\na ok`)).toThrow(new InvalidEnvFileLine('a ok'))
  })
})
