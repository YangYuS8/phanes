/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Phanes',
  tagline: 'Local-first offline runtime, resource-cache builder, and launcher framework',
  favicon: 'img/favicon.svg',

  url: 'https://geneden.gitlab.io',
  baseUrl: '/phanes/',
  trailingSlash: false,

  organizationName: 'Geneden',
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
          editUrl: 'https://gitlab.com/Geneden/phanes/-/edit/main/',
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
          href: 'https://gitlab.com/Geneden/phanes',
          label: 'GitLab',
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
            { label: 'GitLab', href: 'https://gitlab.com/Geneden/phanes' },
            { label: 'License', href: 'https://gitlab.com/Geneden/phanes/-/blob/main/LICENSE' },
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
