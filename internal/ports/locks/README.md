# ports/locks (locking service)

Purpose
- Define pessimistic locking (checkout/checkin/status).

Contains
- Interface `Service` (Status, Checkout, Heartbeat, Checkin, ForceUnlock)

Does NOT contain
- Lock storage/DB/ACL mechanics
