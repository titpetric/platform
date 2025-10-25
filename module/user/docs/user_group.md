# User Group

Stores user group information using ULID as primary key.

| Name       | Type         | Key | Comment                   |
|------------|--------------|-----|---------------------------|
| id         | char(26)     | PRI | Primary key: ULID string  |
| title      | varchar(100) | UNI | Group name/title          |
| created_at | timestamp    |     | Record creation timestamp |
| updated_at | timestamp    |     | Record update timestamp   |
