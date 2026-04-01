# 07 - Deployment and Hosting Strategy

Deploying a Go application like Gost is straightforward because Go compiles down to a single, statically linked binary. This means you do not need a runtime environment (like Node.js or Python) installed on the target server.

---

## 1. Preparing for Production

Before deploying, ensure you compile the application. Go's cross-compilation is powerful.
If you are deploying to a standard Linux server (Ubuntu/Debian) via GitHub Actions or your local machine, run:

```bash
GOOS=linux GOARCH=amd64 go build -o gost-api main.go
```

This produces a file named `gost-api` which is all you need to execute on the server along with your `.env` file!

---

## 2. Hosting Recommendations

Because Gost is incredibly lightweight (usually consuming less than 30MB of RAM at idle), you have several excellent and cost-effective hosting options.

### Option A: Railway.app / Render (PaaS) - Recommended for Beginners/Startups

Platforms as a Service (PaaS) are the easiest way to deploy. They automatically detect Go environments.
**Why choosing them?**

- Zero server config.
- You push to GitHub, they build and deploy.
- Easy to attach managed PostgreSQL and Redis addons with 1-click.

**How to deploy (Render/Railway):**

1. Connect your Github Repository.
2. In the "Build Command", run: `go build -o app main.go`
3. In the "Start Command", run: `./app`
4. Copy your `.env` variables into the platform's "Environment Variables" tab.

### Option B: DigitalOcean Droplets / AWS EC2 (IaaS)

If you require maximum control and lower costs at scale, spinning up a Linux VPS is ideal.

**How to deploy via Docker (Recommended for VPS):**
Since we already have a `docker-compose.yml`, deploying is extremely easy. By adding a `Dockerfile` for the Go App, you can orchestrate everything.

1. **Create a `Dockerfile`** in your root:

```dockerfile
# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]
```

2. **Update your `docker-compose.yml`** on the server to include the app:

```yaml
services:
  api:
    build: .
    ports:
      - "3000:3000"
    env_file: .env
    depends_on:
      - postgres
      - redis
  # ... (postgres and redis remain the same)
```

3. Run `docker-compose up -d --build` on your Droplet!

### Option C: Vercel (Serverless Functions)

**Overview:**
Vercel runs Go as serverless functions. To deploy a monolithic Gin application (like this one) on Vercel, we map all incoming traffic to a single catch-all handler.

**Step 1 — Configure Vercel Routing (`vercel.json`)**
Create a `vercel.json` in your root. A simple rewrite is enough to map everything to your entry point:

```json
{
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/api/index.go"
    }
  ]
}
```

**Step 2 — Exposing the Router in the App**
Your `src/app/app.module.go` must have a `Bootstrap()` function that returns `*gin.Engine` (not one that calls `router.Run()`).

```go
func Bootstrap() *gin.Engine {
    config.LoadEnv()
    config.ConnectDatabase() // Ensure DATABASE_URL is set in Vercel
    // ...
    router := gin.Default()
    // ... setup routes and modules
    return router
}
```

**Step 3 — Handling Translations (The Embed Way)**
Go serverless functions on Vercel cannot always find files on disk. To solve this, we use `go:embed` to compile the translation JSONs into the binary.
1. Put your translation files in `src/common/i18n/locales/*.json`.
2. In your i18n provider, use `//go:embed locales/*.json` to load them.

**Step 4 — Create the Serverless Entrypoint**
Create `api/index.go` at the root of the project:

```go
package handler

import (
    "net/http"
    "api/src/app" // Use "api" as it is the name in your go.mod
)

var engine = app.Bootstrap()

func Handler(w http.ResponseWriter, r *http.Request) {
    engine.ServeHTTP(w, r)
}
```

**Step 5 — Database & Environment Variables**
1. **Critical:** In Vercel Project Settings, add `DATABASE_URL` with your remote connection string (e.g., Supabase or Neon).
2. For Postgres on Vercel, it is recommended to append `?sslmode=require` to your connection string.

**Step 6 — Test and Deploy**
- Push your changes to GitHub (if connected to Vercel).
- Or run `vercel` from your terminal to deploy manually.

**⚠️ Critical Limitations for Vercel:**
- **Persistent Connections:** WebSockets and background goroutines (like `hub.Run`) may not work as expected across different requests.
- **Cold Starts:** The first request after a period of inactivity may take longer as the function boots.
- **Ephemeral Storage:** Local file writing is not supported. Use S3 or similar for uploads.


---

## 3. Reverse Proxies (Nginx / Caddy)

If you deploy on a VPS (Option B), never expose port 3000 directly. Use a reverse proxy to handle SSL/TLS (HTTPS).

**Example with Caddy (Extremely Easy Mux):**
Install Caddy and map your domain to the internal Gost port:

```text
api.yourdomain.com {
    reverse_proxy localhost:3000
}
```

Caddy will automatically provision and renew Let's Encrypt SSL certificates for you!

---

## 4. Environment Checklist for Prod

- Ensure `DATABASE_URL` and `REDIS_HOST` point to the correct production credentials.
- Ensure `ALLOWED_CORS` has your exact Frontend domains listed, removing `*`.
- Set Gin to Release Mode by adding `GIN_MODE=release` to your `.env` file to prevent Gin from printing debug logs into standard output in production.
