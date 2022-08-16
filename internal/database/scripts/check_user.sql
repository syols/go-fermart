SELECT id
FROM users
WHERE login = :login
  AND password = :password;