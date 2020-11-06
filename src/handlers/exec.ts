import { getEnvFromBackend } from '../backend'
import { Flags } from '../types'
import { checkIfResticIsAvailable, exec as execCLI, parseBackend } from '../utils'

export default function exec({ backend, all }: Flags, args: string[]) {
  checkIfResticIsAvailable()
  const backends = parseBackend(backend, all)
  for (const [name, backend] of Object.entries(backends)) {
    console.log(`\n${name}:\n`.grey.underline)
    const env = getEnvFromBackend(backend)
    const { out, err } = execCLI('restic', args, { env })
    console.log(out, err)
  }
}
