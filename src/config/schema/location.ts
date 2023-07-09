import { z } from 'zod'
import { NonEmptyString, OptionallyArray } from './common'
import { HooksSchema } from './hooks'
import { OptionsSchema } from './options'

export const LocationSchema = z
  .strictObject({
    from: OptionallyArray(NonEmptyString.describe('local path to backup')),
    to: OptionallyArray(NonEmptyString.describe('repository to backup to')),
    copy: z
      .record(
        NonEmptyString.describe('source repository from which to copy from'),
        OptionallyArray(NonEmptyString.describe('destination repository'))
      )
      .optional(),

    // adapter:
    cron: NonEmptyString.describe('execute backups for the given cron job').optional(),
    hooks: HooksSchema.optional(),
    options: OptionsSchema.optional(),
    forget: z
      .union([
        z.boolean().describe('automatically run "forget" when backing up'),
        z.literal('prune').describe('also prune when forgetting'),
      ])
      .optional(),
  })
  .describe('Location')
