# domain/models (pure domain types)

Purpose
- Define FreePDM’s domain data structures (the ubiquitous language).
- No I/O, no business logic, no external dependencies beyond Go stdlib.

Contains
- `VaultInfo`, `Entry`, `Lock`, `Version`, `Parameters`, `UserIdentity`
- Query helpers like `SearchQuery`, `Range`, `SearchResult`

Does NOT contain
- File/DB/network access
- GUI code
- Implementation details of any adapter
