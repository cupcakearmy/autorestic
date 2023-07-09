import { Logger } from 'pino'
import { z } from 'zod'
import { relativePath } from '../config/resolution'
import { Config } from '../config/schema/config'
import { RepositorySchema } from '../config/schema/repository'
import { ResticError } from '../errors'
import { Log } from '../logger'
import { execute } from '../restic'
import { Context } from './context'

export class Repository {
  l: Logger

  constructor(public ctx: Context, public name: string, public data: z.infer<typeof RepositorySchema>) {
    this.l = Log.child({ repository: this.name })
  }

  get repository(): string {
    switch (this.data.type) {
      case 'local':
        return relativePath(this.ctx.config.meta.path, this.data.path)
      case 'b2':
      case 'azure':
      case 'gs':
      case 's3':
      case 'sftp':
      case 'rclone':
      case 'swift':
      case 'rest':
        return `${this.data.type}:${this.data.path}`
        break
    }
  }

  get env() {
    return {
      ...this.data.env,
      RESTIC_PASSWORD: this.data.key,
      RESTIC_REPOSITORY: this.repository,
    }
  }

  /**
   * true if initialized
   * false if already initialized
   */
  async init(): Promise<boolean> {
    this.l.trace('initializing')
    const output = await execute({ command: 'restic', args: ['init'], env: this.env })
    if (!output.ok) {
      if (output.stderr.includes('config file already exists')) {
        this.l.debug('already initialized')
        return false
      }
      throw new ResticError([output.stderr])
    }
    this.l.debug('initialized repository')
    return true
  }

  async check() {
    this.l.trace('checking')
    const output = await execute({ command: 'restic', args: ['check'], env: this.env })
    if (!output.ok) throw new ResticError(['could not check repository', output.stderr])
  }
}
