package token

const getUserByLoginDML = `SELECT id, login, password, role
FROM users
WHERE login = $1;`