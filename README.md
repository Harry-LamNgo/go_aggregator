## Prerequisites

Before running gator, you'll need the following installed on your machine:

- **Go** (1.22 or later) — https://go.dev/dl/
- **PostgreSQL** (14 or later) — https://www.postgresql.org/download/
- **Goose** (database migration tool) — install with:

      go install github.com/pressly/goose/v3/cmd/goose@latest

Verify your installations with:

    go version
    psql --version
    goose --version

## Installation

Clone this repository and build gator executable:

    git clone https://github.com/Harry-LamNgo/go_aggregator.git
    cd go_aggregator
    go build -o ~/go/bin/gator  # or windows: go build -o %USERPROFILE%\go\bin\gator.exe

This compiles gator into a standalone binary and places it in your Go bin directory (usually `$HOME/go/bin`). Make sure that directory is on your PATH:

    export PATH="$PATH:$(go env GOPATH)/bin"

## Configuration

gator expects a configuration file (e.g. `config.yaml`) with the following fields:

    db_url: "postgresql://user:password@localhost:5432/gator?sslmode=disable"
    current_user: "your_username"

- **`db_url`** — The PostgreSQL connection string.

- **`current_user`** — The username to use for authenticated commands.

## Database Setup

Create the database, then run the Goose migrations:

    createdb gator

Navigate to the directory containing your schema files and run migration up to latest version of database (005_posts.sql):

    cd sql/schema/
    goose postgres "<connection_string>" up-to 5

## Running gator

After building or installing, run gator with a command:

    gator <command> [args...]

### Available Commands

**User Management:**
- **`register <username>`** - Create a new user account (and set that user account as current user)
- **`login <username>`** - Authenticate as a user (sets the active session)
- **`users`** - List all registerd users
- **`reset`** - Clear all users data

**Feed Management:**
- **`addfeed <feed_name> <feed_url>`** - Add a new RSS feed to current user _(requires login)_
- **`feeds`** - List all feeds in the system
- **`follow <feed_url>`** - Follow a specific feed through url that current user does not follow but that RSS feed is already existed in the system _(requires login)_
- **`following`** - List RSS feeds followed by the current user _(requires login)_
- **`unfollow <feed_url>`** - Stop following a RSS feed _(requires login)_

**Content:**
- **`agg <time_between_reqs>`** - Aggregate and fetch posts from all followed feeds
- **`browse [<limit]`** - Browse posts from followed feeds (default limit is 2) _(requires login)_

### Example Workflow

```bash
# Register and login
gator register username
gator login username

# Add a feed
gator addfeed "The Verge" "https://www.theverge.com/rss/index.xml"

# Follow feed of other users (feed that is added from another account not current account)
# Gather all feeds in system
gator feeds
(
    terminal shows
    ==================================================
    ** Feed ID: d641691a-f9d9-420a-a21f-fc19e5b5a2b3
    - Name: Hacker News
    - URL: https://news.ycombinator.com/rss
    - User Followed: PoPanda
    ==================================================
    ** Feed ID: c0e0805e-45b6-441f-adc7-546602fb1f59
    - Name: The Verge
    - URL: https://www.theverge.com/rss/index.xml
    - User Followed: HarryKing
)

(terminal current user - HarryKing)
gator follow "https://news.ycombinator.com/rss"

# Check current user following feeds
gator following
(
    terminal shows
    User HarryKing following feeds:
    - The Verge
    - Hacker News
)

# Aggregate and browse content
gator agg 10s
(
   the system is fetching new posts every 10 seconds
)
Important: agg command is stuck in an infinite loop
-> Use Ctrl+C to stop the system

gator browse

```