# ports/vaults — Multi-vault Manager (Interfaces Only)

**Purpose**  
Define a small, stable interface to:
- Discover and access multiple vaults (`List`, `Get`, `Add`, `Remove`)
- Perform **cross-vault** operations (`Copy`, `Move`)

This isolates the GUI from how vaults are discovered/stored (local config, DB, remote API) and centralizes policies that span more than one vault (e.g., read-only enforcement on destination vaults).

---

## What belongs here

- `Manager` interface
- Imports only:
  - `internal/domain/models` (types)
  - `internal/ports/vaultfs` (single-vault FS interface)

**No implementations here.** Concrete code lives under `internal/adapters/...`.

---

## Why a Manager?

- **Single entry point** for the GUI to discover vaults and obtain a `vaultfs.FS`.
- **Cross-vault operations** live in one place (don’t leak other vaults into a single `VaultFS`).
- **Policy & validation** (read-only, user permissions, naming rules) can be consistently enforced before delegating to the underlying vault file systems.
- **Wiring flexibility**: today local adapters, tomorrow remote/rsync/DB—GUI code unchanged.

