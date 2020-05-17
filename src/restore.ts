import { Writer } from 'clitastic'
import { resolve } from 'path'

import { config } from './autorestic'
import { getEnvFromBackend } from './backend'
import { LocationFromPrefixes } from './config'
import { Backend } from './types'
import {
	checkIfDockerVolumeExistsOrFail,
	decodeLocationFromPrefix,
	exec,
	execPlain,
	getPathFromVolume,
} from './utils'



export const restoreToFilesystem = (from: string, to: string, backend: Backend) => {
	exec(
		'restic',
		['restore', 'latest', '--path', resolve(from), '--target', to],
		{ env: getEnvFromBackend(backend) },
	)
}

export const restoreToVolume = (volume: string, backend: Backend) => {
	const tmp = getPathFromVolume(volume)
	try {
		restoreToFilesystem(tmp, tmp, backend)
		try {
			checkIfDockerVolumeExistsOrFail(volume)
		} catch {
			execPlain(`docker volume create ${volume}`)
		}

		// For incremental backups. Unfortunately due to how the docker mounts work the permissions get lost.
		// execPlain(`docker run --rm -v ${volume}:/data -v ${tmp}:/backup alpine cp -aT /backup /data`)
		execPlain(`docker run --rm -v ${volume}:/data -v ${tmp}:/backup alpine tar xf /backup/archive.tar -C /data`)
	} finally {
		execPlain(`rm -rf ${tmp}`)
	}
}

export const restoreSingle = (locationName: string, from: string, to?: string) => {
	const location = config.locations[locationName]

	const baseText = locationName.green + '\t\t'
	const w = new Writer(baseText + `Restoring...`)

	let backendName: string = Array.isArray(location.to) ? location.to[0] : location.to
	if (from) {
		if (!location.to.includes(from)) {
			w.done(baseText + `Backend ${from} is not a valid location for ${locationName}`.red)
			return
		}
		backendName = from
		w.replaceLn(baseText + `Restoring from ${backendName.blue}...`)
	} else if (Array.isArray(location.to) && location.to.length > 1) {
		w.replaceLn(baseText + `Restoring from ${backendName.blue}...\tTo select a specific backend pass the ${'--from'.blue} flag`)
	}
	const backend = config.backends[backendName]

	const [type, value] = decodeLocationFromPrefix(location.from)
	switch (type) {

		case LocationFromPrefixes.Filesystem:
			if (!to) throw new Error(`You need to specify the restore path with --to`.red)
			restoreToFilesystem(value, to, backend)
			break

		case LocationFromPrefixes.DockerVolume:
			restoreToVolume(value, backend)
			break

	}
	w.done(locationName.green + '\t\tDone 🎉')
}

