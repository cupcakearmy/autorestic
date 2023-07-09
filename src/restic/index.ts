import { execFile } from 'node:child_process'
import { Log } from '../logger'
import { BinaryNotAvailable } from '../errors'

export type ExecutionContext = {
  command: string
  args?: string[]
  env?: Record<string, string>
}
export async function execute({
  env,
  args,
  command,
}: ExecutionContext): Promise<{ code: number; stderr: string; stdout: string; ok: boolean }> {
  return new Promise((resolve) => {
    execFile(command, args ?? [], { env }, (err, stdout, stderr) => {
      const code = err?.code ?? 0
      resolve({
        code,
        ok: code === 0,
        stderr,
        stdout,
      })
    })
  })
}

export async function isBinaryAvailable(command: string): Promise<boolean> {
  const l = Log.child({ command })
  try {
    l.trace('checking if command is installed')
    const result = await execute({ command })
    return result.ok
  } catch {
    l.trace('not installed')
    return false
  }
}

export async function isResticAvailable() {
  const bin = 'restic'
  const installed = await isBinaryAvailable(bin)
  if (!installed) throw new BinaryNotAvailable(bin)
}
