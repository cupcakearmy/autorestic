import { Writer } from 'clitastic'

import { config, VERBOSE } from './autorestic'
import { Backend, Backends, Locations } from './types'
import { exec, ConfigError, pathRelativeToConfigFile } from './utils'



const ALREADY_EXISTS = /(?=.*already)(?=.*config).*/

export const getPathFromBackend = (backend: Backend): string => {
	switch (backend.type) {
		case 'local':
			return pathRelativeToConfigFile(backend.path)
		case 'b2':
		case 'azure':
		case 'gs':
		case 's3':
		case 'sftp':
			return `${backend.type}:${backend.path}`
		case 'rest':
			throw new Error(`Unsupported backend type: "${backend.type}"`)
		default:
			throw new Error(`Unknown backend type.`)
	}
}

export const getEnvFromBackend = (backend: Backend) => {
	const { type, path, key, ...rest } = backend
	return {
		RESTIC_PASSWORD: key,
		RESTIC_REPOSITORY: getPathFromBackend(backend),
		...rest,
	}
}

export const getBackendsFromLocations = (locations: Locations): string[] => {
	const backends = new Set<string>()
	for (const to of Object.values(locations).map(location => location.to))
		Array.isArray(to) ? to.forEach(t => backends.add(t)) : backends.add(to)
	return Array.from(backends)
}

export const checkAndConfigureBackend = (name: string, backend: Backend) => {
	const writer = new Writer(name.blue + ' : ' + 'Configuring... ⏳')
	try {
		const env = getEnvFromBackend(backend)

		const { out, err } = exec('restic', ['init'], { env })

		if (err.length > 0 && !ALREADY_EXISTS.test(err))
			throw new Error(`Could not load the backend "${name}": ${err}`)

		if (VERBOSE && out.length > 0) console.log(out)

		writer.done(name.blue + ' : ' + 'Done ✓'.green)
	} catch (e) {
		writer.done(name.blue + ' : ' + 'Error ⚠️ ' + e.message.red)
	}
}

export const checkAndConfigureBackends = (backends?: Backends) => {
	if (!backends) {
		if (!config) throw ConfigError
		backends = config.backends
	}

	console.log('\nConfiguring Backends'.grey.underline)
	for (const [name, backend] of Object.entries(backends))
		checkAndConfigureBackend(name, backend)
}
