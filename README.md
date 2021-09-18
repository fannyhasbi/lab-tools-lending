# Lab Tools Lending

# Configuration
Create a new file and name it `.env`. Copy the content from `.env.example` file to `.env` and change the values.

## Database Configuration
### Migration
This project use [golang-migrate](https://github.com/golang-migrate/migrate) tool to make migration. Please install the tool before running these commands in development environment.

**Migrating Up**
```
make migrate-up
```

**Creating New Migration**
```
migrate create -ext sql -dir database/migration -seq example_create_users
```

**Dirty Error**
Sometimes the migration is failed and raise a dirty error when the migration command is being executed again. The error could look like this.
```log
2021/05/18 14:53:59 error: Dirty database version 5. Fix and force version.
make: *** [Makefile:27: migrate-up] Error 1
```

Pay attention to the version. The above example error appears on **version 5**, so you have to force a rollback to previous migration version, which is **version 4**.
```bash
VERSION=4 make migrate-force
```

## Testing
```
make test
```

## Deployment
### Staging
1. Build and run the container
```bash
make container
```
2. Run ngrok on port 3000, then copy the HTTPS url.
  ```bash
  make ngrok
  ```
3. Set the webhook url to the newest ngrok url.
  ```bash
  make change-server URL=https://example.com
  ```

### Production
1. `heroku container:push web -a app-name`
2. `heroku container:release web -a app-name`