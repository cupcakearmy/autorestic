import { runCron } from '../cron'
import { checkIfResticIsAvailable } from '../utils'

export function cron() {
  checkIfResticIsAvailable()
  runCron()
}
