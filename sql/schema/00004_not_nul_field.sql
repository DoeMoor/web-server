-- +goose Up
-- +goose StatementBegin
alter table users
    add constraint user_email_pk
        unique (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    drop constraint user_email_pk;
-- +goose StatementEnd
