import { checkAndConfigureBackends } from '../backend'
import { Flags } from '../types'
import { checkIfResticIsAvailable, parseBackend } from '../utils'

export default function check({ backend, all }: Flags) {
  checkIfResticIsAvailable()
  const backends = parseBackend(backend, all)
  checkAndConfigureBackends(backends)
}
