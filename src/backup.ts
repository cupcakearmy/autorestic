import { Writer } from 'clitastic'

import { config, VERBOSE } from './autorestic'
import { getEnvFromBackend } from './backend'
import { Locations, Location } from './types'
import {
	exec,
	ConfigError,
	pathRelativeToConfigFile,
	getFlagsFromLocation,
	makeArrayIfIsNot,
	execPlain,
	MeasureDuration, fill,
} from './utils'



export const backupSingle = (name: string, to: string, location: Location) => {
	if (!config) throw ConfigError
	const delta = new MeasureDuration()
	const writer = new Writer(name + to.blue + ' : ' + 'Backing up... ⏳')

	const backend = config.backends[to]
	const path = pathRelativeToConfigFile(location.from)

	const cmd = exec(
		'restic',
		['backup', path, ...getFlagsFromLocation(location, 'backup')],
		{ env: getEnvFromBackend(backend) },
	)

	if (VERBOSE) console.log(cmd.out, cmd.err)
	writer.done(`${name}${to.blue} : ${'Done ✓'.green} (${delta.finished(true)})`)
}

export const backupLocation = (name: string, location: Location) => {
	const display = name.yellow + ' ▶ '
	const filler = fill(name.length + 3)
	let first = true

	if (location.hooks && location.hooks.before)
		for (const command of makeArrayIfIsNot(location.hooks.before)) {
			const cmd = execPlain(command)
			if (cmd) console.log(cmd.out, cmd.err)
		}

	for (const t of makeArrayIfIsNot(location.to)) {
		backupSingle(first ? display : filler, t, location)
		if (first) first = false
	}

	if (location.hooks && location.hooks.after)
		for (const command of makeArrayIfIsNot(location.hooks.after)) {
			const cmd = execPlain(command)
			if (cmd) console.log(cmd.out, cmd.err)
		}
}

export const backupAll = (locations?: Locations) => {
	if (!locations) {
		if (!config) throw ConfigError
		locations = config.locations
	}

	console.log('\nBacking Up'.underline.grey)
	for (const [name, location] of Object.entries(locations))
		backupLocation(name, location)
}
