---
slug: /
---

# Phanes Documentation

Phanes is a local-first offline runtime, resource-cache builder, and launcher framework for sandbox research.

This documentation site is the canonical browsable view of the contract-first design docs in this repository.

## Start here

- [Scope and boundaries](adr/scope-and-boundaries)
- [Module boundaries](contracts/module-boundaries)
- [Protobuf contracts](contracts/protobuf)
- [Compliance checklist](contracts/compliance-checklist)

## Hard boundaries

- Runtime defaults to literal `127.0.0.1`.
- Runtime does not download resources.
- GC-Resources is reference-only and not required.
- No client patching, anti-cheat bypass, protection bypass, or official service authentication bypass.
- Save data and resource cache remain separate.
