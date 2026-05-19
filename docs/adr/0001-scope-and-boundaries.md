# ADR 0001: Scope and Boundaries

## Status

Accepted

## Context

Phanes is intended to become a local-first offline runtime, resource-cache builder, and launcher framework for sandbox research. It should feel lightweight and local/offline, while still allowing a local runtime service where that keeps the architecture simpler.

Historical projects such as Grasscutter and Cultivation are useful references for concepts like dispatch, process supervision, and local account/session flow. They are not implementation bases for Phanes.

## Decision

Phanes will be designed as a new project with these mandatory boundaries:

- Runtime is localhost-only by default and binds to literal `127.0.0.1`.
- `0.0.0.0` and public server hosting are out of scope.
- Runtime startup reads local save data and local resource cache only.
- Runtime must not download or build resources during startup.
- GC-Resources is reference-only and must not be required.
- Full copyrighted game resources must not be bundled or redistributed.
- Local cache generation uses user-provided local inputs and normalized local output.
- No client patching, deep patching guides, anti-cheat bypass, protection bypass, official authentication bypass, unauthorized online interaction, or commercial service replacement.
- External projects may inform concepts, but source code must not be copied without explicit license review and maintainer approval.

## Consequences

- The runtime can fail fast when a cache is missing or invalid, with diagnostics that tell the user to run the builder first.
- Builder and runtime remain separate responsibilities.
- Launcher supervises local processes and diagnostics; it does not mutate system settings or patch clients.
- Network-first, public hosting, and bypass-oriented features require an explicit scope change by the maintainer before discussion or implementation.

## Compliance checks

- Initial config validation accepts only literal `127.0.0.1` as a bind host. `0.0.0.0`, public/non-loopback addresses, `localhost`, `::1`, and other aliases require a documented validation update before use.
- Runtime contracts contain no remote resource URL fields.
- Builder input contracts accept local paths, not remote repositories, as the default path.
- Cache deletion commands cannot touch save data.
- Reviews check for copied code from reference projects.

## Related

- [Module boundaries](../contracts/module-boundaries)
- [Compliance checklist](../contracts/compliance-checklist)
