package config

const AddNewUser = `insert into "users"("user_id", "user_name", "user_tgid", "start_time") values($1, $2, $3, $4)`
const UserDB = `select * from users where user_id = $1`
const UsersFromDB = `select * from users`
