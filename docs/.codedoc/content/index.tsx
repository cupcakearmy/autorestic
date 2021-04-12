import { RendererLike } from '@connectv/html'
import { File } from 'rxline/fs'
import {
  Page,
  Meta,
  ContentNav,
  Fonts,
  ToC,
  GithubSearch$,
} from '@codedoc/core/components'

import { config } from '../config'
import { Header } from './header'
import { Footer } from './footer'

export function content(
  _content: HTMLElement,
  toc: HTMLElement,
  renderer: RendererLike<any, any>,
  file: File<string>
) {
  return (
    <Page
      title={config.page.title.extractor(_content, config, file)}
      favicon={config.page.favicon}
      meta={<Meta {...config.page.meta} />}
      fonts={<Fonts {...config.page.fonts} />}
      scripts={config.page.scripts}
      stylesheets={config.page.stylesheets}
      header={<Header {...config} />}
      footer={<Footer {...config} />}
      toc={
        <ToC
          default={'open'}
          search={
            config.misc?.github ? (
              <GithubSearch$
                repo={config.misc.github.repo}
                user={config.misc.github.user}
                root={config.src.base}
                pick={config.src.pick.source}
                drop={config.src.drop.source}
              />
            ) : (
              false
            )
          }
        >
          {toc}
        </ToC>
      }
    >
      {_content}
      <ContentNav content={_content} />
    </Page>
  )
}
