You need Postgres and Go installed to run the program.
Install the gator CLI using go install.
Set up the config file (.gatorconfig.json as {
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}) and database schema, then run the program. 
You can run the commands: login, register, reset, users, agg, addfeed, feeds, follow, following, unfollow, browse
