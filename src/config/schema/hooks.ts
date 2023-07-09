import { z } from 'zod'
import { NonEmptyString, OptionallyArray } from './common'

const Command = NonEmptyString.describe('command to be executed')
const Commands = OptionallyArray(Command).describe('list of commands')

export const HooksSchema = z
  .strictObject({
    before: Commands.optional(),
    after: Commands.optional(),
    failure: Commands.optional(),
    success: Commands.optional(),
  })
  .describe('hooks to be executed')
