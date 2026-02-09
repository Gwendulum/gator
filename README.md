Gator is a CLI blog aggregator that periodically collects RSS feeds, and then stores and organizes them with postgreSQL. 

INSTALLATION

In order to compile gator on your machine, you first need to:
    1. install Go (https://go.dev/doc/install)
    2. install postgreSQL 

Then you install gator from the command line using go install github.com/Gwendulum/gator@latest

CONFIGURATION

Gator needs a .gatorconfig.json file in your home directory. 

1. first create the .gatorconfig.json
2. then copy in the address to your database in this format:
    
{
  "db_url": "postgres://<user>:<password>@localhost:5432/<database>sslmode=disable"
}

This file will also hold the name of the currently logged in user once one has been created.

USAGE

Once installed you operate gator using a list of commands. The [brackets] are necessary arguments for that command.
* register ["username"]         ---registers a local user
* login ["username"]            ---logs in a registered user
* users                         ---shows the list of users
* feeds                         ---shows current feeds
* addfeed [feed name][url]      ---adds feed to track
* following                     ---shows current user's followed feeds
* unfollow [url]                ---unfollow a url
* browse                        ---browse currently stored posts