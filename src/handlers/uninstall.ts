import { unlinkSync } from 'fs'

import { INSTALL_DIR } from '..'

export function uninstall() {
  for (const bin of ['restic', 'autorestic'])
    try {
      unlinkSync(INSTALL_DIR + '/' + bin)
      console.log(`Finished! ${bin} was uninstalled`)
    } catch (e) {
      console.log(`${bin} is already uninstalled`.red)
    }
}
