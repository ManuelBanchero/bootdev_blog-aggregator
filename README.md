## Instalation

- PostgreSQL:

```Shell
brew install postgresql@15
```

- PostgreSQL Server:

```Shell
brew services start postgresql@15
```

- Project Packages:

```Shell
go install
```

## Create the follow file to use the Project

```Shell

cd ~ && touch .gatorconfig.json && nano ~/.gatorconfig.json
```

Inside of it, add the follow:

```json
{
  "db_url": "postgres://yourMachineUsername@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Then:

```Shell
^O
enter
^X
```

## Project Commands

- `gator login [username]`: Set the actual user to `.gatorconfig.json`.

- `gator register [username]``: Create a new user

- `gator reset`: Reset de user database

- `gator users`: List of all users

- `gator agg [time between requests]`: Loops follow urls and save posts to DB.

- `gator addfeed [Title] [url]`: Create a new feed using a Title and a URL. Automatically follows the new url.

- `gator feeds`: Get all feeds

- `gator follow [url]`

- `gator following`: List all the following urls

- `gator unfollow [url]`

- `browse [optional: response limit]`: List of the posts saved at the DB
