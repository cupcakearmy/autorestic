import { Writer } from 'clitastic'

import { config, VERBOSE } from './autorestic'
import { getEnvFromBackend } from './backend'
import { LocationFromPrefixes } from './config'
import { Locations, Location, Flags } from './types'
import {
	exec,
	pathRelativeToConfigFile,
	getFlagsFromLocation,
	makeArrayIfIsNot,
	fill, decodeLocationFromPrefix, getPathFromVolume,
} from './utils'



export const forgetSingle = (name: string, to: string, location: Location, dryRun: boolean) => {
	const base = name + to.blue + ' : '
	const writer = new Writer(base + 'Removing old snapshots… ⏳')

	const backend = config.backends[to]
	const flags = getFlagsFromLocation(location, 'forget')

	const [type, value] = decodeLocationFromPrefix(location.from)
	let path: string
	switch (type) {
		case LocationFromPrefixes.Filesystem:
			path = pathRelativeToConfigFile(value)
			break
		case LocationFromPrefixes.DockerVolume:
			path = getPathFromVolume(value)
			break
	}

	if (flags.length == 0) {
		writer.done(base + 'Skipping, no policy declared')
		return
	}
	if (dryRun) flags.push('--dry-run')

	writer.replaceLn(base + 'Forgetting old snapshots… ⏳')
	const cmd = exec(
		'restic',
		['forget', '--path', path, '--prune', ...flags],
		{ env: getEnvFromBackend(backend) },
	)

	if (VERBOSE) console.log(cmd.out, cmd.err)
	writer.done(base + 'Done ✓'.green)
}

export const forgetLocation = (name: string, backup: Location, dryRun: boolean) => {
	const display = name.yellow + ' ▶ '
	const filler = fill(name.length + 3)
	let first = true

	for (const t of makeArrayIfIsNot(backup.to)) {
		const nameOrBlankSpaces: string = first ? display : filler
		forgetSingle(nameOrBlankSpaces, t, backup, dryRun)
		if (first) first = false
	}
}

export const forgetAll = (backups?: Locations, flags?: Flags) => {
	if (!backups) {
		backups = config.locations
	}

	console.log('\nRemoving old snapshots according to policy'.underline.grey)
	const dryRun = flags ? flags['dry-run'] : false
	if (dryRun) console.log('Running in dry-run mode, not touching data\n'.yellow)

	for (const [name, backup] of Object.entries(backups))
		forgetLocation(name, backup, dryRun)
}
