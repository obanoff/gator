# Gator

## Overview

**Gator** is a simple application that allows to subscribe on RSS feeds and read the last posts

---

## Features

- Supports multiple users. 
- Allows to add specific feeds.
- Follow/unfollow feeds.
- List all feeds a user subscribed to.
- Read recent posts from the particular feeds.

---

## Getting Started 

### Prerequisites

Before installing Gator, ensure you have the following:

- **Go** 1.2 or higher.
- **Postgresql** installed.
- **Bash** (for migrations scripts).

---

### Installation 

Install the application:

```bash
go install github.com/obanoff/gator
```
Or clone the repo, and run from there:

```bash
git clone https://github.com/obanoff/gator.git
```

---

### Configuration 

1. Users credentials    
Create `.gatorconfig.json` file in your home directory to store user credentials:

```bash
touch ~/.gatorconfig.json
```

2. Database URL    
Specify the database connection string in a `.env` file:

```bash
echo "DATABASE_URL=postgres://username:password@localhost:port/dbname" >> .env
```

3. Run DB migrations    
Run the migration script to set up your database schema: 

```bash
./migrate-up.sh
```

---

## Usage 

### Running the application

* If installed via > go install:

```bash
gator <command> <options>
```

* If cloned from the repository:

```bash
go run main.go <command> <options>
```

---

### CLI Commands

| Command                     | Description                                     |
|-----------------------------|-------------------------------------------------|
| `gator register <username>` | Register a new user and log in.                 |
| `gator login <username>`    | Log in with an existing user.                   |
| `gator reset`               | Remove all users.                               |
| `gator users`               | List all registered users (current user highlighted). |
| `gator addfeed <name> <url>`| Add a new RSS feed with a given name and URL.   |
| `gator feeds`               | List all RSS feeds available.                  |
| `gator aggr <time_interval>`| Fetch recent posts and display them.           |
| `gator follow <url>`        | Follow a specific feed by its URL.             |
| `gator unfollow <url>`      | Unfollow a specific feed by its URL.           |
| `gator following`           | List all feeds the current user follows.       |
| `gator browse <limit>`      | View posts from feeds the user follows.        |



