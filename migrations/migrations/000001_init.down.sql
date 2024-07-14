BEGIN;

--DROP EXTENSION IF EXISTS "uuid-ossp";

DROP INDEX IF EXISTS idx__sessions__identifier;
DROP INDEX IF EXISTS idx__sessions__user_id;
DROP INDEX IF EXISTS idx__users__email;

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;

COMMIT;
