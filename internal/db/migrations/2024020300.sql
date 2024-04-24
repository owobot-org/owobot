/* plugins stores information about all the plugins defined for this bot. */
/* This will be used to let plugins perform actions when they're updated  */
CREATE TABLE plugins (
	name        TEXT NOT NULL,
	version     TEXT NOT NULL,
	description TEXT NOT NULL,
	UNIQUE(name) ON CONFLICT REPLACE
);

/* Add a column to allow guilds to enable whichever plugins they want */
ALTER TABLE guilds ADD COLUMN enabled_plugins TEXT NOT NULL DEFAULT '';
