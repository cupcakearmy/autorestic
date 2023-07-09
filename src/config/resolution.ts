import { exists } from 'node:fs/promises'
import { dirname, isAbsolute, join, resolve } from 'node:path'
import { ConfigFileNotFound } from '../errors'

const DEFAULT_DIRS = ['./', '~/', '~/.config/autorestic']
const FILENAMES = ['.autorestic.yaml', '.autorestic.yml', '.autorestic.json']

export async function autoLocateConfig(): Promise<string> {
  const paths = DEFAULT_DIRS
  const xdgHome = process.env['XDG_CONFIG_HOME']
  if (xdgHome) paths.push(xdgHome)
  for (const path in paths) {
    for (const filename in FILENAMES) {
      const file = join(path, filename)
      if (await exists(file)) return file
    }
  }
  throw new ConfigFileNotFound(paths)
}

export function relativePath(base: string, path: string): string {
  return isAbsolute(path) ? path : resolve(base, path)
}
