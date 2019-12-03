import { Writer } from 'clitastic'

import { config, VERBOSE } from './autorestic'
import { getEnvFromBackend } from './backend'
import { Locations, Location, ForgetPolicy, Flags } from './types'
import { exec, ConfigError } from './utils'

export const forgetSingle = (dryRun: boolean, name: string, from: string, to: string, policy: ForgetPolicy) => {
  if (!config) throw ConfigError
  const writer = new Writer(name + to.blue + ' : ' + 'Removing old spnapshots… ⏳')
  const backend = config.backends[to]
  const flags = [] as any[]
  for (const [name, value] of Object.entries(policy)) {
    flags.push(`--keep-${name}`)
    flags.push(value)
  }
  if (dryRun) {
    flags.push('--dry-run')
  }
  const env = getEnvFromBackend(backend)
  writer.appendLn(name + to.blue + ' : ' + 'Forgeting old snapshots… ⏳')
  const cmd = exec('restic', ['forget', '--path', from, '--prune', ...flags], {env})

  if (VERBOSE) console.log(cmd.out, cmd.err)
  writer.done(name + to.blue + ' : ' + 'Done ✓'.green)
}

export const forgetLocation = (dryRun: boolean, name: string, backup: Location, policy?: ForgetPolicy) => {
  const display = name.yellow + ' ▶ '
  if (!policy) {
    console.log(display + 'skipping, no policy declared')
  }
  else {
    if (Array.isArray(backup.to)) {
      let first = true
      for (const t of backup.to) {
        const nameOrBlankSpaces: string = first
          ? display
          : new Array(name.length + 3).fill(' ').join('')
        forgetSingle(dryRun, nameOrBlankSpaces, backup.from, t, policy)
        if (first) first = false
      }
    } else forgetSingle(dryRun, display, backup.from, backup.to, policy)
  }
  }

export const forgetAll = (dryRun: boolean, backups?: Locations) => {
  if (!config) throw ConfigError
  if (!backups) {
    backups = config.locations
  }

  console.log('\nRemoving old shapshots according to policy'.underline.grey)
  if (dryRun) console.log('Running in dry-run mode, not touching data\n'.yellow)

  for (const [name, backup] of Object.entries(backups)) {
    var policy = config.locations[name].keep
    forgetLocation(dryRun, name, backup, policy)
  }
}
