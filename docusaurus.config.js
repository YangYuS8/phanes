/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Phanes',
  tagline: 'Local-first offline runtime, resource-cache builder, and launcher framework',
  favicon: 'img/favicon.svg',

  url: 'https://yangyus8.github.io',
  baseUrl: '/phanes/',
  trailingSlash: false,

  organizationName: 'YangYuS8',
  projectName: 'phanes',

  onBrokenLinks: 'throw',
  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'zh-Hans'],
    localeConfigs: {
      en: {
        label: 'English',
      },
      'zh-Hans': {
        label: '简体中文',
      },
    },
  },

  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
          editUrl: 'https://github.com/YangYuS8/phanes/edit/main/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],

  themeConfig: {
    navbar: {
      title: 'Phanes',
      logo: {
        alt: 'Phanes',
        src: 'img/favicon.svg',
      },
      items: [
        { to: '/', label: 'Docs', position: 'left' },
        {
          type: 'localeDropdown',
          position: 'right',
        },
        {
          href: 'https://github.com/YangYuS8/phanes',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Contracts',
          items: [
            { label: 'Module Boundaries', to: '/contracts/module-boundaries' },
            { label: 'Compliance Checklist', to: '/contracts/compliance-checklist' },
          ],
        },
        {
          title: 'Project',
          items: [
            { label: 'GitHub', href: 'https://github.com/YangYuS8/phanes' },
            { label: 'License', href: 'https://github.com/YangYuS8/phanes/blob/main/LICENSE' },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} The Phanes contributors. Licensed under Apache-2.0.`,
    },
    prism: {
      additionalLanguages: ['go', 'protobuf', 'bash', 'json', 'yaml'],
    },
  },
};

module.exports = config;
