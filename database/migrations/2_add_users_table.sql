-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id uuid PRIMARY KEY default gen_random_uuid(),
  email varchar(255) UNIQUE NOT NULL,
  name varchar(255) NOT NULL DEFAULT 'default_name',
  last_name varchar(255) NOT NULL DEFAULT 'default_last_name',
  password varchar(255) NOT NULL DEFAULT 'default_password',
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create trigger update_users_update_at
before update on users for each row execute procedure update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
