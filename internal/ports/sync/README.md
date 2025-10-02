# ports/sync (blob transport)

Purpose
- Sync file contents (blobs) to/from a remote (e.g., rsync over SSH).

Contains
- Interface `Transport`, types `SSHOptions`, `Report`

Does NOT contain
- Shelling out or SSH specifics (that’s for adapters)
