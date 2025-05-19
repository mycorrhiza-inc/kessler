-- +goose Up
ALTER TYPE stage_state
ADD
    VALUE 'errored';

-- +goose Down
-- ALTER TYPE stage_state DROP VALUE 'errored' CASCADE;
