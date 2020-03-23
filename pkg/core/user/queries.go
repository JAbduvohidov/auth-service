package user

const addUserDML = `INSERT INTO users (name, surname, login, password, avatar)
VALUES ($1, $2, $3, $4, $5);`
