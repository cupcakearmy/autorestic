import { funcTransport } from '@connectv/sdh/transport';
import { useTheme } from '@codedoc/core/transport';

import { theme } from '../theme';


export function installTheme() { useTheme(theme); }
export const installTheme$ = /*#__PURE__*/funcTransport(installTheme);
