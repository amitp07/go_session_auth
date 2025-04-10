alter table users 
add column email VARCHAR(250) not null UNIQUE,
add column mfa_enabled boolean not null DEFAULT false,
add column enabled boolean not null DEFAULT true,
add column email_verified boolean not null DEFAULT false;