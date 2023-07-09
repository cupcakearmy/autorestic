import { z } from 'zod'

export const LockfileSchema = z.strictObject({
  version: z.number().min(0).describe('lockfile version'),
  running: z
    .record(z.string().describe('repository'), z.boolean().describe('whether repository is running'))
    .describe('running information for each repository'),
  cron: z
    .record(z.string().describe('location'), z.number().describe('timestamp of last backup'))
    .describe('information about last run for a given location. in milliseconds'),
})

export type Lockfile = z.infer<typeof LockfileSchema>
