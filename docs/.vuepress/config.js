const { description } = require('../package')

module.exports = {
  // We're deploying to https://uber-go.github.io/cff/
  // so base should be /cff/.
  base: '/cff/',
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#title
   */
  title: 'cff',
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#description
   */
  description: description,

  dest: 'dist', // Publish built website to dist. We'll feed this to GitHub.

  /**
   * Extra tags to be injected to the page HTML `<head>`
   *
   * ref：https://v1.vuepress.vuejs.org/config/#head
   */
  head: [
    ['meta', { name: 'theme-color', content: '#3eaf7c' }],
    ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
    ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black' }]
  ],

  /**
   * Theme configuration, here is the default theme configuration for VuePress.
   *
   * ref：https://v1.vuepress.vuejs.org/theme/default-theme-config.html
   */
  themeConfig: {
    repo: 'uber-go/cff',
    editLinks: true,
    docsDir: 'docs',
    docsBranch: 'main',
    lastUpdated: true,
    nav: [
      {
        text: 'Guide',
        link: '/intro',
      },
      {
        text: 'API Reference',
        link: 'https://pkg.go.dev/go.uber.org/cff'
      }
    ],
    sidebar: [
      {
        title: 'Get Started',
        path: '/get-started/',
        children: [
          'get-started/flow.md',
        ],
      },
      {
        title: 'Introduction',
        path: '/intro',
        children: [
          'use-cases.md',
          'non-use-cases.md',
        ],
      },
      'setup.md',
      {
        title: 'Concepts',
        path: '/concepts',
        children: [
          'architecture.md',
        ],
      },
      'faq.md',
      {
        title: 'More',
        children: [
          'contributing.md',
          'changelog.md',
        ],
      },
    ]
  },

  /**
   * Apply plugins，ref：https://v1.vuepress.vuejs.org/zh/plugin/
   */
  plugins: [
    '@vuepress/plugin-back-to-top',
    '@vuepress/plugin-medium-zoom',
    'fulltext-search',
    'vuepress-plugin-mermaidjs',
  ]
}
