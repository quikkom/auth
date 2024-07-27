# Auth

This repository contains authentication service of Quikkom. It written in Go and uses PostgreSQL as database. It implements protobuf [definitions](https://github.com/quikkom/auth-proto).

## How to run?

Auth service needs a PostgreSQL database which has a schema described within the SQL [files](https://github.com/quikkom/auth/tree/master/sql).

> You can start a PostgreSQL instance by using Docker;
>
> ```shell
> docker run \
>    --restart always \
>    --name postgres \
>    -e POSTGRES_USER=postgres \
>    -e POSTGRES_PASSWORD=postgres \
>    -e POSTGRES_DB=postgres \
>    -p 5432:5432 \
>    -v postgres_data:/var/lib/postgresql/data \
>    -d \
>    postgres
> ```

Once you got a database instance, apply the SQL scripts to achieve desired schema. Open up a terminal and go to the `sql` directory:

```shell
$ PGPASSWORD="<PASSWORD>" psql -U "<DB_USER>" -h "<DB_HOST>" < 00-init-schema.sql
```

> If you started the database instance using Docker by the command above, you can use this command;
>
> ```shell
> $ PGPASSWORD="postgres" psql -U "postgres" -h "localhost" < 00-init-schema.sql
> ```

Then you have to define environment variables either using .env file with [dotenvx](https://dotenvx.com/) or inside the shell environment.

If you are going to use dotenvx, first of all install it by following the instructions at its website. Then create a .env file:

```dotenv
DATABASE_URL=postgres://username:password@host:port/database?sslmode=disable
AUTH_SECRET="A complex string"
```

> You can take a look to `.env.example` file.

Change database URL according to your connection variables.

Now go to `src` directory and run the service:

```shell
$ cd src
$ dotenvx -f ../.env run -- go run .
```

## Building container image

You can build a container image of the service. Type in terminal at root of the repository:

```shell
$ docker build -t quikkom-auth:$(cat VERSION) .
```

Now you can use the service as a container.

## Continuos Deployment

The repository has a pipeline which automatically build a container and pushes it to [Quikkom's](https://github.com/orgs/quikkom/packages) packages named `auth`.
