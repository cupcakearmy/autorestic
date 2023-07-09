import { exists, readFile } from 'node:fs/promises'
import yaml from 'yaml'
import { ConfigFileNotFound, CustomError, InvalidConfigFile } from '../errors'
import { enrichConfig } from './env/file'
import { autoLocateConfig } from './resolution'
import { Config, ConfigWithMetaSchema } from './schema/config'
import { basename } from 'node:path'

export async function loadConfig(customPath?: string): Promise<Config> {
  let path: string
  if (customPath) {
    path = customPath
    if (!(await exists(path))) throw new ConfigFileNotFound([path])
  } else {
    path = await autoLocateConfig()
  }

  const rawConfig = await readFile(path, 'utf-8')
  const config = yaml.parse(rawConfig)
  await enrichConfig(config, path)
  config.meta = { path: basename(path) }
  const parsed = ConfigWithMetaSchema.safeParse(config)
  if (!parsed.success)
    throw new InvalidConfigFile(parsed.error.errors.map((e) => `${e.path.join(' > ')}: ${e.message}`))

  // Check for semantics

  return parsed.data
}
