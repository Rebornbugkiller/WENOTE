# WeNote Project Rules

## Project Structure

- `wenote-backend/` - Go 1.23 backend (Gin + GORM + MySQL 8.0)
- `wenote-frontend/` - Vue 3 frontend (Vite + Element Plus + Pinia + TailwindCSS)

## Deployment

Production server: http://47.97.78.32

- Frontend: Nginx, port 80
- Backend: wenote service, port 8080
- Database: MySQL, port 3306

### Update & Deploy

```bash
# Build frontend locally
cd wenote-frontend && npm run build

# Build backend locally (cross-compile for Linux)
cd wenote-backend && GOOS=linux GOARCH=amd64 go build -o server-linux ./cmd/server

# Upload frontend
scp -r wenote-frontend/dist/* root@47.97.78.32:/opt/wenote/frontend/

# Upload backend
scp wenote-backend/server-linux root@47.97.78.32:/opt/wenote/backend/server

# Restart backend service
ssh root@47.97.78.32 "systemctl restart wenote"
```

### Test Account

- Username: test
- Password: test123
