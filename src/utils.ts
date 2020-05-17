import axios from 'axios'
import { spawnSync, SpawnSyncOptions } from 'child_process'
import { createHash, randomBytes } from 'crypto'
import { createWriteStream, renameSync, unlinkSync } from 'fs'
import { homedir, tmpdir } from 'os'
import { dirname, isAbsolute, join, resolve } from 'path'
import { Duration, Humanizer } from 'uhrwerk'

import { CONFIG_FILE, LocationFromPrefixes } from './config'
import { Location } from './types'



export const exec = (command: string, args: string[], { env, ...rest }: SpawnSyncOptions = {}) => {
	const { stdout, stderr, status } = spawnSync(command, args, {
		...rest,
		env: {
			...process.env,
			...env,
		},
	})

	const out = stdout && stdout.toString().trim()
	const err = stderr && stderr.toString().trim()

	return { out, err, status }
}

export const execPlain = (command: string, opt: SpawnSyncOptions = {}) => {
	const split = command.split(' ')
	if (split.length < 1) throw new Error(`The command ${command} is not valid`.red)

	return exec(split[0], split.slice(1), opt)
}

export const checkIfResticIsAvailable = () =>
	checkIfCommandIsAvailable(
		'restic',
		'restic is not installed'.red +
		'\nEither run ' + 'autorestic install'.green +
		'\nOr go to https://restic.readthedocs.io/en/latest/020_installation.html#stable-releases',
	)

export const checkIfCommandIsAvailable = (cmd: string, errorMsg?: string) => {
	if (spawnSync(cmd).error)
		throw new Error(errorMsg ? errorMsg : `"${cmd}" is not installed`.red)
}

export const makeObjectKeysLowercase = (object: Object): any =>
	Object.fromEntries(
		Object.entries(object).map(([key, value]) => [key.toLowerCase(), value]),
	)


export function rand(length = 32): string {
	return randomBytes(length / 2).toString('hex')
}


export const filterObject = <T>(
	obj: { [key: string]: T },
	filter: (item: [string, T]) => boolean,
): { [key: string]: T } =>
	Object.fromEntries(Object.entries(obj).filter(filter))

export const filterObjectByKey = <T>(
	obj: { [key: string]: T },
	keys: string[],
) => filterObject(obj, ([key]) => keys.includes(key))

export const downloadFile = async (url: string, to: string) =>
	new Promise<void>(async res => {
		const { data: file } = await axios({
			method: 'get',
			url: url,
			responseType: 'stream',
		})

		const tmp = join(tmpdir(), rand(64))
		const stream = createWriteStream(tmp)

		const writer = file.pipe(stream)
		writer.on('close', () => {
			stream.close()
			try {
				// Delete file if already exists. Needed if the binary wants to replace itself.
				// Unix does not allow to overwrite a file that is being executed, but you can remove it and save other one at its place
				unlinkSync(to)
			} catch {
			}
			renameSync(tmp, to)
			res()
		})
	})

// Check if is an absolute path, otherwise get the path relative to the config file
export const pathRelativeToConfigFile = (path: string): string => isAbsolute(path)
	? path
	: resolve(dirname(CONFIG_FILE), path)

export const resolveTildePath = (path: string): string | null =>
	(path.length === 0 || path[0] !== '~')
		? null
		: join(homedir(), path.slice(1))

export const getFlagsFromLocation = (location: Location, command?: string): string[] => {
	if (!location.options) return []

	const all = {
		...location.options.global,
		...(location.options[command || ''] || {}),
	}

	let flags: string[] = []
	// Map the flags to an array for the exec function.
	for (let [flag, values] of Object.entries(all))
		for (const value of makeArrayIfIsNot(values)) {
			const stringValue = String(value)
			const resolvedTilde = resolveTildePath(stringValue)
			flags = [...flags, `--${String(flag)}`, resolvedTilde === null ? stringValue : resolvedTilde]
		}

	return flags
}

export const makeArrayIfIsNot = <T>(maybeArray: T | T[]): T[] => Array.isArray(maybeArray) ? maybeArray : [maybeArray]

export const fill = (length: number, filler = ' '): string => new Array(length).fill(filler).join('')

export const capitalize = (string: string): string => string.charAt(0).toUpperCase() + string.slice(1)

export const treeToString = (obj: Object, highlight = [] as string[]): string => {
	let cleaned = JSON.stringify(obj, null, 2)
		.replace(/[{}"\[\],]/g, '')
		.replace(/^ {2}/mg, '')
		.replace(/\n\s*\n/g, '\n')
		.trim()

	for (const word of highlight)
		cleaned = cleaned.replace(word, capitalize(word).green)

	return cleaned
}


export class MeasureDuration {
	private static Humanizer: Humanizer = [
		[d => d.hours() > 0, d => `${d.hours()}h ${d.minutes()}min`],
		[d => d.minutes() > 0, d => `${d.minutes()}min ${d.seconds()}s`],
		[d => d.seconds() > 0, d => `${d.seconds()}s`],
		[() => true, d => `${d.milliseconds()}ms`],
	]

	private start = Date.now()


	finished(human?: false): number
	finished(human?: true): string
	finished(human?: boolean): number | string {
		const delta = Date.now() - this.start

		return human
			? new Duration(delta, 'ms').humanize(MeasureDuration.Humanizer)
			: delta
	}

}


export const decodeLocationFromPrefix = (from: string): [LocationFromPrefixes, string] => {
	const firstDelimiter = from.indexOf(':')
	if (firstDelimiter === -1) return [LocationFromPrefixes.Filesystem, from]

	const type = from.substr(0, firstDelimiter)
	const value = from.substr(firstDelimiter + 1)

	switch (type.toLowerCase()) {
		case 'volume':
			return [LocationFromPrefixes.DockerVolume, value]
		case 'path':
			return [LocationFromPrefixes.Filesystem, value]
		default:
			throw new Error(`Could not decode the location from: ${from}`.red)
	}
}

export const hash = (plain: string): string => createHash('sha1').update(plain).digest().toString('hex')

export const getPathFromVolume = (volume: string) => pathRelativeToConfigFile(hash(volume))

export const checkIfDockerVolumeExistsOrFail = (volume: string) => {
	const cmd = exec('docker', [
		'volume', 'inspect', volume,
	])
	if (cmd.err.length > 0)
		throw new Error('Volume not found')
}
