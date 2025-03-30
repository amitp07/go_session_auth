CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(20) UNIQUE NOT NULL CHECK(name IN ('admin', 'super_user', 'contributor', 'reader')),
    description varchar(250),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
)