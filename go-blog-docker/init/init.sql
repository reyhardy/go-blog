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
BEGIN
    PERFORM pg_notify('db_changes', row_to_json(NEW)::text);
	-- PERFORM pg_notify('db_changes', NEW);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER row_change
AFTER INSERT ON post
FOR EACH ROW
EXECUTE FUNCTION notify_trigger();
