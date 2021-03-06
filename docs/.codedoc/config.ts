import { configuration } from '@codedoc/core'

export const config = configuration({
  src: {
    base: 'markdown',
  },
  dest: {
    html: './build',
    assets: './build',
    bundle: './_',
    styles: './_',
  },
  page: {
    title: {
      base: 'Autorestic',
    },
  },
  misc: {
    github: {
      user: 'cupcakearmy',
      repo: 'autorestic',
    },
  },
})
