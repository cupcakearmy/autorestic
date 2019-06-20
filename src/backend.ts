import { Writer } from 'clitastic'

import { config, VERBOSE } from './autorestic'
import { Backend, Backends } from './types'
import { exec } from './utils'

const ALREADY_EXISTS = /(?=.*exists)(?=.*already)(?=.*config).*/


export const getPathFromBackend = (backend: Backend): string => {
	switch (backend.type) {
		case 'local':
			return backend.path
		case 'b2':
		case 'azure':
		case 'gs':
		case 's3':
			return `${backend.type}:${backend.path}`
		case 'sftp':
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


export const checkAndConfigureBackend = (name: string, backend: Backend) => {
	const writer = new Writer(name.blue + ' : ' + 'Configuring... ⏳')
	const env = getEnvFromBackend(backend)

	const { out, err } = exec('restic', ['init'], { env })

	if (err.length > 0 && !ALREADY_EXISTS.test(err))
		throw new Error(`Could not load the backend "${name}": ${err}`)

	if (VERBOSE && out.length > 0) console.log(out)

	writer.done(name.blue + ' : ' + 'Done ✓'.green)
}


export const checkAndConfigureBackends = (backends: Backends = config.backends) => {
	console.log('\nConfiguring Backends'.grey.underline)
	for (const [name, backend] of Object.entries(backends))
		checkAndConfigureBackend(name, backend)
}