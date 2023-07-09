import { exists, readFile } from 'node:fs/promises'
import { InvalidEnvFileLine } from '../../errors'
import { setByPath } from '../../utils/path'
import { relativePath } from '../resolution'

export function parseFile(contents: string) {
  const variables: Record<string, string> = {}
  const lines = contents
    .trim()
    .split('\n')
    .map((l) => l.trim())
  const matcher = /^\s*(?<variable>\w+)\s*=(?<value>.*)$/
  for (const line of lines) {
    if (!line) continue
    const match = matcher.exec(line)
    if (!match) throw new InvalidEnvFileLine(line)
    variables[match.groups!.variable] = match.groups!.value.trim()
  }
  return variables
}

const PREFIX = 'AUTORESTIC_'

function envVariableToObjectPath(env: string): string {
  if (env.startsWith(PREFIX)) env = env.replace(PREFIX, '')
  return (
    env
      // Convert to object path
      .replaceAll('_', '.')
      // Escape the double unterscore. __ -> .. -> _
      .replaceAll('..', '_')
      .toLowerCase()
  )
}

/**
 * Fill the config file with the env file variables.
 * These take precedence before the config file itself.
 */
export async function enrichConfig(rawConfig: any, path: string) {
  const envFilePath = relativePath(path, '.autorestic.env')
  let variables: Record<string, string> = {}

  if (await exists(envFilePath)) {
    const envFile = parseFile(await readFile(envFilePath, 'utf-8'))
    Object.assign(variables, envFile)
  }

  Object.assign(variables, process.env)

  for (const [key, value] of Object.entries(variables)) {
    if (!key.startsWith(PREFIX)) continue
    setByPath(rawConfig, envVariableToObjectPath(key), value)
  }
}
