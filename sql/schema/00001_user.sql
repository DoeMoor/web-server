-- +goose Up
-- +goose StatementBegin
create table users (
  id UUID not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  email text not null,
  primary key (id)
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table users;

-- +goose StatementEnd
