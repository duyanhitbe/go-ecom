CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   username VARCHAR NOT NULL,
   password VARCHAR NOT NULL,
   created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
   updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
