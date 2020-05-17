import { getRenderer } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/transport/renderer.js';
import { initJssCs } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/transport/setup-jss.js';initJssCs();
import { installTheme } from '/Users/nicco/Documents/git/autorestic/.codedoc/content/theme.ts';installTheme();
import { codeSelection } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/code/selection.js';codeSelection();
import { sameLineLengthInCodes } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/code/same-line-length.js';sameLineLengthInCodes();
import { initHintBox } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/code/line-hint/index.js';initHintBox();
import { initCodeLineRef } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/code/line-ref/index.js';initCodeLineRef();
import { initSmartCopy } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/code/smart-copy.js';initSmartCopy();
import { copyHeadings } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/heading/copy-headings.js';copyHeadings();
import { contentNavHighlight } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/page/contentnav/highlight.js';contentNavHighlight();
import { loadDeferredIFrames } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/transport/deferred-iframe.js';loadDeferredIFrames();
import { smoothLoading } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/transport/smooth-loading.js';smoothLoading();
import { tocHighlight } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/page/toc/toc-highlight.js';tocHighlight();
import { postNavSearch } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/page/toc/search/post-nav/index.js';postNavSearch();
import { ToCPrevNext } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/page/toc/prevnext/index.js';
import { CollapseControl } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/collapse/collapse-control.js';
import { GithubSearch } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/misc/github/search.js';
import { ToCToggle } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/page/toc/toggle/index.js';
import { DarkModeSwitch } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/components/darkmode/index.js';
import { ConfigTransport } from '/Users/nicco/Documents/git/autorestic/.codedoc/node_modules/@codedoc/core/dist/es6/transport/config.js';

const components = {
  'qcaKEY878Mn2dFQW/lSrDg==': ToCPrevNext,
  'fz894w7KG2/tX4kLbbA1Kg==': CollapseControl,
  '+SrlfVhZ/PRQ5WhUlZbTaA==': GithubSearch,
  'XsNW3ht5ee+RmVUActEo9g==': ToCToggle,
  'Y1WWvCKxkgk1yh8xbCfXqw==': DarkModeSwitch,
  'v641FmLj+AeGp0uuFTI6ug==': ConfigTransport
};

const renderer = getRenderer();
const ogtransport = window.__sdh_transport;
window.__sdh_transport = function(id, hash, props) {
  if (hash in components) {
    const target = document.getElementById(id);
    renderer.render(renderer.create(components[hash], props)).after(target);
    target.remove();
  }
  else if (ogtransport) ogtransport(id, hash, props);
}
