# User Group Member

Stores user memberships in groups using ULID for IDs.

| Name          | Type      | Key | Comment                              |
|---------------|-----------|-----|--------------------------------------|
| user_group_id | char(26)  | PRI | Reference to user_group.id (ULID)    |
| user_id       | char(26)  | PRI | Reference to user.id (ULID)          |
| joined_at     | timestamp |     | Timestamp when user joined the group |
