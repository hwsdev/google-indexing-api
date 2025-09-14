# 🚀 Quick Start Guide

Simple Google Indexing API tanpa authentication.

## 📦 Deployment

```bash
# Clone project
git clone <your-repo>
cd google-indexing-api

# Deploy dengan Docker
chmod +x deploy.sh
./deploy.sh start
```

## 🔥 Usage

### Submit Single URL
```bash
curl -X POST http://localhost:8080/api/v1/index \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/page",
    "service_account": {
      "type": "service_account",
      "project_id": "your-project-id",
      "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
      "client_email": "service@project.iam.gserviceaccount.com",
      "client_id": "123456789",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token"
    }
  }'
```

### Submit Batch URLs
```bash
curl -X POST http://localhost:8080/api/v1/index/batch \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://example.com/page1", "https://example.com/page2"],
    "service_account": {...}
  }'
```

### Health Check
```bash
curl http://localhost:8080/api/health
```

## 🎯 Key Features

- ✅ No API key required
- ✅ Service account via request body
- ✅ Batch processing
- ✅ Auto caching
- ✅ Docker ready

## 🛠️ Commands

```bash
./deploy.sh start    # Start application
./deploy.sh stop     # Stop application  
./deploy.sh logs     # View logs
./deploy.sh restart  # Restart application
```

## 📊 Endpoints

- `GET /api/health` - Health check
- `POST /api/v1/index` - Submit single URL
- `POST /api/v1/index/batch` - Submit batch URLs
- `GET /api/v1/status/{url}` - Check URL status
- `GET /api/v1/cache/stats` - Cache statistics
- `POST /api/v1/cache/clear` - Clear cache

## 🔧 Requirements

- Docker & Docker Compose
- Google Service Account with Indexing API access

That's it! 🎉
