import { restoreSingle } from '../restore'
import { Flags } from '../types'
import { checkIfResticIsAvailable, checkIfValidLocation } from '../utils'

export default function restore({ location, to, from }: Flags) {
  checkIfResticIsAvailable()
  checkIfValidLocation(location)
  restoreSingle(location, from, to)
}
