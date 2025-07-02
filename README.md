# lock-bot

**lock-bot** is a free and open-source Telegram bot (client) for interacting with **lock-box**, a simple self-hosted file server.
Built in Go with modular architecture and minimal dependencies.

## Features

* Upload, list, download, and delete files via Telegram
* Authenticated access using Telegram ID
* Works with any lock-box API endpoint
* Written in pure Go (no external frameworks except Telegram API lib)
* Built-in logger wrapper
* Docker-ready backend (lock-box)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/lnkssr/lock-bot.git
cd lock-bot
```

### 2. Set up Environment Variables

Create a `.env` file in the root directory or set variables manually:

```bash
# Telegram bot token
TOKEN="0123456789:abcdefghijklmnopqrstuvwxyz"

# lock-box API URL
API_URL="http://localhost:5000/api/"
```

### 3. Run the Bot

**Initialize and start the bot (build & run):**

```bash
make init
```

**Development mode (rebuild and run with debug logging):**

```bash
make dev
```

**Clean up build artifacts:**

```bash
make delete
```

## License

**lock-bot** is released under the **AGPLv3** license.
See [LICENSE](./LICENSE.md) for more information.
