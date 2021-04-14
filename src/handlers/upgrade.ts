import { chmodSync } from 'fs'

import axios from 'axios'
import { Writer } from 'clitastic'
import semver from 'semver'

import { INSTALL_DIR, VERSION } from '..'
import { checkIfResticIsAvailable, downloadFile, exec } from '../utils'

export async function upgrade() {
  checkIfResticIsAvailable()
  const w = new Writer('Checking for latest restic version... ⏳')
  exec('restic', ['self-update'])

  w.replaceLn('Checking for latest autorestic version... ⏳')
  const { data: json } = await axios({
    method: 'get',
    url: 'https://api.github.com/repos/cupcakearmy/autorestic/releases/latest',
    responseType: 'json',
  })

  const latest = semver.coerce(json.tag_name)
  const current = semver.coerce(VERSION)
  if (!latest || !current) throw new Error('Could not parse versions numbers.')
  if (semver.gt(latest, current)) {
    if (semver.major(latest) === semver.major(current)) {
      // Update to compatible
      const platformMap: { [key: string]: string } = {
        darwin: 'macos',
      }

      const name = `autorestic_${platformMap[process.platform] || process.platform}_${process.arch}`
      const dl = json.assets.find((asset: any) => asset.name === name)

      const to = INSTALL_DIR + '/autorestic'
      w.replaceLn('Downloading binary... 🌎')
      await downloadFile(dl.browser_download_url, to)

      chmodSync(to, 0o755)
    } else {
      w.appendLn('Newer major version available, will not install automatically.')
    }
  }

  w.done('All up to date! 🚀')
}
