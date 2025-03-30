create table user_group_members (
    user_id uuid references users(id) on delete cascade,
    group_id uuid references roles(id) on delete cascade,
    primary key(user_id, group_id)
)