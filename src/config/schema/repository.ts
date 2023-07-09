import { z } from 'zod'
import { NonEmptyString } from './common'
import { OptionsSchema } from './options'

export const RepositorySchema = z.strictObject({
  type: z.enum(['local', 'sftp', 'rest', 'swift', 's3', 'b2', 'azure', 'gs', 'rclone']).describe('type of repository'),
  path: NonEmptyString.describe('restic path'),
  key: NonEmptyString.describe('encryption key for the repository'),
  env: z
    .record(
      NonEmptyString.describe('environment variable'),
      NonEmptyString.describe('value of the environment variable')
    )
    .transform((env) => Object.fromEntries(Object.entries(env).map(([key, value]) => [key.toUpperCase(), value])))
    .describe('environment variables')
    .optional(),
  options: OptionsSchema.describe('options').optional(),
})
