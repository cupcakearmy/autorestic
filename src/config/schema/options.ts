import { z } from 'zod'
import { NonEmptyString, OptionallyArray } from './common'

const OptionSchema = z.record(
  NonEmptyString.describe('native restic option'),
  z.union([z.literal(true).describe('boolean flag'), OptionallyArray(NonEmptyString)]).describe('value of option')
)

export const OptionsSchema = z
  .strictObject({
    all: OptionSchema.optional(),
    backup: OptionSchema.optional(),
    forget: OptionSchema.optional(),
  })
  .describe('native restic options')
