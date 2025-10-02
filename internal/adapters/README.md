# adapters (implementations of ports)

Purpose
- Provide concrete implementations for the port interfaces (filesystems, sync, locks, params).

Examples (future)
- `localfs`   → implements vaultfs.FS using the local OS filesystem (or your existing logic)
- `manager`   → implements vaults.Manager (in-memory or config-backed)
- `memorylocks` → implements locks.Service in-memory (dev/testing)
- `rsync`     → implements sync.Transport via rsync over SSH
- `db`  → implements params.Store and/or locks.Service against your DB

Rules
- Adapters import `internal/ports/*` and `internal/domain/models`, plus any tech libs they need.
- Adapters do NOT import GUI code, and ideally do not import other adapters.
