# User

Stores user profile information using ULID as primary key.

| Name       | Type         | Key | Comment                               |
|------------|--------------|-----|---------------------------------------|
| id         | char(26)     | PRI | Primary key: ULID string              |
| first_name | varchar(100) |     | User first name                       |
| last_name  | varchar(100) |     | User last name                        |
| deleted_at | timestamp    | MUL | Soft delete timestamp, NULL if active |
| created_at | timestamp    |     | Record creation timestamp             |
| updated_at | timestamp    |     | Record update timestamp               |
