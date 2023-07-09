import { Command, Help, Option, program } from '@commander-js/extra-typings'
import { loadConfig } from './config'
import { CustomError, NotImplemented } from './errors'
import { Log, LogLevel, setLevelFromFlag } from './logger'
import { check } from './cmd/check'
import { Context } from './models/context'
import { backup } from './cmd/backup'

export const helpConfig: Partial<Help> = {
  showGlobalOptions: true,
  sortOptions: true,
  sortSubcommands: true,
  helpWidth: 1,
}

program
  .name('autorestic')
  .description('configuration manager and runner for restic')
  .version('2.0.0-alpha.0')
  .configureHelp(helpConfig)
  .allowExcessArguments(false)
  .allowUnknownOption(false)

// Global options
program.option('-c, --config <file>', 'specify custom configuration file')
program.option('-v', 'verbosity', (_, previous) => previous + 1, 1)
program.addOption(new Option('--ci', 'CI mode').env('CI').default(false))

// Common Options
const specificLocation = new Option('-l, --location <locations...>', 'location name, multiple possible')
specificLocation.variadic = true
const allLocations = new Option('-a, --all', 'all locations')
const specificRepo = new Option('-r, --repository <names...>', 'repository name, multiple possible')
specificLocation.variadic = true
const allRepos = new Option('-a, --all', 'all repositories')

function mergeOptions<T extends {}>(local: T, p: Command) {
  const globals = p.optsWithGlobals() as { config?: string; verbosity: number; ci: boolean }
  return {
    ...globals,
    ...local,
  }
}

program.hook('preAction', (command) => {
  // @ts-ignore
  const v: number = command.opts().v
  setLevelFromFlag(v)
})

program
  .command('check')
  .description('check if the config is valid and sets up the repositories')
  .configureHelp(helpConfig)
  .action(async (options, p) => {
    const merged = mergeOptions(options, p)
    const config = await loadConfig(merged.config)
    const ctx = new Context(config)
    await check(ctx)
  })

program
  .command('backup')
  .description('create backups')
  .configureHelp(helpConfig)
  .addOption(specificLocation)
  .addOption(allLocations)
  .action(async (options, p) => {
    // throw new NotImplemented('backup')
    const merged = mergeOptions(options, p)
    const config = await loadConfig(merged.config)
    const ctx = new Context(config)
    await backup(ctx)
  })

program
  .command('exec')
  .description('execute arbitrary native restic commands for given repositories')
  .configureHelp(helpConfig)
  .addOption(specificRepo)
  .addOption(allRepos)
  .allowExcessArguments(true)
  .action((options, p) => {
    throw new NotImplemented('exec')
  })

program
  .command('forget')
  .description('forget snapshots according to the specified policies')
  .configureHelp(helpConfig)
  .addOption(specificLocation)
  .addOption(allLocations)
  // Pass natively
  // .option('--dry-run', 'do not write changes, show what would be affected')
  // .option('--prune', 'also prune repository')
  .action((options) => {
    throw new NotImplemented('backup')
  })

program
  .command('restore')
  .description('restore a snapshot to a given location')
  .option('--force', 'overwrite target folder')
  .option('--from <repository>', 'repository from which to restore')
  .option('--to <path>', 'path where to restore the data')
  .option('-l, --location <location>', 'location to be restored')
  .argument('[snapshot-id]', 'snapshot to be restored. if empty latest will be taken')
  .action(() => {
    throw new NotImplemented('restore')
  })

const self = new Command('self').description('utility commands for managing autorestic').configureHelp(helpConfig)
self.command('install').action(() => {
  throw new NotImplemented('install')
})
self.command('uninstall').action(() => {
  throw new NotImplemented('uninstall')
})
self.command('upgrade').action(() => {
  throw new NotImplemented('upgrade')
})
self.command('completion').action(() => {
  throw new NotImplemented('completion')
})
program.addCommand(self)

try {
  await program.parseAsync()
} catch (e) {
  if (e instanceof CustomError) {
    Log.fatal(e.message)
  } else if (e instanceof Error) {
    Log.fatal(`unknown error: ${e.message}`)
  }
  process.exit(1)
} finally {
  // TODO: Unlock
}
