# A Simple Bank

### Features

- Create and manage account
  - Owner
  - Balance
  - Currency
- Record all balance changes for each account
  - Create an account entry for each change for each account
- Money transfer transaction
  - Perform money transfer between 2 accounts consistently within a transaction

### Pre-requisites

- Install `docker for desktop`
- Execute `brew install golang-migrate sqlc`

### Database Design

- Design DB schema using dbdiagram.io
  - Export the queries onto `/dbdiagrams.io`
- Save and share DB diagram within the team
- Generate SQL code to create database in a target database engine i.e. postgres/MySQL/SQLServer

### Docker and Postgres

- Execute `docker pull postgres:12-alpine` to get the postgres image
- Execute `docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine` to run the postgres container
- Execute `docker logs postgres12` to see the logs
- Execute `docker exec -it postgres12 psql -U root` to connect to the postgres container and login as `root` user
- Connect to postgres container and execute the queries from `/dbdiagrams.io` to create the tables

### DB migration

- Execute `migrate -version` to verify that the `golang-migrate` has been installed
- Execute `migrate create -ext sql -dir db/migration -seq init_schema` to generate migration files
  - `*.up.sql` is used to migrate up to a new version using `migrate up`
  - `*.down.sql` is used to migrate down to an older version using `migrate down`
  - Copy the sql quesries generated from `dbdiagram.io` to `*.up.sql`
  - Add `DROP TABLE` queries to `*.down.sql`
- Execute `make migrateup` to migrate data upwards to a new version
- Execute `make migratedown` to revert migration to a previous version
- Manage migrations in future with `migrtion up/down` commands

### DB and Docker Setup for development

- Execute `make postgres` to run postgres container on local docker setup
- Execute `make createdb` to create the `simple_bank` database
- Execute `make migrateup` to setup tables and initial database state
- If required,
  - Execute `make dropdb` to drop database
  - Execute `make migratedown` to migrate or revert database state to a previous version

### Generate CRUD Golang code from SQL

- Execute `make sqlc` to auto generate CRUD functionalities
