import { ZodTypeAny, z } from 'zod'
import { Log } from '../../logger'

export const NonEmptyString = z
  .string()
  .min(1)
  // Extrapolate env variables from a string
  .transform((s) => {
    return s.replaceAll(/\$(\w+)|\${(\w+)}/g, (_, g0, g1) => {
      const variable = g0 || g1
      const value = process.env[variable] ?? ''
      if (!value) Log.error(`cannot find environment variable "${variable}" to replace in ${s}`)
      return value
    })
  })
  .describe('non-empty string that can extrapolate env variables inside it')

export function OptionallyArray<T extends ZodTypeAny>(type: T) {
  return z.union([type, z.array(type).min(1)])
}
