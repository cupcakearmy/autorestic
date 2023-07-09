import { readFileSync, writeFileSync } from 'node:fs'
import yaml from 'yaml'
import { relativePath } from '../config/resolution'
import { LockfileAlreadyLocked } from '../errors'
import { Log } from '../logger'
import { Context } from '../models/context'
import { Lockfile, LockfileSchema } from './schema'
import { wait } from '../utils/time'

const LOCKFILE = '.autorestic.lock'
const VERSION = 2
const l = Log.child({ command: 'lock' })

function load(ctx: Context): Lockfile {
  const defaultLockfile = { version: VERSION, cron: {}, running: {} }
  try {
    const path = relativePath(ctx.config.meta.path, LOCKFILE)
    l.trace('looking for lock file', { path })
    // throw new Error(path)
    const rawConfig = readFileSync(path, 'utf-8')
    const config = yaml.parse(rawConfig)
    const parsed = LockfileSchema.safeParse(config)
    if (!parsed.success) return defaultLockfile
    if (parsed.data.version < VERSION) {
      l.debug('lockfile is old and will be overwritten')
      return defaultLockfile
    }
    return parsed.data
  } catch {
    return defaultLockfile
  }
}

function write(ctx: Context, lockfile: Lockfile) {
  const path = relativePath(ctx.config.meta.path, LOCKFILE)
  writeFileSync(path, yaml.stringify(lockfile), 'utf-8')
}

export function lockRepo(ctx: Context, repo: string) {
  const lock = load(ctx)
  l.trace('trying to lock repository', { repo })
  if (lock.running[repo]) throw new LockfileAlreadyLocked(repo)
  lock.running[repo] = true
  write(ctx, lock)
}

/**
 * Waits for a repo to become unlocked, and errors if it does not succeed in the given timeout.
 *
 * @param [timeout=10] max seconds to wait for repo to become unlocked
 */
export async function waitForRepo(ctx: Context, repo: string, timeout = 10) {
  const now = Date.now()
  while (Date.now() - now < timeout * 1_000) {
    try {
      lockRepo(ctx, repo)
      l.trace('repo is free again', { repo })
      break
    } catch {
      l.trace('waiting for repo to be unlocked', { repo })
      await wait(0.1) // Wait for 100ms
    }
  }
  throw new LockfileAlreadyLocked(repo)
}

export function updateLastRun(ctx: Context, location: string) {
  const lock = load(ctx)
  lock.cron[location] = Date.now()
  write(ctx, lock)
}

export function unlockRepo(ctx: Context, repo: string) {
  l.trace('unlocking repository', { repo })
  const lock = load(ctx)
  lock.running[repo] = false
  write(ctx, lock)
}
