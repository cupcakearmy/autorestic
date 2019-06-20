import { readFileSync, writeFileSync } from 'fs'
import yaml from 'js-yaml'
import { CONFIG_FILE } from './autorestic'
import { Backend, Config } from './types'
import { makeObjectKeysLowercase, rand } from './utils'


export const normalizeAndCheckBackends = (config: Config) => {
	config.backends = makeObjectKeysLowercase(config.backends)

	for (const [name, { type, path, key, ...rest }] of Object.entries(config.backends)) {

		if (!type || !path) throw new Error(`The backend "${name}" is missing some required attributes`)

		const tmp: any = {
			type,
			path,
			key: key || rand(128),
		}
		for (const [key, value] of Object.entries(rest))
			tmp[key.toUpperCase()] = value

		config.backends[name] = tmp as Backend
	}
}


export const normalizeAndCheckBackups = (config: Config) => {
	config.locations = makeObjectKeysLowercase(config.locations)
	const backends = Object.keys(config.backends)

	const checkDestination = (backend: string, backup: string) => {
		if (!backends.includes(backend))
			throw new Error(`Cannot find the backend "${backend}" for "${backup}"`)
	}

	for (const [name, { from, to, ...rest }] of Object.entries(config.locations)) {
		if (!from || !to) throw new Error(`The backup "${name}" is missing some required attributes`)

		if (Array.isArray(to))
			for (const t of to)
				checkDestination(t, name)
		else
			checkDestination(to, name)
	}
}


export const init = (): Config => {
	const raw: Config = makeObjectKeysLowercase(yaml.safeLoad(readFileSync(CONFIG_FILE).toString()))

	normalizeAndCheckBackends(raw)
	normalizeAndCheckBackups(raw)

	writeFileSync(CONFIG_FILE, yaml.safeDump(raw))

	return raw
}
