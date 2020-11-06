import { build } from '@codedoc/core';

import { config } from './config';
import { installTheme$ } from './content/theme';
import { content } from './content';


build(config, content, installTheme$, {
  resolve: {
    modules: ['.codedoc/node_modules']
  },
  resolveLoader: {
    modules: ['.codedoc/node_modules']
  }
});
