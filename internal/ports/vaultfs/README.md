# ports/vaultfs (single-vault file system interface)

Purpose
- Define the interface to browse and mutate entries within one vault.

Contains
- Interface `FS` (Info, List, Stat, Rename/Move/Copy/Delete)

Does NOT contain
- OS/DB/Network access
- Concrete implementations
