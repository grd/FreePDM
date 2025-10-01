# apps/fpg (GUI)

Purpose
- Fyne GUI for FreePDM that depends only on ports (interfaces) and models (types).

Examples
- Choose `localfs` for vaultfs.FS, `memorylocks` for locks, `paramsstub` for params, `rsync` (dry-run) for sync.
- Read env/flags, create the adapters, pass them into GUI constructors.

Rules
- No direct filesystem/db/ssh code.
- Receive concrete adapters via constructors (passed in from cmd/fpg composition).
- Log any external calls via the port interfaces (observability).
- Keep this thin and declarative.
- No business logic; no GUI code beyond wiring.