import { defineConfig } from 'vitepress'

function siteBase(): string {
  return process.env.VITEPRESS_BASE || '/'
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
    socialLinks: [{ icon: 'github', link: 'https://github.com/navio/bookmarks' }]
  }
})
