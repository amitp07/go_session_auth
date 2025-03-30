create table permissions (
    id uuid primary key default gen_random_uuid(),
    name varchar(20) unique not null check (name in ('read', 'write', 'execute')),
    description text
)