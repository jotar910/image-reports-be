# Image Reports API
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)
[![License: ISC](https://img.shields.io/badge/License-ISC-yellow.svg)](#)

> This API is designed for managing user-generated images, providing tools to prevent offensive content and ensuring a safer digital environment.

## Get Started

To launch the application:
```sh
sh scripts/start.sh
```

To populate databases:
```sh
sh scripts/seed.sh
```

You can also run each service independently using `Air` or `go`:

```sh
air -c api-gateway/.air.toml
# OR
go build -o ./dist/api-gateway ./api-gateway && ./dist/api-gateway -m prod
```

## Project Overview

This backend application is part of a broader project designed to practice advanced concepts in software development, including microservices, message streaming architecture, server-sent events (SSE), JWT authentication, and Golang backend development. The project is ongoing and continually evolving.

### Frontend Backoffice
For the corresponding frontend backoffice, please visit our repository: [Image Reports FE Backoffice](https://github.com/jotar910/image-reports-fe-backoffice).
This frontend integration showcases the practical application of the backend services in a user-facing environment.

## API Overview

- **User Authentication (POST /v1/auth/login):** Login and receive a JWT token for user authentication.
- **Report List (GET /v1/reports):** Retrieve a paginated list of reports with query parameters for pagination.
- **Individual Report (GET /v1/reports/:id):** Fetch details of a specific report.
- **Create Report (POST /v1/reports):** Submit a new report to the system.
- **Update Report (PATCH /v1/reports/:id):** Modify the approval status of a report (Admins only).
- **Retrieve Image (GET /v1/storage/:id):** Access the image associated with a report.

For a comprehensive guide and detailed definition of all API endpoints, please refer to our full API documentation.
This includes in-depth information on request formats, response structures, and authentication procedures.
It's a valuable resource for developers looking to integrate with our API or understand its capabilities in detail.

[Access the Full API Documentation here](./docs/api.yaml)

## Architecture Planning

![Architecture Diagram](./docs/diagram.png)

### Architecture and Design
The system is structured around a **Microservices Architecture** and adheres to **Hexagonal Architecture** principles for scalability, maintainability, and loose coupling. Services are designed to scale independently, facilitating efficient resource management.

### Services
* **Gateway**: Serves as the primary entry for all API requests, handling user authentication, authorization, and data aggregation.
* **Realtime Service**: Utilizes Server-Sent Events (SSE) to push real-time updates to clients, ensuring immediate reflection of state changes.
* **Storage Service**: Manages file storage and operations, acting as a bridge between the application and the file system.
* **User Service**: Responsible for user account management, including data handling and operations related to user accounts.
* **Processing Service**: Applies AI to assess and process reports, playing a key role in content analysis and moderation.
* **Reporter Service**: Manages the entire lifecycle of reports from creation to data retrieval and status updates.

### RESTful API
The system's communication leverages RESTful principles, which standardize interactions between microservices and external clients, simplifying the data exchange process.

