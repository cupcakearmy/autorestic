import CronParser from 'cron-parser'

import { config } from './autorestic'
import { checkAndConfigureBackendsForLocations } from './backend'
import { Location } from './types'
import { backupLocation } from './backup'
import { readLock, writeLock } from './lock'


const runCronForLocation = (name: string, location: Location) => {
  const lock = readLock()
  const parsed = CronParser.parseExpression(location.cron || '')
  const last = parsed.prev()

  if (!lock.crons[name] || last.toDate().getTime() > lock.crons[name].lastRun) {
    backupLocation(name, location)
    lock.crons[name] = { lastRun: Date.now() }
    writeLock(lock)
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