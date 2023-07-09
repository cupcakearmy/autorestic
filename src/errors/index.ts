function formatLines(lines: string[]) {
  return lines.map((p) => `  ▶ ${p}`).join('\n')
}

export class CustomError extends Error {}

export class InvalidEnvFileLine extends CustomError {
  constructor(line: string) {
    super(`invalid env file line: "${line}"`)
  }
}

export class NotImplemented extends CustomError {
  constructor(functionality: string) {
    super(`not implemented: ${functionality}`)
  }
}

export class ConfigFileNotFound extends CustomError {
  constructor(paths: string[]) {
    super(`could not locate config file.\nthe following paths were tried:\n${formatLines(paths)}`)
  }
}

export class InvalidConfigFile extends CustomError {
  constructor(errors: string[]) {
    super(`could not parse the config file.\n${formatLines(errors)}`)
  }
}

export class BinaryNotAvailable extends CustomError {
  constructor(binary: string) {
    super(`binary "${binary}" is not available in $PATH`)
  }
}

export class ResticError extends CustomError {
  constructor(errors: string[]) {
    super(`internal restic error.\n${formatLines(errors)}`)
  }
}

export class LockfileAlreadyLocked extends CustomError {
  constructor(repo: string) {
    super(`cannot acquire lock for repository "${repo}", already in use`)
  }
}
