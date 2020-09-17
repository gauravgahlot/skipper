SET ROLE tinkerbell;

CREATE TABLE IF NOT EXISTS events (
	id UUID UNIQUE NOT NULL primary key
	, resource_id UUID NOT NULL
	, resource_type int NOT NULL
	, event_type int NOT NULL
	, created_at TIMESTAMPTZ
	, data JSONB
);

CREATE INDEX IF NOT EXISTS idx_eid ON events (id);
CREATE INDEX IF NOT EXISTS idx_etype ON events (event_type);
CREATE INDEX IF NOT EXISTS idx_rid ON events (resource_id);
CREATE INDEX IF NOT EXISTS idx_rtype ON events (resource_type);

CREATE OR REPLACE FUNCTION notify_event_changes()
RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify('events_channel', row_to_json(NEW)::text);
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_channel
AFTER INSERT ON events
FOR EACH ROW EXECUTE PROCEDURE notify_event_changes()
