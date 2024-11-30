# Gator

## Overview

Gator is a simple application that allows to subscribe on RSS feeds and read the last posts


## Features

* Supports multiple users 
* Allows to add specific feeds
* Follow/unfollow feeds
* List all feeds a user subscribed on 
* Read recent posts from the particular feeds


## Getting Started 

### Prerequisites

* Go 1.2 or higher
* Postgresql installed


### Installation 

Install the application:

```
go install github.com/obanoff/gator
```
Or clone the repo, and run from there:

```
git clone https://github.com/obanoff/gator.git
```


### Configuration 

Users credentials stored in >.gatorconfig.json file, so you need to create it in your home directory:

```
touch ~/.gatorconfig.json
```


## Usage 

### Running the application

If installed, you can run the app like so:

```
gator <command> <options>
```

If clonned:

```
go run main.go <command> <options>
```


### CLI Commands

Register a new user and log in:

```
gator register <username>
```

Log in (for existing users):

```
gator login <username>
```

Remove all users:

```
gator reset
```

List all users registered (current user highlighted):

```
gator users
```

Add feed:

```
gator addfeed <name> <url>
```

List all feeds:

```
gator feeds
```

Aggregate (fetch the actual posts from feeds and print them):

```
gator aggr <time_interval>
```

Follow:

```
gator follow <url>
```

Unfollow:

```
gator unfollow <url>
```

List all feeds the current user follows:

```
gator following
```

View all the posts from the feeds the user follows:

```
gator browse <limit>
```



