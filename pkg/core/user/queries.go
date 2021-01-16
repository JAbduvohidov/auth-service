package user

const addUserDML = `INSERT INTO users (name, surname, login, password, avatar)
VALUES ($1, $2, $3, $4, $5);`

const getUserDML = `SELECT id, name, surname, login, avatar FROM users WHERE id = $1;`

const updateUserProfileDML = `UPDATE users
SET name    = $1,
    surname = $2
WHERE id = $3;`

const updateUserPasswordDML = `UPDATE users
SET password = $1
WHERE id = $2;`

const deleteUserDML = `UPDATE users
SET removed = 'true'
WHERE id = $1;`

const getUsersDML = `SELECT id, name, surname, login, role, avatar FROM users WHERE role <> 'MODERATOR' AND removed = false;`

const upgradeUserDML = `UPDATE users
SET role = $1
WHERE id = $2;`