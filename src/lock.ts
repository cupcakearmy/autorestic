import fs from 'fs'

import { pathRelativeToConfigFile } from './utils'
import { Lockfile } from './types'

export const getLockFileName = () => {
  const LOCK_FILE = '.autorestic.lock'
  return pathRelativeToConfigFile(LOCK_FILE)
}

export const readLock = (): Lockfile => {
  const name = getLockFileName()
  let lock = {
    running: false,
    crons: {},
  }
  try {
    lock = JSON.parse(fs.readFileSync(name, { encoding: 'utf-8' }))
  } catch {}
  return lock
}
export const writeLock = (lock: Lockfile) => {
  const name = getLockFileName()
  fs.writeFileSync(name, JSON.stringify(lock, null, 2), { encoding: 'utf-8' })
}

export const unlock = () => {
  writeLock({
    ...readLock(),
    running: false,
  })
}
