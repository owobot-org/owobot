# Installing owobot

## Discord setup

Before installing the bot, you need to do some set up in your Discord account. Start by opening the [Developer Portal](https://discord.com/developers/applications) and creating a new application. Once you've created it, click on your new application. Then, go to `Bot` in the sidebar and enable the Presence Intent, Server Members Intent, and Message Content intent.

To get the token that you have to pass to `owobot`, click on the `Reset Token` button to generate a new token, and then copy it and store it somewhere safe (you won't be able to see it later).

## Linux packages

The [latest release](https://gitea.elara.ws/owobot/owobot/releases/latest) contains RPM, Deb, and Arch packages, and owobot is available [on the AUR](https://aur.archlinux.org/packages/owobot-bin/) as well. Choose whichever one of those you need and install it with your package manager.

Once it's installed, there should be a default config file at `/etc/owobot.toml`. You can edit that to add your token, change the activity text, etc. and then run the bot by running `sudo systemctl enable --now owobot`. Systemd will now start running the bot and monitoring it to make sure it doesn't go down.

That's it! Your bot should be up and running!

## Docker

This guide will use Docker, but `owobot` should work with any other OCI-compatible container engine, such as Podman. The container image is hosted on [Gitea](https://gitea.elara.ws/owobot/-/packages/container/owobot/latest).

There's a [`docker-compose.yml`](docker-compose.yml) file provided in this repo as a starting point. Here's how you can use it:

1. First, make sure `docker` and `docker-compose` are installed and working
2. Create a new folder for owobot to use to store its data
3. Put the example `docker-compose.yml` file into the new folder
4. Edit the `docker-compose.yml` file to set the token and anything else you may want to change
5. Make sure the directory can be accessed by the container's user (`sudo chown -R 65532:65532 folder`)
6. Run `docker-compose up -d`
7. That's it! Your bot should now be running.