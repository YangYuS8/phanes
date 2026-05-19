# Compliance Checklist

Use this checklist for design reviews, implementation reviews, and handoffs.

## Localhost-only runtime

- [ ] Default bind host is `127.0.0.1`.
- [ ] `0.0.0.0` is rejected.
- [ ] Non-loopback/public bind addresses are rejected.
- [ ] Runtime status reports the actual bind host and port.
- [ ] No public server hosting feature is introduced without maintainer-approved scope change.

## Runtime does not download resources

- [ ] Runtime config has no remote resource URL fields.
- [ ] Missing/invalid cache fails with a diagnostic instead of fetching.
- [ ] Runtime startup does not invoke builder.
- [ ] Runtime package does not introduce remote resource clients.

## No required GC-Resources

- [ ] Builder supports `embedded-minimal` and local inputs without GC-Resources.
- [ ] `external-import` is optional and explicit.
- [ ] Docs do not tell users to download GC-Resources as a requirement.
- [ ] Tests do not depend on remote resource repositories.

## No copyrighted resource bundling

- [ ] Fixtures are tiny, synthetic, or otherwise safe to redistribute.
- [ ] Builder output is generated from user-provided local inputs.
- [ ] Repository does not commit full copyrighted game resource data.

## No unsafe bypass work

- [ ] No anti-cheat bypass.
- [ ] No client protection bypass.
- [ ] No official service authentication bypass.
- [ ] No unauthorized online interaction.
- [ ] No commercial service replacement.
- [ ] No deep client patching instructions.
- [ ] Launcher does not silently mutate system proxy/settings.

## Clean-room/reference posture

- [ ] Grasscutter and Cultivation are used for conceptual reference only.
- [ ] No copied source code from GPL/AGPL/reference projects without explicit license review and maintainer approval.
- [ ] Package layout and handler design are not blindly cloned from reference projects.

## Contract-first workflow

- [ ] Contract updated before implementation.
- [ ] Examples or fixtures added for changed contract.
- [ ] Validation/tests planned or added for changed contract.
- [ ] Cross-module data is typed or documented.
- [ ] No undocumented JSON shapes, untyped maps, or silent schema drift for core boundaries.

## Storage separation

- [ ] Save DB and cache DB are separate.
- [ ] Cache clean/rebuild cannot delete save data.
- [ ] Save migrations are explicit and backup-aware.
- [ ] Runtime opens cache read-only where possible.
