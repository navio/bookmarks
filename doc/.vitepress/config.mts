import { defineConfig } from 'vitepress'

function siteBase(): string {
  return process.env.VITEPRESS_BASE || '/'
}

const personalSiteIcon = {
  svg: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 0 1 0 20"/><path d="M12 2a15.3 15.3 0 0 0 0 20"/></svg>'
}

export default defineConfig({
  lang: 'en-US',
  title: 'bm',
  description: 'Bookmark directories from the terminal.',
  base: siteBase(),

  themeConfig: {
    nav: [
      { text: 'Getting Started', link: '/getting-started' },
      { text: 'Commands', link: '/commands' },
      { text: 'Usage', link: '/workflows' },
      { text: 'Store', link: '/store' }
    ],
    sidebar: [
      {
        text: 'Docs',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'Commands', link: '/commands' },
          { text: 'Usage Patterns', link: '/workflows' },
          { text: 'Store & Format', link: '/store' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/navio/bookmarks' },
      { icon: personalSiteIcon, link: 'https://alberto.pub' }
    ],
    footer: {
      message:
        '<a href="https://alberto.pub" target="_blank" rel="noreferrer">Alberto Navarro</a>'
    }
  }
})
