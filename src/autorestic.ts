import 'colors'
import minimist from 'minimist'

import { init } from './config'
import handlers, { error, help } from './handlers'



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
		d: 'dry-run',
	},
	boolean: ['a', 'd'],
	string: ['l', 'b'],
})

export const VERSION = '0.16'
export const INSTALL_DIR = '/usr/local/bin'
export const VERBOSE = flags.verbose

export const config = init()


async function main() {
	if (commands.length < 1 || commands[0] === 'help') return help()

	const command: string = commands[0]
	const args: string[] = commands.slice(1)

	const fn = handlers[command] || error
	await fn(args, flags)
}


main().catch((e: Error) => console.error(e.message))
