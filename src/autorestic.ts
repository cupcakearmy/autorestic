import 'colors'
import minimist from 'minimist'
import { resolve } from 'path'

import { init } from './config'
import handlers, { error, help } from './handlers'
import { Config } from './types'


process.on('uncaughtException', err => {
	console.log(err.message)
	process.exit(1)
})

export const { _: commands, ...flags } = minimist(process.argv.slice(2), {
	alias: {
		'c': 'config',
		'v': 'verbose',
		'h': 'help',
		'a': 'all',
		'l': 'location',
		'b': 'backend',
	},
	boolean: ['a'],
	string: ['l', 'b'],
})

export const VERSION = '0.1'
export const DEFAULT_CONFIG = '~/.autorestic.yml'
export const INSTALL_DIR = '/usr/local/bin'
export const CONFIG_FILE: string = resolve(flags.config || DEFAULT_CONFIG)
export const VERBOSE = flags.verbose

export const config: Config = init()

if (commands.length < 1)
	help()
else {
	const command: string = commands[0]
	const args: string[] = commands.slice(1)
	;(handlers[command] || error)(args, flags)
}
