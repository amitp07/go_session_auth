create table user_group_members (
    user_id uuid references users(id) on delete cascade,
    user_group_id uuid references user_groups(id) on delete cascade,
    primary key(user_id, user_group_id)
)