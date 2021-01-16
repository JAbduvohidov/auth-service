package token

const getUserByLoginDML = `SELECT id, name, login, password, role
FROM users
WHERE login = $1;`