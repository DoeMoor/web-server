-- +goose Up
-- +goose StatementBegin
alter table chirps alter column user_id set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
