CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role AS ENUM (
    'participant',
    'jury',
    'admin'
);

CREATE TABLE IF NOT EXISTS accounts(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    email text UNIQUE NOT NULL,
    name text NOT NULL,
    password varchar(128) NOT NULL,
    role role NOT NULL DEFAULT 'participant'::role,
    is_banned bool NOT NULL DEFAULT false,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TYPE language AS ENUM (
    'rust',
    'cpp',
    'python',
    'java',
    'html'
);

CREATE TABLE IF NOT EXISTS contests(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    name text NOT NULL,
    start_time timestamptz NOT NULL,
    end_time timestamptz NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS problems(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    contest_id uuid REFERENCES contests(id),
    position int,

    name text NOT NULL,
    content_path text NOT NULL,
    checker_language language,
    checker_path text,
    time_limit int,
    memory_limit int,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS test_cases(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    problem_id uuid NOT NULL REFERENCES problems(id),

    input_path text NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS problem_languages(
    problem_id uuid NOT NULL REFERENCES problems(id),
    language language NOT NULL,

    PRIMARY KEY (problem_id, language)
);

CREATE TYPE verdict AS ENUM (
    'accepted',
    'wrong_answer',
    'time_limit_exceeded',
    'compilation_error',
    'memory_limit_exceeded',
    'runtime_error',
    'idle_time_limit_exceeded'
);

CREATE TABLE IF NOT EXISTS submissions(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    problem_id uuid NOT NULL REFERENCES problems(id),
    account_id uuid NOT NULL REFERENCES accounts(id),

    language language NOT NULL,
    code text NOT NULL,
    verdict verdict,
    score real,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
