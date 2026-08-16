# Secrets Management

## Development (Environment Variables)

For development, use `.env` file (already in `.gitignore`):

```env
POSTGRES_USER=amitsharma
POSTGRES_PASSWORD=Asha12@Ashok24
POSTGRES_DB=rojgarsetu2
JWT_SECRET=your-jwt-secret-min-32-chars
DATABASE_URL=postgres://amitsharma:Asha12@Ashok24@localhost:5432/rojgarsetu2?sslmode=disable
REDIS_URL=redis://localhost:6379
```

## Production (Docker Secrets)

For production, use Docker secrets with `docker-compose.prod.yml`:

### Create Secrets

```bash
echo "amitsharma" | docker secret create db_user -
echo "Asha12@Ashok24" | docker secret create db_password -
echo "your-production-jwt-secret-min-32-chars" | docker secret create jwt_secret -
```

### Deploy with Secrets

```bash
docker compose -f deployment/docker-compose.prod.yml up -d
```

### Secret Usage in Services

- PostgreSQL: Uses `POSTGRES_USER_FILE` and `POSTGRES_PASSWORD_FILE`
- Backend: Uses `JWT_SECRET_FILE` for JWT secret
- Database URL constructed from secrets in production

### Security Notes

- Never commit secrets to git
- Use strong, unique secrets for each environment
- Rotate secrets regularly
- Use secret management services (AWS Secrets Manager, HashiCorp Vault) for production
- Enable RBAC for secret access