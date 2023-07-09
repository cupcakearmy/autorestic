import { Config } from '../config/schema/config'
import { CustomError } from '../errors'
import { Location } from './location'
import { Repository } from './repository'

export class Context {
  repos: Repository[]
  locations: Location[]

  constructor(public config: Config) {
    this.repos = Object.entries(config.repos).map(([name, r]) => new Repository(this, name, r))
    this.locations = Object.entries(config.locations).map(([name, l]) => new Location(this, name, l))
  }

  getRepo(name: string) {
    const repo = this.repos.find((r) => r.name === name)
    if (!repo) throw new CustomError(`could not find backend "${name}"`)
    return repo
  }
  getLocation(name: string) {
    const location = this.locations.find((l) => l.name === name)
    if (!location) throw new CustomError(`could not find location "${name}"`)
    return location
  }
}
