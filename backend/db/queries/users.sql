-- name: GetUserByID :one
select id, auth_provider_id, handle, display_name, created_at, updated_at
from users
where id = $1;

-- name: CreateUser :one
insert into users (id, auth_provider_id, handle, display_name)
values ($1, $2, $3, $4)
returning id, auth_provider_id, handle, display_name, created_at, updated_at;

-- name: CountUsers :one
select count(*) from users;
