CREATE TABLE post (
		id VARCHAR(40) PRIMARY KEY NOT NULL,
		title TEXT,
		content TEXT,
		author varchar(250),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
INSERT INTO post (id, title, content, author, created_at, updated_at) VALUES
    ('1', 'First Post', 'This is the content of the first post.', 'Author1', NOW(), NOW()),
    ('2', 'Second Post', 'This is the content of the second post.', 'Author2', NOW(), NOW()),
    ('3', 'Third Post', 'This is the content of the third post.', 'Author3', NOW(), NOW());

CREATE OR REPLACE FUNCTION notify_trigger()
RETURNS TRIGGER AS $$
DECLARE
	notification json;
BEGIN
	IF (TG_OP = 'INSERT') THEN
		-- notification := row_to_json(NEW);
		notification := json_build_object(
			'id', NEW.id,
			'title', NEW.title,
			'content', NEW.content,
			'author', NEW.author,
			'updated_at', to_char(NEW.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
			'created_at', to_char(NEW.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')
		);
	ELSIF (TG_OP = 'DELETE') THEN
		-- notification := row_to_json(OLD);
		notification := json_build_object(
			'id', OLD.id,
			'deleted', true
		);
	END IF;
    PERFORM pg_notify('db_changes', notification::text);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER row_change
AFTER INSERT OR DELETE ON post
FOR EACH ROW
EXECUTE FUNCTION notify_trigger();
