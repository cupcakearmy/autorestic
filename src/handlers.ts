import axios from 'axios'
import { Writer } from 'clitastic'
import { unlinkSync } from 'fs'
import { tmpdir } from 'os'
import { join, resolve } from 'path'

import { config, INSTALL_DIR, VERSION } from './autorestic'
import { checkAndConfigureBackends, getBackendsFromLocations, getEnvFromBackend } from './backend'
import { backupAll } from './backup'
import { forgetAll } from './forget'
import { Backends, Flags, Locations } from './types'
import {
	checkIfCommandIsAvailable,
	checkIfResticIsAvailable,
	downloadFile,
	exec,
	filterObjectByKey,
	ConfigError, makeArrayIfIsNot,
} from './utils'



export type Handlers = {
	[command: string]: (args: string[], flags: Flags) => void
}

const parseBackend = (flags: Flags): Backends => {
	if (!config) throw ConfigError
	if (!flags.all && !flags.backend)
		throw new Error(
			'No backends specified.'.red +
			'\n--all [-a]\t\t\t\tCheck all.' +
			'\n--backend [-b] myBackend\t\tSpecify one or more backend',
		)
	if (flags.all) return config.backends
	else {
		const backends = makeArrayIfIsNot<string>(flags.backend)
		for (const backend of backends)
			if (!config.backends[backend])
				throw new Error('Invalid backend: '.red + backend)
		return filterObjectByKey(config.backends, backends)
	}
}

const parseLocations = (flags: Flags): Locations => {
	if (!config) throw ConfigError
	if (!flags.all && !flags.location)
		throw new Error(
			'No locations specified.'.red +
			'\n--all [-a]\t\t\t\tBackup all.' +
			'\n--location [-l] site1\t\t\tSpecify one or more locations',
		)

	if (flags.all) {
		return config.locations
	} else {
		const locations = makeArrayIfIsNot<string>(flags.location)
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
		if (!config) throw ConfigError
		checkIfResticIsAvailable()
		const locations: Locations = parseLocations(flags)

		checkAndConfigureBackends(
			filterObjectByKey(config.backends, getBackendsFromLocations(locations)),
		)
		backupAll(locations)

		console.log('\nFinished!'.underline + ' 🎉')
	},
	restore(args, flags) {
		if (!config) throw ConfigError
		checkIfResticIsAvailable()

		if (!flags.to) {
			console.log(`You need to specify the restore path with --to`.red)
			return
		}

		const locations = parseLocations(flags)
		for (const [name, location] of Object.entries(locations)) {
			const baseText = name.green + '\t\t'
			const w = new Writer(baseText + `Starting...`)

			let backend: string = Array.isArray(location.to) ? location.to[0] : location.to
			if (flags.from) {
				if (!location.to.includes(flags.from)) {
					w.done(baseText + `Backend ${flags.from} is not a valid location for ${name}`.red)
					continue
				}
				backend = flags.from
				w.replaceLn(baseText + `Restoring from ${backend.blue}...`)
			} else if (Array.isArray(location.to) && location.to.length > 1) {
				w.replaceLn(baseText + `Restoring from ${backend.blue}...\tTo select a specific backend pass the ${'--from'.blue} flag`)
			}
			const env = getEnvFromBackend(config.backends[backend])

			exec(
				'restic',
				['restore', 'latest', '--path', resolve(location.from), '--target', flags.to],
				{ env },
			)
			w.done(name.green + '\t\tDone 🎉')
		}
	},
	forget(args, flags) {
		if (!config) throw ConfigError
		checkIfResticIsAvailable()
		const locations: Locations = parseLocations(flags)

		checkAndConfigureBackends(
			filterObjectByKey(config.backends, getBackendsFromLocations(locations)),
		)
		forgetAll(locations, flags)

		console.log('\nFinished!'.underline + ' 🎉')
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
			x32: '386',
			x64: 'amd64',
		}

		w.replaceLn('Downloading binary... 🌎')
		const name = `${json.name.replace(' ', '_')}_${process.platform}_${
			archMap[process.arch]
		}.bz2`
		const dl = json.assets.find((asset: any) => asset.name === name)
		if (!dl)
			return console.log(
				'Cannot get the right binary.'.red,
				'Please see https://bit.ly/2Y1Rzai',
			)

		const tmp = join(tmpdir(), name)
		const extracted = tmp.slice(0, -4) //without the .bz2

		await downloadFile(dl.browser_download_url, tmp)

		w.replaceLn('Decompressing binary... 📦')
		exec('bzip2', ['-dk', tmp])
		unlinkSync(tmp)

		w.replaceLn(`Moving to ${INSTALL_DIR} 🚙`)
		exec('chmod', ['+x', extracted])
		exec('mv', [extracted, INSTALL_DIR + '/restic'])

		w.done(
			`\nFinished! restic is installed under: ${INSTALL_DIR}`.underline + ' 🎉',
		)
	},
	uninstall() {
		for (const bin of ['restic', 'autorestic'])
			try {
				unlinkSync(INSTALL_DIR + '/' + bin)
				console.log(`Finished! ${bin} was uninstalled`)
			} catch (e) {
				console.log(`${bin} is already uninstalled`.red)
			}
	},
	async update() {
		checkIfResticIsAvailable()
		const w = new Writer('Checking for latest restic version... ⏳')
		exec('restic', ['self-update'])

		w.replaceLn('Checking for latest autorestic version... ⏳')
		const { data: json } = await axios({
			method: 'get',
			url:
				'https://api.github.com/repos/cupcakearmy/autorestic/releases/latest',
			responseType: 'json',
		})

		if (json.tag_name != VERSION) {
			const platformMap: { [key: string]: string } = {
				darwin: 'macos',
			}

			const name = `autorestic_${platformMap[process.platform] || process.platform}_${process.arch}`
			const dl = json.assets.find((asset: any) => asset.name === name)

			const to = INSTALL_DIR + '/autorestic'
			w.replaceLn('Downloading binary... 🌎')
			await downloadFile(dl.browser_download_url, to)

			exec('chmod', ['+x', to])
		}

		w.done('All up to date! 🚀')
	},
	version() {
		console.log('version'.grey, VERSION)
	},
}

export const help = () => {
	console.log(
		'\nAutorestic'.blue +
		` - ${VERSION} - Easy Restic CLI Utility` +
		'\n' +
		'\nOptions:'.yellow +
		`\n  -c, --config                                                          Specify config file. Default: .autorestic.yml` +
		'\n' +
		'\nCommands:'.yellow +
		'\n  check    [-b, --backend]  [-a, --all]                                 Check backends' +
		'\n  backup   [-l, --location] [-a, --all]                                 Backup all or specified locations' +
		'\n  forget   [-l, --location] [-a, --all] [--dry-run]                     Forget old snapshots according to declared policies' +
		'\n  restore  [-l, --location] [--from backend] [--to <out dir>]           Restore all or specified locations' +
		'\n' +
		'\n  exec     [-b, --backend]  [-a, --all] <command> -- [native options]   Execute native restic command' +
		'\n' +
		'\n  install                                                               install restic' +
		'\n  uninstall                                                             uninstall restic' +
		'\n  update                                                                update restic' +
		'\n  help                                                                  Show help' +
		'\n' +
		'\nExamples: '.yellow +
		'https://git.io/fjVbg' +
		'\n',
	)
}

export const error = () => {
	help()
	console.log(
		`Invalid Command:`.red.underline,
		`${process.argv.slice(2).join(' ')}`,
	)
}

export default handlers
