# Gator

A command-line RSS feed aggregator built with Go and PostgresSQL.

Gator allows you to follow RSS feeds, automatically fetches new posts and lets you read them 
in the CLI.

## Features

- multi-user support
- follow multiple RSS feeds
- automatic feed aggregation
- PostgresSQL storage for posts and feeds

## Installation

1. Clone the repository:

```bash
git clone https://github.com/eqo-m/gator.git
cd gator
```

2. Set up the database :

```bash
createdb gator
psql < schema.sql
```

3. Build the application:
```bash
go build -o gator
```

4. Create your config file (`~/.gatorconfig.json`):
```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Quick Start

```bash
./gator register alice
./gator addfeed "Hacker News RSS" "https://hnrss.org/newest"
./gtor agg 1m
```

## Usage

###Users: 
**Register a new user**

```bash
./gator register alice
```
**Login as an existing user:**
```bash
./gator login alice
```
**List all users:**
```bash
./gator users
```
Shows all registered users with (current) next to the active user.

### Feed Management

**Add a new feed:**
```bash
./gator addfeed  "Feedname" "FeedURL"
```

**List all available feeds:**
```bash
./gator feeds
```

**Follow an existing feed:**
```bash
./gator follow "FeedURL"
```

**Unfollow a feed:**
```bash
./gator unfollow "FeedURL"
```

**List feeds you're following:**
```bash
./gator following
```

### Aggregation

**Start the aggregator:**
```bash
./gator agg 10s
```
Fetches posts from all feeds at the specified interval.

## Configuration

The configuration file is located at `~/.gatorconfig.json`:
```json
{
  "db_url": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
  "current_user_name": "alice"
}
```
- `db_url`: PostgreSQL connection string
- `current_user_name`: The currently active user (managed by `login` command)