import fs from 'fs'

import CronParser from 'cron-parser'

import { config } from './autorestic'
import { checkAndConfigureBackendsForLocations } from './backend'
import { Location, Lockfile } from './types'
import { backupLocation } from './backup'
import { pathRelativeToConfigFile } from './utils'


const getLockFileName = () => {
  const LOCK_FILE = '.autorestic.lock'
  return pathRelativeToConfigFile(LOCK_FILE)
}

const readLock = (): Lockfile => {
  const name = getLockFileName()
  let lock = {}
  try {
    lock = JSON.parse(fs.readFileSync(name, { encoding: 'utf-8' }))
  } catch { }
  return lock
}
const writeLock = (diff: Lockfile = {}) => {
  const name = getLockFileName()
  const newLock = Object.assign(
    readLock(),
    diff
  )
  fs.writeFileSync(name, JSON.stringify(newLock, null, 2), { encoding: 'utf-8' })
}

const runCronForLocation = (name: string, location: Location) => {
  const lock = readLock()[name]
  const parsed = CronParser.parseExpression(location.cron || '')
  const last = parsed.prev()

  if (!lock || last.toDate().getTime() > lock.lastRun) {
    backupLocation(name, location)
    writeLock({
      [name]: {
        ...lock,
        lastRun: Date.now()
      }
    })
  } else {
    console.log(`${name.yellow} ▶ Skipping. Sheduled for: ${parsed.next().toString().underline.blue}`)
  }
}

export const runCron = () => {
  const locationsWithCron = Object.entries(config.locations).filter(([name, { cron }]) => !!cron)
  checkAndConfigureBackendsForLocations(Object.fromEntries(locationsWithCron))

  console.log('\nRunning cron jobs'.underline.gray)
  for (const [name, location] of locationsWithCron)
    runCronForLocation(name, location)

  console.log('\nFinished!'.underline + ' 🎉')
}