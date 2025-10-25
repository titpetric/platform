# User Auth

Stores user authentication credentials.

| Name       | Type         | Key | Comment                     |
|------------|--------------|-----|-----------------------------|
| user_id    | char(26)     | PRI | Reference to user.id (ULID) |
| email      | varchar(255) | UNI | User email address, unique  |
| password   | varchar(255) |     | Hashed password             |
| created_at | timestamp    |     | Record creation timestamp   |
| updated_at | timestamp    |     | Record update timestamp     |
