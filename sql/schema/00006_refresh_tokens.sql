-- +goose Up
-- +goose StatementBegin
create table refresh_tokens(
  token text not null primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  user_id UUID references users (id) on delete cascade not null,
  expires_at timestamp not null,
  revoked_at timestamp default null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table refresh_tokens;
-- +goose StatementEnd
