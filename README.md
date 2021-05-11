# Lab Tools Lending

## Unit Test
```
make test
```

## Deployment
### Staging
1. Build and run the container
```bash
make container
```
2. Run ngrok on port 3000, thne copy the HTTPS url.
  ```bash
  make ngrok
  ```
3. Set the webhook url to the newest ngrok url.
  ```bash
  make change-server URL=https://example.com
  ```