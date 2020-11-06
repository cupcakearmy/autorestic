import { CodedocConfig } from '@codedoc/core';
import { Header as _Header, GithubButton, Watermark } from '@codedoc/core/components';


export function Header(config: CodedocConfig, renderer: any) {
  return (
    <_Header>{config.misc?.github ?
      <fragment>
        <GithubButton action={config.misc.github.action || 'Star'}
          repo={config.misc.github.repo}
          user={config.misc.github.user}
          large={config.misc.github.large === true}
          count={config.misc.github.count !== false}
          standardIcon={config.misc.github.standardIcon !== false}/>
        <br/><br/>
      </fragment>
      : ''}
      <Watermark/>
    </_Header>
  )
}
