# Knights Of Discord

A Discord bot used to bring life to your servers through an RPG-styled in-chat game. **Knights Of Discord is still work in progress and is currently unplayable!**

## Adding The Bot To Your Server

**Knights Of Discord is still work in progress and is currently unplayable!**

Note that adding the bot requires allowing the following permissions:
1. **Manage Channels** - Allows management and editing of channels.
2. **Add Reactions** - Allows for the addition of reactions to messages.
3. **View Channel** - Allows guild members to view a channel, which includes reading messages in text channels.
4. **Send Messages** - Allows for sending messages in a channel.
5. **Manage Messages** - Allows for deletion of other users messages.
6. **Mention Everyone** - Allows for using the `@everyone` tag to notify all users in a channel, and the `@here` tag to notify all online users in a channel.
7. **Use External Emoji** - Allows the usage of custom emojis from other servers.
8. **Manage Roles** - Allows management and editing of roles.


[Click here to add Knights Of Discord to your server/guild](https://discordapp.com/oauth2/authorize?client_id=487744442531315712&scope=bot&permissions=268840016)

## Current State

The project is work in progress.

The current step if building and maintaining the game environment in all servers the Bot is present.

## Contributing
1. Fork this repository
2. Make code changes
3. Compile and test in your own test server
4. Submit a PR to this repository

To get the bot running locally:
1. Create a directory for GoLang projects (example: `/go`)
2. Create an `src` directory inside (example: `/go/src`)
3. Create the path to this repository ot your fork (example: `/go/src/github.com/Noxdew`)
4. Clone the repository (example command: `git clone git@github.com:Noxdew/Knights-Of-Discord.git`)
5. Add your test bot authentication token to the settings.
6. Run the bot (example: `go run main.go` or check the `makefile`)
