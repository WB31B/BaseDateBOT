package config

const ADDNEWUSER = `insert into "users"("user_id", "user_name", "user_tgid", "start_time") values($1, $2, $3, $4)`
const USERDB = `select * from users where user_id = $1`
const USERSDB = `select * from users`
