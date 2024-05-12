# ON Esports Bot for Discord

Welcome to the implementation of the Discord Bot for [ONEsports.gg](https://onesports.gg/). Along the way, I implemented a Clean Code Architecture and made project structure separated into modules to provide future extendibility and flexibility.

To get started with using this bot, please follow the setup instructions below.

## So, what this bot can do?
This bot is set up for many different use cases that can be added on demand. However, for now, bot's main funcitonality is:

- Joining a newly created voice channel, recording audio for specified amount of time and saving audio file in .ogg format to AWS S3 and retaining reference in PostgreSQL.
- Sending the recorded audio file to `# voice-chat-recordings` text channel.

## How can I try using this bot?

Follow these steps to set up the ON Esports Bot in your Discord server:

**1. Create a Discord Server**  
First of all, you will need a server in Discord where you are going to add the bot.

**2. Create a Discord Application**  
Next, you need to navigate to [Discord Developer Portal](https://discord.com/developers/applications) and do following steps:

-   Create a new application.
-   Go to OAuth2, add a bot and give him Administrator permissions.
-   Save Application Token (we will need it later).

**3. Add bot to your server**

-   Navigate to the generated link.
-   Login and choose the server we created in Step 1.

**4. Clone the Repository**  
 Start by cloning the bot repository to your local machine. Use the following command:

```bash
git clone https://github.com/nrakhay/ONEsports.git
cd ONEsports
cp .env-example .env # change *cp* to *copy* if you are using Windows
```

Then, go to .env file and put token from Step 2 instead of {YOUR_TOKEN_HERE}. 

Also, do not forget to populate AWS variables and put ID of a Text Channel where bot should send audio recordings. 

**5. Setup Database**  

```bash
cd local/
docker compose up -d
```

**6. Run the application**

```bash
 go build
 go run main.go
```

## What's specific about project's structure?
This bot's starting point is `main.go`, where configurations are retrieved, and methods for connecting to PostgreSQL and Discord websocket are called. I tried to keep this file clean and clear for easy navigation starting from entry point.

Entry point of bot's main functional is `internals/bot/bot.go`, where event Discord connections are made handlers are added.

Event handlers are specified in `internals/handlers` folder, while methods used by those handlers are in `internals/commands`.

External services like AWS S3 are located inside `internals/service` directory.

In this project, I used `sqlx`, an extension to Golang's `database/sql`. However, I chose not to use any ORM since there was no need.

Also, errors are handled gracefully and `log/slog` module is used for structured logging during the execution

## Miscellaneous
Thank you for such a cool task :D

Working on this project was a wonderful experience as it showed me how building a discord bot is simple and interesting