-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
    name, age, 	Password, Address, Country_code, Email, Phone, Created_at, Updated_at
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
    RETURNING *;

-- name: CreateUserProfiles :one
INSERT INTO user_profiles (
    user_id, company_name, country, state
) VALUES (
             $1, $2, $3, $4
         )
RETURNING *;



-- name: CreateUserRoles :one
INSERT INTO user_roles (
    user_id, role
) VALUES (
             $1, $2
         )
RETURNING *;





-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;