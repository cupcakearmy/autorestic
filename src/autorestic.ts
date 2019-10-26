import 'colors'
import minimist from 'minimist'

import { init } from './config'
import handlers, { error, help } from './handlers'
import { Config } from './types'

process.on('uncaughtException', err => {
  console.log(err.message)
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
  },
  boolean: ['a'],
  string: ['l', 'b'],
})

export const VERSION = '0.5'
export const INSTALL_DIR = '/usr/local/bin'
export const VERBOSE = flags.verbose

if (flags.version) {
  console.log('version'.grey, VERSION)
  process.exit(0)
}

export const config = init()

function main() {
  if (commands.length < 1) return help()

  const command: string = commands[0]
  const args: string[] = commands.slice(1)
  ;(handlers[command] || error)(args, flags)
}

main()
