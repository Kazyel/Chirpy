# Chirpy API

This is the API built for the **Chirpy application**.

It is written in **Go** and uses these two tools: **[sqlc](https://github.com/kyleconroy/sqlc)** to generate the database queries and models and **[goose](https://github.com/pressly/goose)** to manage the database migrations.

## Setup

### Database

The database is managed by **[PostgreSQL](https://www.postgresql.org/)**. You can install it on your machine using the following command (Mac & Linux):

```bash
brew install postgresql
sudo apt-get install postgresql
```

Once installed and started, you can create a database named `chirpy` with the following command:

```bash
createdb chirpy
```

Then, you should run the migrations to create the tables and columns:

```bash
goose postgres dburl -dir ./sql/schema up
```

Now, you can create a .env file in the root of the project with the following content:

```bash
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=secret
POLKA_KEY=polka
```

## API Endpoints & Resources

For more detailed documentation, please refer to the [API documentation](https://github.com/Kazyel/chirpy-bootdev/blob/main/docs/api.md).
