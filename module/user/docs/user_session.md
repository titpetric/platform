# User Session

Stores immutable user sessions with ULID IDs and expiration.

| Name       | Type      | Key | Comment                                              |
|------------|-----------|-----|------------------------------------------------------|
| id         | char(26)  | PRI | Primary key: ULID string, also used as session token |
| user_id    | char(26)  | MUL | Reference to user.id (ULID)                          |
| expires_at | timestamp |     | Session expiration timestamp                         |
| created_at | timestamp |     | Session creation timestamp                           |
