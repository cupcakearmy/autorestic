import { Writer } from 'clitastic'
import { mkdirSync } from 'fs'

import { config, VERBOSE } from './autorestic'
import { getEnvFromBackend } from './backend'
import { LocationFromPrefixes } from './config'
import { Locations, Location, Backend } from './types'
import {
	exec,
	ConfigError,
	pathRelativeToConfigFile,
	getFlagsFromLocation,
	makeArrayIfIsNot,
	execPlain,
	MeasureDuration,
	fill,
	decodeLocationFromPrefix,
	hash, checkIfDockerVolumeExistsOrFail,
} from './utils'



export const backupFromFilesystem = (from: string, location: Location, backend: Backend, tags?: string[]) => {
	const path = pathRelativeToConfigFile(from)

	const cmd = exec(
		'restic',
		['backup', '.', ...getFlagsFromLocation(location, 'backup')],
		{ env: getEnvFromBackend(backend), cwd: path },
	)

	if (VERBOSE) console.log(cmd.out, cmd.err)
}

export const backupFromVolume = (volume: string, location: Location, backend: Backend) => {
	const tmp = pathRelativeToConfigFile(hash(volume))
	try {
		mkdirSync(tmp)
		checkIfDockerVolumeExistsOrFail(volume)

		// For incremental backups. Unfortunately due to how the docker mounts work the permissions get lost.
		// execPlain(`docker run --rm -v ${volume}:/data -v ${tmp}:/backup alpine cp -aT /data /backup`)
		execPlain(`docker run --rm -v ${volume}:/data -v ${tmp}:/backup alpine tar cf /backup/archive.tar -C /data .`)

		backupFromFilesystem(tmp, location, backend)
	} finally {
		execPlain(`rm -rf ${tmp}`)
	}
}

export const backupSingle = (name: string, to: string, location: Location) => {
	if (!config) throw ConfigError
	const delta = new MeasureDuration()
	const writer = new Writer(name + to.blue + ' : ' + 'Backing up... ⏳')

	try {
		const backend = config.backends[to]
		const [type, value] = decodeLocationFromPrefix(location.from)

		switch (type) {

			case LocationFromPrefixes.Filesystem:
				backupFromFilesystem(value, location, backend)
				break

			case LocationFromPrefixes.DockerVolume:
				backupFromVolume(value, location, backend)
				break

		}

		writer.done(`${name}${to.blue} : ${'Done ✓'.green} (${delta.finished(true)})`)
	} catch (e) {
		writer.done(`${name}${to.blue} : ${'Failed!'.red} (${delta.finished(true)}) ${e.message}`)
	}
}

export const backupLocation = (name: string, location: Location) => {
	const display = name.yellow + ' ▶ '
	const filler = fill(name.length + 3)
	let first = true

	if (location.hooks && location.hooks.before)
		for (const command of makeArrayIfIsNot(location.hooks.before)) {
			const cmd = execPlain(command, {})
			console.log(cmd.out, cmd.err)
		}

	for (const t of makeArrayIfIsNot(location.to)) {
		backupSingle(first ? display : filler, t, location)
		if (first) first = false
	}

	if (location.hooks && location.hooks.after)
		for (const command of makeArrayIfIsNot(location.hooks.after)) {
			const cmd = execPlain(command)
			console.log(cmd.out, cmd.err)
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
