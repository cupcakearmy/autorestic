import { Logger } from 'pino'
import { z } from 'zod'
import { LocationSchema } from '../config/schema/location'
import { Log } from '../logger'
import { asArray } from '../utils/array'
import { Context } from './context'
import { execute } from '../restic'

export class Location {
  l: Logger

  constructor(public ctx: Context, public name: string, public data: z.infer<typeof LocationSchema>) {
    this.l = Log.child({ location: name })
  }

  async backup() {
    this.l.trace('backing up location')
    for (const name of asArray(this.data.to)) {
      const repo = this.ctx.getRepo(name)
      this.l.debug(repo.name)
      await execute({
        command: 'restic',
        args: ['backup', '--dry-run'],
        env: repo.env,
      })
    }
  }
}
