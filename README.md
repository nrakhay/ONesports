# ON Esports Bot for Discord

Welcome to the implementation of the Discord Bot for [ONEsports.gg](https://onesports.gg/). To get started with using this bot, please follow the setup instructions below.

## Project Setup

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

**5. Run the application**

```bash
 go build
 go run main.go
```

## What this bot can do?

To be written :D
