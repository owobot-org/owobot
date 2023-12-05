# owobot - The coolest bot ever written

[![owobot-bin AUR package](https://img.shields.io/aur/version/owobot-bin?label=owobot-bin&logo=archlinux)](https://aur.archlinux.org/packages/owobot-bin/)
[![OCI Image Badge](https://img.shields.io/badge/oci-images-24184C?logo=opencontainersinitiative)](https://gitea.elara.ws/owobot/-/packages/container/owobot/latest)

## Introduction

owobot is a powerful Discord bot designed to handle a wide range of tasks in your server, from moderation to entertainment. It takes advantage of several cutting-edge discord features, such as message buttons and modals.

## Features

### Vetting

In order to catch trolls and other troublemakers before they get access to your server, owobot can be configured to require new users to go through a vetting process before gaining access to the server.

To create your vetting message, just choose any message, click `More > Apps > Make Vetting Message`, and that's it! owobot will delete the message and post a new one with a message button which can be used by users to request vetting.

When users click the request vetting button, owobot will send a vetting request in the vetting request channel.

If a moderator accepts the request, a new ticket will be created in which mods can talk to the user. When they're finished, they can either kick the user which will automatically close the ticket, or they can approve the user using the `/approve` command.

Commands:

- `/vetting role` can be used by anyone with the `Manage Server` permission to set the server's vetting role. owobot will assign this role to all new users.
- `/vetting req_channel` can be used by anyone with the `Manage Server` permission to set the server's vetting request channel. This is where owobot will post vetting requests.
- `/approve` can be used by anyone with the `Kick Members` permission to approve users that are in vetting.

### Tickets

owobot can create tickets to allow users to privately talk to your server's moderators. Only one ticket per user can exist at any time. When a ticket is closed, a log containing all the messages in the ticket is sent to the event log ticket channel.

Commands:

- `/ticket` can be used by any user to create a ticket for themselves
- `/mod_ticket` can be used by anyone with the `Manage Channels` permission to create a ticket for another user
- `/close_ticket` can be used by anyone with the `Manage Channels` permission to close a user's ticket
- `/ticket_category` can be used by anyone with the `Manage Server` permission to set the category in which ticket channels are created.

### Eventlog

The eventlog sends important events such as kicks/bans, role changes, etc. to a configurable discord channel.

Commands:

- `/eventlog channel` can be used by anyone with the `Manage Server` permission to set the channel for the event log
- `/eventlog ticket_channel` can be used by anyone with the `Manage Server` permission to set the channel in which ticket conversations logs will be sent

### Reactions

owobot has a very powerful reaction system which can find content inside of messages and then react with an emoji or reply with text.

A single reaction consists of a match type, match, reaction type, reaction, and an optional random chance.

The match type can either be `contains` or `regex`. The `contains` matcher checks if a message contains the match. The `regex` matcher checks if a message matches a regular expression and extracts any submatches. If you're using the `regex` matcher with the `text` reaction type, you can include submatches in your reply by putting the submatch index in curly braces (for example: `{1}` or `{5}`).

The optional random chance allows you to add reactions that only occur a certain percentage of the time. Setting it to `10`, for example, means the reaction will only happen in 10% of detected messages.

Commands:

- `/reactions add` can be used by anyone with the `Manage Expressions` permission to add new reactions
- `/reactions list` can be used by anyone with the `Manage Expressions` permission to get a list of all existing reactions
- `/reactions delete` can be used by anyone with the `Manage Expressions` permission to delete an existing reaction

> [!NOTE]
> The `/reactions delete` command has an extra check to make sure the user has `Manage Expressions`, so you can safely add a role override on the `/reactions` command without allowing anyone without that permission to delete reactions.

### Reaction Roles

Reaction roles allow users to easily assign roles to themselves using message buttons.

Reaction roles are organized in categories, which can have a name and a description. You can't have more than one category with the same name in a given channel.

Commands:

- `/reaction_roles new_category` can be used by anyone with the `Manage Server` permission to create a new reaction role category in the current channel
- `/reaction_roles remove_category` can be used by anyone with the `Manage Server` permission to remove an existing reaction role category from the current channel
- `/reaction_roles add` can be used by anyone with the `Manage Server` permission to add a reaction role to a category
- `/reaction_roles remove` can be used by anyone with the `Manage Server` permission to remove an existing reaction role from a category

> [!TIP]
> There's a `/neopronoun` command that any user can use to assign themselves a pronoun role. It will never assign a role that provides any permissions, so it's safe to allow for everyone.

### Polls

owobot can easily create polls for your members to vote in. Polls use message components and privacy tokens to ensure that votes are always private and even the person running the bot can't find out who voted for what.

A poll can be created using the `/poll` command. owobot will create a message with just the title and two buttons: `Add Options` and `Finish`. Clicking the `Add Options` button opens a modal (pop up) where you can type the text for a new option. Once that's done, owobot edits the message and asks the poll owner to react with the emoji they'd like to use for that poll. Once they react, that option is added. Options can keep being added until the Finish button is clicked, which finalizes the poll, creates a thread, and opens it up to votes.

Commands:

- `/poll` can be used by any user to create a poll

### Starboard

The starboard is a way for your users to feature the messages they like. Users can react to messages with stars, and once a configurable threshold of stars is reached, the message will be posted to the starboard channel.

Commands:

- `/starboard stars` can be used by anyone with the `Manage Server` permission to set the star reaction threshold for the starboard (The default is 3)
- `/starboard channel` can be used by anyone with the `Manage Server` permission to set the starboard channel for the server.

### Rate Limiting

owobot will rate limit events such as channel deletions, kicks, and bans, ensuring that compromised mod accounts can't destroy the server. If a user gets near the rate limit, they'll receive two warnings and then they'll be kicked from the server.

Here are the current rate limits:

- `channel_delete`: 10 / minute
- `kick`: 10 / minute
- `ban`: 7 / 5 minutes