## gator

gator is a basic cli rss-aggregator that allows you to follow feeds for multiple users, based on the boot.dev tutorial "Build a Blog Aggregator". The backend database used is PostgresSQL.

### Requirements

You'll need to have PostgresSQL installed:

Go: https://go.dev/doc/install\n
PostgresSQL: https://www.postgresql.org/

### Setup

#### Install gator

Install with command "go install github.com/SaschaRunge/gator"

#### PostgresSQL and config

Create a database with Postgres and a user:
https://www.postgresql.org/docs/current/tutorial-createdb.html
You'll likely need psql aswell.

Create a .gatorconfig.json in your home-directory:

Windows: C:\Users\<username>\.gatorconfig.json\n
Linux: ~/.gatorconfig.json

with contents

```
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

The [connection string](https://www.postgresql.org/docs/18/libpq-connect.html#LIBPQ-CONNSTRING-URIS) should look something like this: 
"postgres://user:password@localhost:port/databasename?sslmode=disable"

### Usage

Run with command "gator <command> <args>". The following commands are supported and will show a help message if there are missing arguments:

		"addfeed":      add a feed for the currently logged in user
		"agg":          fetch data for the currently subscribed feeds
		"browse":       shows a number of the most recent feeds
		"feeds":        shows all available feeds
		"follow":       subscribe to a feed
		"following":    shows the feeds the logged in user is subscribed to
		"login":        log the specified user in
		"register":     register a new user. the user will be automatically logged in
		"reset":        this will reset your database. do not use. i was to lazy to research how to remove it for production
		"unfollow":     unfollows from a feed
		"users":        shows the known users





