import { unlockRepo, waitForRepo } from '../lock'
import { Log } from '../logger'
import { Context } from '../models/context'
import { isResticAvailable } from '../restic'

export async function check(ctx: Context) {
  const l = Log.child({ cmd: 'check' })
  l.trace('starting')

  // Restic
  isResticAvailable()

  // Repos
  for (const repo of ctx.repos) {
    await waitForRepo(ctx, repo.name)
    try {
      await repo.init()
      await repo.check()
    } finally {
      unlockRepo(ctx, repo.name)
    }
  }

  l.trace('done')
}
