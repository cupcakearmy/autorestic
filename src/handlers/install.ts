import { join } from 'path'
import { chmodSync, renameSync, unlinkSync } from 'fs'
import { tmpdir } from 'os'

import axios from 'axios'
import { Writer } from 'clitastic'

import { INSTALL_DIR } from '..'
import { checkIfCommandIsAvailable, checkIfResticIsAvailable, downloadFile, exec } from '../utils'

export default async function install() {
  try {
    checkIfResticIsAvailable()
    console.log('Restic is already installed')
    return
  } catch {}

  const w = new Writer('Checking latest version... ⏳')
  checkIfCommandIsAvailable('bzip2')
  const { data: json } = await axios({
    method: 'get',
    url: 'https://api.github.com/repos/restic/restic/releases/latest',
    responseType: 'json',
  })

  const archMap: { [a: string]: string } = {
    x32: '386',
    x64: 'amd64',
  }

  w.replaceLn('Downloading binary... 🌎')
  const name = `${json.name.replace(' ', '_')}_${process.platform}_${archMap[process.arch]}.bz2`
  const dl = json.assets.find((asset: any) => asset.name === name)
  if (!dl) return console.log('Cannot get the right binary.'.red, 'Please see https://bit.ly/2Y1Rzai')

  const tmp = join(tmpdir(), name)
  const extracted = tmp.slice(0, -4) //without the .bz2

  await downloadFile(dl.browser_download_url, tmp)

  w.replaceLn('Decompressing binary... 📦')
  exec('bzip2', ['-dk', tmp])
  unlinkSync(tmp)

  w.replaceLn(`Moving to ${INSTALL_DIR} 🚙`)
  chmodSync(extracted, 0o755)
  renameSync(extracted, INSTALL_DIR + '/restic')

  w.done(`\nFinished! restic is installed under: ${INSTALL_DIR}`.underline + ' 🎉')
}
