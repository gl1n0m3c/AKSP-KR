-- +goose Up
-- +goose StatementBegin
CREATE TABLE positions
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL DEFAULT ''
);

CREATE TABLE teams
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL DEFAULT ''
);

CREATE TABLE units
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL DEFAULT ''
);

CREATE TABLE users
(
    id          UUID PRIMARY KEY,
    first_name  VARCHAR(50)                    NOT NULL,
    last_name   VARCHAR(50)                    NOT NULL,
    patronymic  VARCHAR(50)                    NOT NULL,
    email       VARCHAR(50)                    NOT NULL,
    phone       VARCHAR(50),
    telegram    VARCHAR(50),
    status      VARCHAR(16)                    NOT NULL CHECK (status IN ('active', 'inactive', 'deleting', 'deleted')),
    position_id UUID REFERENCES positions (id) NOT NULL,
    team_id     UUID REFERENCES teams (id)     NOT NULL,
    unit_id     UUID REFERENCES units (id)     NOT NULL,
    head_id     UUID REFERENCES users (id),
    created_at  TIMESTAMPTZ                    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    created_by  UUID REFERENCES users (id)     NOT NULL DEFAULT '30f6afc4-b8eb-4e6b-b0bd-158fcef0dc28',
    updated_by  UUID REFERENCES users (id)
);

CREATE TABLE meetings
(
    id           UUID PRIMARY KEY,
    title        VARCHAR(50)  NOT NULL,
    description  VARCHAR(300) NOT NULL DEFAULT '',
    scheduled_at TIMESTAMPTZ  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ,
    organizer_id UUID REFERENCES users (id),
    team_id      UUID REFERENCES teams (id)
);

INSERT INTO positions (id, name, description)
VALUES ('65b618f2-a563-45aa-8dd8-5976c03a2059', 'Admin', 'Built-in root admin')
ON CONFLICT DO NOTHING;

INSERT INTO teams (id, name, description)
VALUES ('6109edb9-48f8-49fa-b461-b12379459eb1', 'Admin', 'Built-in root team')
ON CONFLICT DO NOTHING;

INSERT INTO units (id, name, description)
VALUES ('e5366c33-df2e-4e9d-9100-e6536786a570', 'Admin', 'Built-in root unit')
ON CONFLICT DO NOTHING;

INSERT INTO users (
    id, first_name, last_name, patronymic,
    email, phone, telegram, status,
    position_id, team_id, unit_id,
    created_at, updated_at
) VALUES (
    '30f6afc4-b8eb-4e6b-b0bd-158fcef0dc28',
    'Admin', 'Admin', 'Admin',
    'admin@example.com', NULL, NULL, 'active',
    '65b618f2-a563-45aa-8dd8-5976c03a2059',
    '6109edb9-48f8-49fa-b461-b12379459eb1',
    'e5366c33-df2e-4e9d-9100-e6536786a570',
    NOW(), NULL
) ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS units;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS positions;
-- +goose StatementEnd
