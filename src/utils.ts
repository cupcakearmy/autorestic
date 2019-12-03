import axios from 'axios'
import { spawnSync, SpawnSyncOptions } from 'child_process'
import { randomBytes } from 'crypto'
import { createWriteStream } from 'fs'
import { isAbsolute, resolve, dirname } from 'path'
import { CONFIG_FILE } from './config'



export const exec = (
	command: string,
	args: string[],
	{ env, ...rest }: SpawnSyncOptions = {},
) => {
	const cmd = spawnSync(command, args, {
		...rest,
		env: {
			...process.env,
			...env,
		},
	})

	const out = cmd.stdout && cmd.stdout.toString().trim()
	const err = cmd.stderr && cmd.stderr.toString().trim()

	return { out, err }
}

export const checkIfResticIsAvailable = () =>
	checkIfCommandIsAvailable(
		'restic',
		'Restic is not installed'.red +
		' https://restic.readthedocs.io/en/latest/020_installation.html#stable-releases',
	)

export const checkIfCommandIsAvailable = (cmd: string, errorMsg?: string) => {
	if (require('child_process').spawnSync(cmd).error)
		throw new Error(errorMsg ? errorMsg : `"${errorMsg}" is not installed`.red)
}

export const makeObjectKeysLowercase = (object: Object): any =>
	Object.fromEntries(
		Object.entries(object).map(([key, value]) => [key.toLowerCase(), value]),
	)


export function rand(length = 32): string {
	return randomBytes(length / 2).toString('hex')
}


export const singleToArray = <T>(singleOrArray: T | T[]): T[] =>
	Array.isArray(singleOrArray) ? singleOrArray : [singleOrArray]

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

		const stream = createWriteStream(to)

		const writer = file.pipe(stream)
		writer.on('close', () => {
			stream.close()
			res()
		})
	})

// Check if is an absolute path, otherwise get the path relative to the config file
export const pathRelativeToConfigFile = (path: string): string => isAbsolute(path)
	? path
	: resolve(dirname(CONFIG_FILE), path)

export const ConfigError = new Error('Config file not found')
