import { Log } from '../logger'
import { Context } from '../models/context'

export async function backup(ctx: Context) {
  const log = Log.child({ cmd: 'check' })
  log.trace('starting')

  // Locations
  for (const location of ctx.locations) {
    await location.backup()
  }

  log.trace('done')
}
