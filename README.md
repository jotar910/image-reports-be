# Welcome to Image Reports API 👋
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)
[![License: ISC](https://img.shields.io/badge/License-ISC-yellow.svg)](#)

> API for managing user generated images and prevent offensive content!

## Start

In order to start the whole application you'll need to run:
```sh
sh scripts/start.sh
```

After that, you can populate the databases by running:
```sh
sh scripts/seed.sh
```

You can also start each service separately with `Air` or building and running each service with `go`.

```sh
air -c api-gateway/.air.toml
# OR
go build -o ./dist/api-gateway ./api-gateway && ./dist/api-gateway -m prod
```

## API

- #### `POST /v1/auth/login`
  - User login on the application, returns a JWT Token to be used to authenticate users.

```json
{
  "user": {
    "id": 1,
    "email": "admin@email.com",
    "role": "ADMIN"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjM0MzAxMTMsImlkIjoxLCJlbWFpbCI6ImFkbWluQGVtYWlsLmNvbSIsInJvbGUiOiJBRE1JTiJ9.5wEQ4Y3C42FeYJI0jDpDrXGV0xoF-nqUSQ3l5dBi7Ek"
}
```

- #### `GET /v1/reports`
  - Returns a paginated list of reports.
  - Needs Authorization header.
  - Query Parameters
    - page (optional): `int`
    - count (optional): `int`
  - Response Payload - `PaginatedList<ReportItem>`


```json
{
  "content": [
    {
      "id": 2,
      "name": "Report #2",
      "user": "admin@email.com",
      "image": "image-2-uuid",
      "status": "PUBLISHED",
      "approval": {
        "userId": 1,
        "status": "REJECTED",
        "date": 1663340018
      },
      "evaluation": {
        "grade": 20,
        "categories": [
          "Mythic",
          "Violence",
          "Abuse"
        ]
      },
      "creationDate": 1663340018
    }
  ],
  "page": 1,
  "totalPages": 2,
  "totalElements": 2,
  "numberOfElements": 1
}
```

- #### `GET /v1/reports/:id`
  - Returns a single report for the given id.
  - Needs Authorization header.
  - Parameters
    - id: `int`
  - Response Payload - `ReportItem`

```json
{
  "id": 2,
  "name": "Report #2",
  "user": "user1@email.com",
  "image": "image-2-uuid",
  "status": "PUBLISHED",
  "approval": {
    "userId": 1,
    "status": "REJECTED",
    "date": 1663340018
  },
  "evaluation": {
    "grade": 20,
    "categories": [
      "Mythic",
      "Violence",
      "Abuse"
    ]
  },
  "creationDate": 1663340018
}
```

- #### `POST /v1/reports`
  - Creates a new report on the system.
  - Needs Authorization header.
  - Request Body - multipart/form-data with the following fields
    - name: `string`
    - callback: `string`
    - type: `"file" | "url`
    - url (required if type is url): `string`
    - file (required if type is file): `binary`

- #### `PATCH /v1/reports/:id`
  - Replaces the approval state for the given report.
  - Needs Authorization header.
  - Only available for Admins.
  - Request Payload - `PatchApprovalStatus`

```JSON
{ "approvalStatus":"APPROVED" }
```

- #### `GET /v1/storage/:id`
  - Returns the image that corresponds to a given report id.
  - Needs Authorization header.
  - Parameters
    - id: `int`
  - Response payload: `binary`.

## Architecture Planning


![Architecture Diagram](./docs/diagram.png)
