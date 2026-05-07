CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role AS ENUM (
    'participant',
    'judge',
    'admin'
);

CREATE TABLE IF NOT EXISTS accounts(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(128) UNIQUE NOT NULL,
    password varchar(128) NOT NULL,
    role role NOT NULL DEFAULT 'participant'::role,
    name varchar(128) NOT NULL,
    school varchar(64),
    student_id varchar(32),
    phone_number varchar(16),
    is_verified bool NOT NULL DEFAULT false,
    is_banned bool NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
