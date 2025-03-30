create table user_groups_roles (
    user_group_id uuid references user_groups(id) on delete cascade,
    role_id uuid references roles(id) on delete cascade,
    primary key(user_group_id, role_id)
)