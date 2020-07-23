import 'colors'
import minimist from 'minimist'

import { init } from './config'
import handlers, { error, help } from './handlers'
import { Config } from './types'
import { readLock, writeLock, unlock } from './lock'

process.on('uncaughtException', (err) => {
  console.log(err.message)
  unlock()
  process.exit(1)
})

export const { _: commands, ...flags } = minimist(process.argv.slice(2), {
  alias: {
    c: 'config',
    v: 'version',
    h: 'help',
    a: 'all',
    l: 'location',
    b: 'backend',
    d: 'dry-run',
  },
  boolean: ['a', 'd'],
  string: ['l', 'b'],
})

export const VERSION = '0.20'
export const INSTALL_DIR = '/usr/local/bin'
export const VERBOSE = flags.verbose

export let config: Config

async function main() {
  config = init()

  // Don't let 2 instances run on the same config
  const lock = readLock()
  if (lock.running) {
    console.log('An instance of autorestic is already running for this config file'.red)
    return
  }
  writeLock({
    ...lock,
    running: true,
  })

  // For dev
  // return await handlers['cron']([], { ...flags, all: true })

  if (commands.length < 1 || commands[0] === 'help') return help()

  const command: string = commands[0]
  const args: string[] = commands.slice(1)

  const fn = handlers[command] || error
  await fn(args, flags)
}

main()
  .catch((e: Error) => console.error(e.message))
  .finally(unlock)
