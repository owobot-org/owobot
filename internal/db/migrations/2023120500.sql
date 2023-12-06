/* add time format column */
ALTER TABLE guilds ADD COLUMN time_format TEXT NOT NULL DEFAULT 'discord';

/* add welcome message columns */
ALTER TABLE guilds ADD COLUMN welcome_chan_id TEXT NOT NULL DEFAULT '';
ALTER TABLE guilds ADD COLUMN welcome_msg TEXT NOT NULL DEFAULT '';