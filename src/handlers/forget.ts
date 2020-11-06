import { checkAndConfigureBackendsForLocations } from '../backend'
import { forgetAll } from '../forget'
import { Flags, Locations } from '../types'
import { checkIfResticIsAvailable, parseLocations } from '../utils'

export default function forget({ location, all, dryRun }: Flags) {
  checkIfResticIsAvailable()
  const locations: Locations = parseLocations(location, all)
  checkAndConfigureBackendsForLocations(locations)
  forgetAll(locations, dryRun)

  console.log('\nFinished!'.underline + ' 🎉')
}
