-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_participants (
    meeting_id UUID REFERENCES meetings(id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (meeting_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meeting_participants;
-- +goose StatementEnd

