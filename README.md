# Welcome to Image Reports API 👋
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)
[![License: ISC](https://img.shields.io/badge/License-ISC-yellow.svg)](#)

> API for image 

## Start Locally

In order to start the whole application you'll need to run:
```sh
sh scripts/start.sh
```

After that, you can populate the databases by running:
```sh
sh scripts/seed.sh
```

### Users, accounts and roles.

This will setup two containers: server and database.

## Populate Database

After everything is running, you may populate the database with some initial data. For that you'll need to run:
```sh
npm run seed
```

## API

- #### `POST /v1/auth/login`
  - User login on the application, returns a JWT Token to be used to authenticate users.
- #### `GET /v1/ott/:accountId/movies`
  - Returns a list of movies paginated that belongs to some account.
- #### `GET /v1/ott/:accountId/movies/:movieId`
  - Returns the details for a movie belonging to some account.
- #### `CREATE /v1/ott/:accountId/movies`
  - Creates a new movie on some account.
- #### `PUT /v1/ott/:accountId/movies/:movieId`
  - Patches details for a movie on some account.
- #### `DELETE /v1/ott/:accountId/movies/:movieId`
  - Deletes a movie on some account.

This will create all the data for you to start using the application. This data includes:

- ### Users, accounts and roles.

There are two types of roles on this application: Admin and User. The Admin will be able to execute all types of requests, either the resources belong to its account or not. Also, it will have write permissions, and so, he will be able to create, update and delete resources. The User, will only be able to read data that belongs to his account.

It also creates three users: `admin@email.com`, `user1@email.com` and `user2@email.com`. The first is an Admin and the others are Users. The passwords will be generated randomly and printed at the execution point. Each user will have an account associated which contain the *accountId* as `1`, `2` and `3`, respectively.

- ### Actors and directors.

Randomly created. Actors and directors have ID and Name (based on the requirements).

- ### Countries.

Randomly created. Countries have ID and Name (based on the requirements).

- ### Movies.

Randomly created. Long list of movies, which properties were based on the requirements.

## API Endpoints

The API endpoints are fully described in `openapi.yaml` (OpenApi document).

### Summary:

- #### `POST /v1/ott/auth/login`
  - User login on the application, returns a JWT Token to be used to authenticate users.
- #### `GET /v1/ott/:accountId/movies`
  - Returns a list of movies paginated that belongs to some account.
- #### `GET /v1/ott/:accountId/movies/:movieId`
  - Returns the details for a movie belonging to some account.
- #### `CREATE /v1/ott/:accountId/movies`
  - Creates a new movie on some account.
- #### `PUT /v1/ott/:accountId/movies/:movieId`
  - Patches details for a movie on some account.
- #### `DELETE /v1/ott/:accountId/movies/:movieId`
  - Deletes a movie on some account.

** You must always add the Authorization header on the request (except the login endpoint itself) in order to access the endpoints (e.g. `Authorization: Bearer <jwt_token>`).

## Configuration

The execution is configured through environment variables.
You can find this configuration on `.env.docker` file (feel free to modify it).

### General config

| Variable | Description | Default |
|---|---|---|
| NODE_ENV | Execution type | production |
| NODE_PORT | Server listening port | 3000 |

### Database TypeORM

| Variable | Description | Default |
|---|---|---|
| TYPEORM_HOST | Database host | 10.7.0.3 |
| TYPEORM_PORT | Database listening port | 5432 |
| TYPEORM_USERNAME | Database username | root |
| TYPEORM_PASSWORD | Database password | password |
| TYPEORM_DATABASE | Database name | ott_movies |
| TYPEORM_SYNCHRONIZE | Synchronize database entities on start | true|
| TYPEORM_LOGGING | Log database operation to the console | false|
| TYPEORM_ENTITIES | Path to find entities for the synchronization | dist/api/components/**/model.js|

### DBMS Postgres

| Variable | Description | Default |
|---|---|---|
| POSTGRES_USER | DBMS user | root |
| POSTGRES_PASSWORD | DBMS password | password |
| POSTGRES_DB | DBMS database name | ott_movies |

## Install

```sh
npm install
```

## Run tests

```sh
npm test
```

## Author

👤 **João Rodrigues**

* Github: [@jotar910](https://github.com/jotar910)