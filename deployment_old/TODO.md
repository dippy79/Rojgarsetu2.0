# RojgarSetu 2.0 Docker Fix - TODO Steps

## Completed Steps: (Track progress by editing)

## Pending Steps:

1. [ ] Run the structure fix script:
   ```
   cd deployment
   .\fix_structure.ps1
   ```
   This moves services/ into deployment/services/, copies .env.

2. [ ] Validate docker-compose config:
   ```
   cd deployment
   docker compose config
   ```

3. [ ] Verify full deployment:
   ```
   .\verify_compose.ps1
   ```
   This cleans, builds, starts, checks health.

4. [ ] Manual checks if needed:
   - docker compose ps
   - docker compose logs backend
   - curl http://localhost:8083/health

## Notes:
- Ensure Docker is running.
- All services should show healthy, no path/env errors.
- Backend connects to postgres/redis via service names.

Updated: $(Get-Date)
