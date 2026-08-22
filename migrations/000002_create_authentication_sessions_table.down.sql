ALTER TABLE authentication_sessions
DROP CONSTRAINT IF EXISTS authentication_sessions_user_id_fkey;

DROP TABLE IF EXISTS authentication_sessions;