# Secrets Management

## Development (Environment Variables)

For development, use `.env` file (already in `.gitignore`):

```env
POSTGRES_USER=your-database-user
POSTGRES_PASSWORD=generate-a-unique-password
POSTGRES_DB=rojgarsetu2
JWT_SECRET=generate-a-unique-random-secret
REFRESH_TOKEN_KEY=generate-a-different-unique-random-secret
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/rojgarsetu2?sslmode=disable
REDIS_URL=redis://localhost:6379
```

## Production (Docker Secrets)

For production, use Docker secrets with `docker-compose.prod.yml`:

### Create Secrets

```bash
openssl rand -base64 48 | docker secret create jwt_secret -
openssl rand -base64 48 | docker secret create refresh_token_key -
# Create db_user and db_password from your organization's secret manager.
```

### Deploy with Secrets

```bash
docker compose -f deployment/docker-compose.prod.yml up -d
```

### Secret Usage in Services

- PostgreSQL: Uses `POSTGRES_USER_FILE` and `POSTGRES_PASSWORD_FILE`
- Backend: Uses `JWT_SECRET_FILE` and `REFRESH_TOKEN_KEY_FILE`; both must be unique secrets
- Backend: Requires `ENVIRONMENT=production`; authentication cookies are marked Secure
- Database URL constructed from secrets in production

### Security Notes

- Never commit secrets to git
- Use strong, unique secrets for each environment
- Rotate secrets regularly
- Use secret management services (AWS Secrets Manager, HashiCorp Vault) for production
- Enable RBAC for secret access