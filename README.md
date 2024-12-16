# blog_aggregator
CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

Note that you'll need Go, Postgres, Goose and SQLC installed to run the program.

To install:
1. Clone the repo locally.
2. Navigate to the blog-aggregator directory
3. Run sqlc generate to generate required objects
4. Create a postgres database to use for the project and run "goose postgres "CONNECTION_STRING" up" until you are up to the latest database migration
4. Run "go install ."
5. Create a "~/.gatorconfig.json" file in your ~ directory pointing to your database, example:
    {"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user":"admin"}
6. You can now use the "blog-aggregator" command to browse RSS feeds. Run "blog-aggregator help" for a list of commands.
