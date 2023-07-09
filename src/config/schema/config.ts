import { z } from 'zod'
import { asArray } from '../../utils/array'
import { RepositorySchema } from './repository'
import { NonEmptyString } from './common'
import { LocationSchema } from './location'
import { OptionsSchema } from './options'

export const ConfigSchema = z.strictObject({
  version: z.number().describe('version number'),
  repos: z.record(NonEmptyString.describe('repository name'), RepositorySchema).describe('available repositories'),
  locations: z.record(NonEmptyString.describe('location name'), LocationSchema).describe('available locations'),
  global: z
    .strictObject({
      options: OptionsSchema.optional(),
    })
    .describe('global configuration')
    .optional(),
  extras: z.any().optional(),
})

const ConfigMeta = z
  .strictObject({
    path: NonEmptyString.describe('The path of the loaded config'),
  })
  .describe('Meta information about the config')

export const ConfigWithMetaSchema = ConfigSchema.extend({
  meta: ConfigMeta,
}).superRefine((config, ctx) => {
  const availableRepos = Object.keys(config.repos)
  for (const [name, location] of Object.entries(config.locations)) {
    const locationPath = [...ctx.path, 'locations', name]
    const toRepos = asArray(location.to)
    // Check if all target repos are valid
    for (const to of toRepos) {
      if (!availableRepos.includes(to)) {
        const message = `location "${name}" has an invalid repository "${to}"`
        ctx.addIssue({ message, code: 'custom', path: [...locationPath, 'to'] })
      }
    }
    // Check copy field
    if (!location.copy) continue
    for (const [source, destinations] of Object.entries(location.copy)) {
      const path = [...locationPath, 'copy', source]
      if (!toRepos.includes(source))
        ctx.addIssue({
          code: 'custom',
          path,
          message: `copy source "${source}" must be also a backup target`,
        })
      for (const destination of asArray(destinations)) {
        if (destination === source)
          ctx.addIssue({
            code: 'custom',
            path: [...path, destination],
            message: `destination repository "${destination}" cannot be also the source in copy field`,
          })
        if (!availableRepos.includes(destination))
          ctx.addIssue({
            code: 'custom',
            path: [...path, destination],
            message: `destination repository "${destination}" does not exist`,
          })
      }
    }
  }
})

export type Config = z.infer<typeof ConfigWithMetaSchema>
