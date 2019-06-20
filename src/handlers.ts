import axios from 'axios'
import { Writer } from 'clitastic'
import { createWriteStream, unlinkSync } from 'fs'
import { arch, platform, tmpdir } from 'os'
import { join, resolve } from 'path'

import { config, INSTALL_DIR, CONFIG_FILE } from './autorestic'
import { checkAndConfigureBackends, getEnvFromBackend } from './backend'
import { backupAll } from './backup'
import { Backends, Flags, Locations } from './types'
import { checkIfCommandIsAvailable, checkIfResticIsAvailable, exec, filterObjectByKey, singleToArray } from './utils'

export type Handlers = { [command: string]: (args: string[], flags: Flags) => void }

const parseBackend = (flags: Flags): Backends => {
	if (!flags.all && !flags.backend)
		throw new Error('No backends specified.'.red
			+ '\n--all [-a]\t\t\t\tCheck all.'
			+ '\n--backend [-b] myBackend\t\tSpecify one or more backend',
		)
	if (flags.all)
		return config.backends
	else {
		const backends = singleToArray<string>(flags.backend)
		for (const backend of backends)
			if (!config.backends[backend])
				throw new Error('Invalid backend: '.red + backend)
		return filterObjectByKey(config.backends, backends)
	}
}

const parseLocations = (flags: Flags): Locations => {
	if (!flags.all && !flags.location)
		throw new Error('No locations specified.'.red
			+ '\n--all [-a]\t\t\t\tBackup all.'
			+ '\n--location [-l] site1\t\t\tSpecify one or more locations',
		)

	if (flags.all) {
		return config.locations
	} else {
		const locations = singleToArray<string>(flags.location)
		for (const location of locations)
			if (!config.locations[location])
				throw new Error('Invalid location: '.red + location)
		return filterObjectByKey(config.locations, locations)
	}
}

const handlers: Handlers = {
	check(args, flags) {
		checkIfResticIsAvailable()
		const backends = parseBackend(flags)
		checkAndConfigureBackends(backends)
	},
	backup(args, flags) {
		checkIfResticIsAvailable()
		const locations: Locations = parseLocations(flags)

		const backends = new Set<string>()
		for (const to of Object.values(locations).map(location => location.to))
			Array.isArray(to) ? to.forEach(t => backends.add(t)) : backends.add(to)

		checkAndConfigureBackends(filterObjectByKey(config.backends, Array.from(backends)))
		backupAll(locations)

		console.log('\nFinished!'.underline + ' 🎉')
	},
	restore(args, flags) {
		checkIfResticIsAvailable()
		const locations = parseLocations(flags)
		for (const [name, location] of Object.entries(locations)) {
			const w = new Writer(name.green + `\t\tRestoring... ⏳`)
			const env = getEnvFromBackend(config.backends[Array.isArray(location.to) ? location.to[0] : location.to])

			exec(
				'restic',
				['restore', 'latest', '--path', resolve(location.from), ...args],
				{ env },
			)
			w.done(name.green + '\t\tDone 🎉')
		}
	},
	exec(args, flags) {
		checkIfResticIsAvailable()
		const backends = parseBackend(flags)
		for (const [name, backend] of Object.entries(backends)) {
			console.log(`\n${name}:\n`.grey.underline)
			const env = getEnvFromBackend(backend)

			const { out, err } = exec('restic', args, { env })
			console.log(out, err)
		}
	},
	async install() {
		try {
			checkIfResticIsAvailable()
			console.log('Restic is already installed')
			return
		} catch (e) {
		}

		const w = new Writer('Checking latest version... ⏳')
		checkIfCommandIsAvailable('bzip2')
		const { data: json } = await axios({
			method: 'get',
			url: 'https://api.github.com/repos/restic/restic/releases/latest',
			responseType: 'json',
		})

		const archMap: { [a: string]: string } = {
			'x32': '386',
			'x64': 'amd64',
		}

		w.replaceLn('Downloading binary... 🌎')
		const name = `${json.name.replace(' ', '_')}_${platform()}_${archMap[arch()]}.bz2`
		const dl = json.assets.find((asset: any) => asset.name === name)
		if (!dl) return console.log(
			'Cannot get the right binary.'.red,
			'Please see https://bit.ly/2Y1Rzai',
		)

		const { data: file } = await axios({
			method: 'get',
			url: dl.browser_download_url,
			responseType: 'stream',
		})

		const from = join(tmpdir(), name)
		const to = from.slice(0, -4)

		w.replaceLn('Decompressing binary... 📦')
		const stream = createWriteStream(from)
		await new Promise(res => {
			const writer = file.pipe(stream)
			writer.on('close', res)
		})
		stream.close()

		w.replaceLn(`Moving to ${INSTALL_DIR} 🚙`)
		// TODO: Native bz2
		// Decompress
		exec('bzip2', ['-dk', from])
		// Remove .bz2
		exec('chmod', ['+x', to])
		exec('mv', [to, INSTALL_DIR + '/restic'])

		unlinkSync(from)

		w.done(`\nFinished! restic is installed under: ${INSTALL_DIR}`.underline + ' 🎉')
	},
	uninstall() {
		try {
			unlinkSync(INSTALL_DIR + '/restic')
			console.log(`Finished! restic was uninstalled`)
		} catch (e) {
			console.log('restic is already uninstalled'.red)
		}
	},
	update() {
		checkIfResticIsAvailable()
		const w = new Writer('Checking for new restic version... ⏳')
		exec('restic', ['self-update'])
		w.done('All up to date! 🚀')
	},
}

export const help = () => {
	console.log('\nAutorestic'.blue + ' - Easy Restic CLI Utility'
		+ '\n'
		+ '\nOptions:'.yellow
		+ `\n  -c, --config                                                          Specify config file. Default: ${CONFIG_FILE}`
		+ '\n'
		+ '\nCommands:'.yellow
		+ '\n  check    [-b, --backend]  [-a, --all]                                 Check backends'
		+ '\n  backup   [-l, --location] [-a, --all]                                 Backup all or specified locations'
		+ '\n  restore  [-l, --location] [-- --target <out dir>]                     Check backends'
		+ '\n'
		+ '\n  exec     [-b, --backend]  [-a, --all] <command> -- [native options]   Execute native restic command'
		+ '\n'
		+ '\n  install                                                               install restic'
		+ '\n  uninstall                                                             uninstall restic'
		+ '\n  update                                                                update restic'
		+ '\n  help                                                                  Show help'
		+ '\n'
		+ '\nExamples: '.yellow + 'https://git.io/fjVbg'
		+ '\n',
	)
}
export const error = () => {
	help()
	console.log(`Invalid Command:`.red.underline, `${process.argv.slice(2).join(' ')}`)
}

export default handlers