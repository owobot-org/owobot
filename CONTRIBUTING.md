# Contributing to owobot

Thanks for your interest in contributing to owobot! This page contains information that you should know before contributing.

## Code structure

### Systems

owobot consists of several independent systems, such as the `starboard` system, `members` system, `commands` system, etc. These systems are what actually interact with users and they're all in the `internal/systems` directory.

All the systems that require initialization have an `Init(*discordgo.Session) error` function, which does things like registers all the commands and handlers, and performs any other initialization steps that need to be done for that system. These `Init` functions are called by `main.go` when the bot starts up.

The `commands` system always starts last because the other systems register commands that it needs to know about before it does its initialization.

### Database

All the database code is in `internal/db`. owobot doesn't use any ORM or framework for the database, it directly executes SQL queries. Database migrations are stored in `internal/db/migrations`. They are sql files whose names contain the date when they were made and an extra number to avoid collisions in case multiple migrations are ever made in the same day.

If you change anything in the database, always make a new migration file rather than editing existing ones. This way, owobot will automatically apply the the changes whenever it's run next. Changing migrations requires a full recompile because they're embedded into the binary.

## Testing

If you want to test out your changes, you'll need to make a test server and bot account. To do that, go to https://discord.com/developers/applications and create a new application. Then, go to `Bot` in the sidebar, and enable the privileged gateway intents for `Message Content` and `Server Members`. Now, go to `OAuth2 > URL Generator`, select `bot` in Scopes, and then `Administrator` in Bot Permissions. That will give you a URL. Next, go to Discord and make a new server that you'll use for testing. Then, paste the URL you generated into your browser and invite your test bot into your new server.