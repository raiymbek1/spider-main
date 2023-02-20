CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table posts(id uuid primary key NOT NULL default uuid_generate_v4(), url text, title text, body text, created_at date);