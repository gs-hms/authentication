CREATE TABLE authentication_sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
)