# Documentation Site Contract

Phanes documentation is managed as repository Markdown and published as a Docusaurus site through GitLab Pages.

## Source of truth

- Documentation source lives under `docs/`.
- Docusaurus uses `docs/intro.md` as the site entry page.
- ADRs live under `docs/adr/`.
- Contract docs live under `docs/contracts/`.
- Sidebar ordering is explicit in `sidebars.js` so important contract docs remain easy to find.

## Local commands

```bash
npm install
npm run start
npm run build
```

`npm run build` must succeed before GitLab Pages deployment.

## GitLab Pages

GitLab CI/CD builds the Docusaurus site and publishes it with the `pages` job.

The site is configured for the project Pages path:

```text
baseUrl: /phanes/
```

If the GitLab Pages URL or project path changes, update `docusaurus.config.js` and validate links with `npm run build`.

## Content boundaries

The documentation site must preserve the repository hard boundaries:

- no client patching instructions
- no anti-cheat/protection bypass guidance
- no official service authentication bypass guidance
- no required GC-Resources setup path
- no bundled copyrighted resource data
- runtime remains localhost-only by default
