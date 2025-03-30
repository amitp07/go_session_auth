create table users_roles (
    user_id uuid references users(id) on delete cascade,
    role_id uuid references roles(id) on delete cascade,
    primary key(user_id, role_id)
)