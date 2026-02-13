import { defineConfig } from 'vitepress'

function pagesBase(): string {
  const envBase = process.env.VITEPRESS_BASE
  if (envBase) return envBase

  if (process.env.GITHUB_ACTIONS) {
    const repo = process.env.GITHUB_REPOSITORY?.split('/')[1]
    if (repo) return `/${repo}/`
  }

  return '/'
}

export default defineConfig({
  lang: 'en-US',
  title: 'bm',
  description: 'Bookmark directories from the terminal.',
  base: pagesBase(),

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
