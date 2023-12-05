/* version holds the current database version for migrations */
CREATE TABLE version (
  current TEXT NOT NULL
);

/* guilds holds all the guilds this bot is in and their settings */
CREATE TABLE guilds (
  id                   TEXT PRIMARY KEY,
  starboard_chan_id    TEXT NOT NULL DEFAULT '',
  starboard_stars      INT  NOT NULL DEFAULT 3,
  log_chan_id          TEXT NOT NULL DEFAULT '',
  ticket_log_chan_id   TEXT NOT NULL DEFAULT '',
  ticket_category_id   TEXT NOT NULL DEFAULT '',
  vetting_req_chan_id  TEXT NOT NULL DEFAULT '',
  vetting_role_id      TEXT NOT NULL DEFAULT ''
);

/* tickets holds all open tickets */
CREATE TABLE tickets (
  user_id    TEXT PRIMARY KEY,
  guild_id   TEXT NOT NULL DEFAULT '',
  channel_id TEXT UNIQUE,
  UNIQUE(user_id, guild_id)
);

/* vetting_requests holds all open vetting requests */
CREATE TABLE vetting_requests (
  user_id    TEXT PRIMARY KEY,
  guild_id   TEXT NOT NULL DEFAULT '',
  msg_id     TEXT UNIQUE,
  UNIQUE(user_id, guild_id)
);

/* starboard holds all the messages that have been added to the starboard */
CREATE TABLE starboard (
  id TEXT PRIMARY KEY
);

/* polls holds all the poll messages created by the bot and their state */
CREATE TABLE polls (
  msg_id     TEXT PRIMARY KEY,
  owner_id   TEXT NOT NULL,
  title      TEXT NOT NULL,
  finished   BOOL NOT NULL DEFAULT false,
  opt_emojis TEXT NOT NULL DEFAULT '',
  opt_text   TEXT NOT NULL DEFAULT ''
);

/* votes holds all the poll votes */
CREATE TABLE votes (
  poll_msg_id TEXT NOT NULL,
  user_token  TEXT NOT NULL,
  option      INT  NOT NULL,
  UNIQUE(poll_msg_id, user_token) ON CONFLICT REPLACE
);

CREATE INDEX idx_votes_option ON votes(option);
CREATE INDEX idx_votes_poll_msg_id ON votes(poll_msg_id);

/* reactions holds all the message reactions */
CREATE TABLE reactions (
  guild_id      TEXT NOT NULL,
  match_type    TEXT NOT NULL,
  match         TEXT NOT NULL,
  reaction_type TEXT NOT NULL,
  reaction      TEXT NOT NULL,
  chance        INT  NOT NULL CHECK (chance >= 1 AND chance <= 100) DEFAULT 100
);

CREATE INDEX idx_reactions_guild_id ON reactions(guild_id);
CREATE INDEX idx_reactions_match ON reactions(match);

/* reaction_role_categories holds the reaction roles */
CREATE TABLE reaction_role_categories (
  msg_id      TEXT PRIMARY KEY,
  channel_id  TEXT NOT NULL,
  name        TEXT NOT NULL,
  description TEXT NOT NULL,
  emoji       TEXT NOT NULL,
  roles       TEXT NOT NULL,
  UNIQUE(channel_id, name)
);