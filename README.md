# Gator

A CLI RSS feed aggregator written in Go. Gator lets you register users, subscribe to RSS feeds, aggregate posts on a schedule, and browse the latest content — all from your terminal.

---

## Dependencies

| Dependency | Purpose |
|---|---|
| [Go](https://go.dev/) ≥ 1.21 | Language runtime |
| [PostgreSQL](https://www.postgresql.org/) | Persistent storage for users, feeds, and posts |
| [`github.com/lib/pq`](https://github.com/lib/pq) | PostgreSQL driver for Go's `database/sql` |
| [`github.com/google/uuid`](https://github.com/google/uuid) | UUID generation for database records |
| [`github.com/araddon/dateparse`](https://github.com/araddon/dateparse) | Flexible RSS publish-date parsing |

---

## Installation

**Install directly with Go:**

```bash
go install https://codeberg.org/TerraLambda/gator@latest
```

**Or clone and build from source:**

```bash
git clone https://codeberg.org/TerraLambda/gator.git
cd gator
go build -o gator .
```

**Create the database:**

```sql
CREATE DATABASE gator;
```

---

## Configuration

Gator reads its config from a file named `.gatorconfig.json` in your **home directory** (`~/.gatorconfig.json`).

**Example:**

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "alice"
}
```

| Field | Description |
|---|---|
| `db_url` | PostgreSQL connection string. Update the username, password, host, port, and database name to match your setup. |
| `current_user_name` | The currently logged-in user. Updated automatically by `login` and `register`, but can also be edited manually. |

---

## Commands

### `register <username>`

Creates a new user and automatically logs in as that user.

```bash
gator register alice
```

---

### `login <username>`

Switches the active session to an existing user.

```bash
gator login alice
```

> Requires the user to already exist. Use `register` to create a new one.

---

### `users`

Lists all registered users. The currently logged-in user is marked with `(current)`.

```bash
gator users
```

---

### `reset`

Wipes all data from the database. Takes no arguments.

```bash
gator reset
```

> ⚠️ This is destructive and cannot be undone.

---

### `addfeed <name> <url>`

Adds a new RSS feed to the database and automatically follows it as the current user. Requires a logged-in user.

```bash
gator addfeed "Hacker News" https://news.ycombinator.com/rss
```

| Argument | Description |
|---|---|
| `name` | A display name for the feed |
| `url` | The full URL of the RSS feed |

---

### `feeds`

Prints all feeds in the database, showing each feed's name, URL, and the user who added it.

```bash
gator feeds
```

---

### `follow <url>`

Subscribes the current user to an existing feed by its URL. The feed must already be in the database (use `addfeed` to add new ones). Requires a logged-in user.

```bash
gator follow https://news.ycombinator.com/rss
```

---

### `following`

Lists all feeds the current user is subscribed to. Requires a logged-in user.

```bash
gator following
```

---

### `unfollow <url>`

Unsubscribes the current user from a feed by its URL. Requires a logged-in user.

```bash
gator unfollow https://news.ycombinator.com/rss
```

---

### `agg <interval>`

Starts the feed aggregator, fetching and storing new posts from all feeds on a repeating interval. Runs continuously until interrupted. Takes a Go duration string as the interval.

```bash
gator agg 30s
gator agg 5m
gator agg 1h
```

| Argument | Description |
|---|---|
| `interval` | How often to fetch feeds, as a Go duration (e.g. `30s`, `5m`, `1h`) |

> Run this in a separate terminal session to keep feeds up to date in the background.

---

### `browse [limit]`

Prints the latest posts from feeds the current user follows. Requires a logged-in user.

```bash
gator browse        # defaults to 2 posts
gator browse 10     # shows up to 10 posts
```

| Argument | Description |
|---|---|
| `limit` *(optional)* | Number of posts to display. Defaults to `2` if omitted. |

Each post is shown with its title and a plain-text version of its description.
