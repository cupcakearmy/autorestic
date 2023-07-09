import Pino, { type LoggerOptions } from 'pino'
import Pretty from 'pino-pretty'

// https://getpino.io/#/docs/api?id=loggerlevel-string-gettersetter
export enum LogLevel {
  Trace = 'trace',
  Debug = 'debug',
  Info = 'info',
  Warn = 'warn',
  Error = 'error',
  Fatal = 'fatal',
  Silent = 'silent',
}

const pretty = !process.env.CI
const options: LoggerOptions = {
  base: undefined,
  level: LogLevel.Info,
}

export const Log = pretty ? Pino(options, Pretty({ colorize: true })) : Pino(options)

export function setLevelFromFlag(flag: number) {
  switch (flag) {
    case 1:
      Log.level = LogLevel.Info
      break
    case 2:
      Log.level = LogLevel.Debug
      break
    case 3:
      Log.level = LogLevel.Trace
      break
    default:
      Log.error('invalid logging level')
  }
}
