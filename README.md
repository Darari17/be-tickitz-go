# Belalai E-Wallet Backend

![badge golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![badge postgresql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![badge redis](https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white)

Welcome to Tickitz Movies!
The ultimate movie booking backend system designed to deliver a seamless and efficient cinema experience. Just like a well-directed film where every scene matters, Tickitz Backend ensures every request, response, and transaction runs smoothly behind the screen.

Built with Go (Golang) for performance and scalability, Tickitz Movies provides powerful APIs to manage users, movies, showtimes, and bookings. This backend serves as the core engine for the Tickitz Frontend Web Application, handling authentication, payment integration, and data persistence with precision and reliability.

From listing movies to confirming your seat, Tickitz Backend keeps the show running—fast, secure, and always in sync. Sit back, relax, and let the backend handle the plot.

## 🔧 Tech Stack

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/install-redis-on-windows/)
- [JWT](https://github.com/golang-jwt/jwt)
- [argon2](https://pkg.go.dev/golang.org/x/crypto/argon2)
- [migrate](https://github.com/golang-migrate/migrate)
- [Docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- [Swagger for API docs](https://swagger.io/) + [Swaggo](https://github.com/swaggo/swag)

## 🗝️ Environment

```bash
# database
DBUSER=<your_database_user>
DBPASS=<your_database_password>
DBNAME=<your_database_name
DBHOST=<your_database_host>
DBPORT=<your_database_port>

# JWT hash
JWT_SECRET=<your_secret_jwt>
JWT_ISSUER=<your_jwt_issuer>

# Redish
RDB_HOST=<your_redis_host>
RDB_PORT=<your_redis_port>

```

## ⚙️ Installation

1. Clone the project

```sh
$ https://github.com/Darari17/be-tickitz-go.git
```

2. Navigate to project directory

```sh
$ cd be-tickitz-full
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Setup your [environment](##-environment)

5. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) for DB migration

6. Do the DB Migration

```sh
$ migrate -database YOUR_DATABASE_URL -path ./db/migrations up
```

or if u install Makefile run command

```sh
$ make migrate-createUp
```

7. Run the project

```sh
$ go run ./cmd/main.go
```

### 📘 API Endpoints

| Method              | Endpoint                   | Auth         | Body / Params                                                                                                                                                             | Description                         |
| ------------------- | -------------------------- | ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| **Auth**            |                            |              |                                                                                                                                                                           |                                     |
| `POST`              | `/auth/login`              |              | `email`, `password`                                                                                                                                                       | Authenticate user                   |
| `POST`              | `/auth/register`           |              | `email`, `password`                                                                                                                                                       | Register new user                   |
| **Profile**         |                            |              |                                                                                                                                                                           |                                     |
| `GET`               | `/profile`                 | Bearer Token | -                                                                                                                                                                         | Get logged-in user profile          |
| `PATCH`             | `/profile`                 | Bearer Token | `{ firstname, lastname, phone_number }`                                                                                                                                   | Update user profile                 |
| `PATCH`             | `/profile/change-avatar`   | Bearer Token | `avatar (file)`                                                                                                                                                           | Upload new profile avatar           |
| `PATCH`             | `/profile/change-password` | Bearer Token | `{ old_password, new_password }`                                                                                                                                          | Change user password                |
| **Movies (Public)** |                            |              |                                                                                                                                                                           |                                     |
| `GET`               | `/movies`                  | -            | `page`, `search`, `genre`                                                                                                                                                 | Get all movies with optional filter |
| `GET`               | `/movies/{id}`             | -            | `id` (path)                                                                                                                                                               | Get movie detail                    |
| `GET`               | `/movies/popular`          | -            | `page`                                                                                                                                                                    | Get popular movies                  |
| `GET`               | `/movies/upcoming`         | -            | `page`                                                                                                                                                                    | Get upcoming movies                 |
| `GET`               | `/movies/genres`           | -            | -                                                                                                                                                                         | Get all available genres            |
| **Admin - Movies**  |                            |              |                                                                                                                                                                           |                                     |
| `GET`               | `/admin/movies`            | Bearer Token | -                                                                                                                                                                         | Get all movies (admin)              |
| `POST`              | `/admin/movies`            | Bearer Token | `multipart/form-data` — includes `title`, `overview`, `director_name`, `duration`, `release_date`, `popularity`, `poster`, `backdrop`, `genres[]`, `casts[]`, `schedules` | Create new movie                    |
| `GET`               | `/admin/movies/{id}`       | Bearer Token | `id` (path)                                                                                                                                                               | Get movie detail by ID              |
| `PATCH`             | `/admin/movies/{id}`       | Bearer Token | `multipart/form-data` — update movie fields                                                                                                                               | Update movie                        |
| `DELETE`            | `/admin/movies/{id}`       | Bearer Token | `id` (path)                                                                                                                                                               | Soft delete movie                   |
| **Orders**          |                            |              |                                                                                                                                                                           |                                     |
| `POST`              | `/orders`                  | Bearer Token | `{ email, fullname, phone, payment_id, schedule_id, seat_codes[] }`                                                                                                       | Create a new order                  |
| `GET`               | `/orders/{id}`             | Bearer Token | `id` (path)                                                                                                                                                               | Get order detail                    |
| `GET`               | `/orders/history`          | Bearer Token | -                                                                                                                                                                         | Get user order history              |
| `GET`               | `/orders/cinemas`          | Bearer Token | -                                                                                                                                                                         | Get all cinemas                     |
| `GET`               | `/orders/locations`        | Bearer Token | -                                                                                                                                                                         | Get all locations                   |
| `GET`               | `/orders/payments`         | Bearer Token | -                                                                                                                                                                         | Get all payment methods             |
| `GET`               | `/orders/schedules`        | Bearer Token | `movie_id` (query)                                                                                                                                                        | Get schedules by movie ID           |
| `GET`               | `/orders/seats`            | Bearer Token | `schedule_id` (query)                                                                                                                                                     | Get available seats                 |
| `GET`               | `/orders/times`            | Bearer Token | -                                                                                                                                                                         | Get available movie times           |

---

## 📄 LICENSE

MIT License

Copyright (c) 2025 Belalai team

## 📧 Contact Info & Contributor

[https://github.com/Darari17](https://github.com/Darari17)

## 🎯 Related Project

[https://github.com/Darari17/fe-tickitz-react](https://github.com/Darari17/fe-tickitz-react)
