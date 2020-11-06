import { checkAndConfigureBackendsForLocations } from '../backend'
import { backupAll } from '../backup'
import { Flags, Locations } from '../types'
import { checkIfResticIsAvailable, parseLocations } from '../utils'

export default function backup({ location, all }: Flags) {
  checkIfResticIsAvailable()
  const locations: Locations = parseLocations(location, all)
  checkAndConfigureBackendsForLocations(locations)
  backupAll(locations)

  console.log('\nFinished!'.underline + ' 🎉')
}
