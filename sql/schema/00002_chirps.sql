-- +goose Up
-- +goose StatementBegin
CREATE TABLE chirps (
  id UUID NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  body text NOT NULL,
  user_id UUID references users (id) on delete cascade,
  PRIMARY KEY (id)
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table chirps;

-- +goose StatementEnd
