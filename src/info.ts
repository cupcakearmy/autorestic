import { config } from './autorestic'
import { fill, treeToString } from './utils'



const showAll = () => {
	console.log('\n\n' + fill(32, '_') + 'LOCATIONS:'.underline)
	for (const [key, data] of Object.entries(config.locations)) {
		console.log(`\n${key.blue.underline}:`)
		console.log(treeToString(
			data,
			['to:', 'from:', 'hooks:', 'options:'],
		))
	}

	console.log('\n\n' + fill(32, '_') + 'BACKENDS:'.underline)
	for (const [key, data] of Object.entries(config.backends)) {
		console.log(`\n${key.blue.underline}:`)
		console.log(treeToString(
			data,
			['type:', 'path:', 'key:'],
		))
	}
}

export default showAll