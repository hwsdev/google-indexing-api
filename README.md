# Google Indexing API Application

🚀 REST API untuk mengintegrasikan Google Indexing API menggunakan Golang dan framework Gin untuk mempercepat proses indexing URL di Google Search.

## 📋 Features

- ✅ Submit single URL untuk indexing
- ✅ Batch processing untuk multiple URLs  
- ✅ **Service Account via Request Body** - Service account wajib melalui request body
- ✅ Service account caching untuk performa optimal
- ✅ **No Authentication Required** - Akses langsung tanpa API key
- ✅ URL status checking
- ✅ Comprehensive logging
- ✅ Docker support
- ✅ Health check endpoint
- ✅ CORS support
- ✅ Error handling
- ✅ Cache management endpoints

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Google API**: Indexing API v3
- **Logging**: Logrus
- **Validation**: Validator v10
- **Environment**: Godotenv

## 🚀 Quick Start

### Prerequisites

1. Go 1.21 atau lebih tinggi
2. Google Cloud Project dengan Indexing API enabled
3. Service Account dengan permission yang tepat

### Installation

1. Clone repository:
```bash
git clone <repository-url>
cd google-indexing-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Setup environment variables:
```bash
cp .env.example .env
# Edit .env file dengan konfigurasi Anda
```

4. Setup environment variables:
```bash
cp .env.example .env
# Edit .env file dengan konfigurasi server (opsional)
```

5. Run aplikasi:
```bash
go run cmd/main.go
```

**Note**: Tidak perlu setup service account di server - service account dikirim via request body.

## 🔧 Configuration

Edit file `.env` dengan konfigurasi berikut (semua opsional):

```env
# Server Configuration
PORT=8080
GIN_MODE=release

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Performance
MAX_BATCH_SIZE=100
CACHE_TTL_MINUTES=60
```

**Note**: Tidak ada konfigurasi service account atau API key yang diperlukan!

## 📚 API Endpoints

### Authentication

**Tidak ada autentikasi API key yang diperlukan!** 

Semua endpoint dapat diakses secara langsung tanpa header Authorization.

### Service Account

API ini **hanya mendukung dynamic service account** melalui request body. Service account credentials **wajib** disertakan di setiap request.

#### Format Service Account Request:
```json
{
  "url": "https://example.com/page",
  "service_account": {
    "type": "service_account",
    "project_id": "your-project-id",
    "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
    "client_email": "your-service@project.iam.gserviceaccount.com",
    "client_id": "your-client-id",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/..."
  }
}
```

### Endpoints

#### Health Check
```http
GET /api/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-09-14T10:30:00Z",
  "version": "1.0.0"
}
```

#### Submit Single URL
```http
POST /api/v1/index
Content-Type: application/json

{
  "url": "https://example.com/page",
  "service_account": {
    "type": "service_account",
    "project_id": "your-project-id",
    "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
    "client_email": "your-service@project.iam.gserviceaccount.com",
    "client_id": "your-client-id",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token"
  }
}
```

Response:
```json
{
  "success": true,
  "message": "URL submitted successfully",
  "url": "https://example.com/page"
}
```

#### Submit Batch URLs
```http
POST /api/v1/index/batch
Content-Type: application/json

{
  "urls": [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3"
  ],
  "service_account": {
    "type": "service_account",
    "project_id": "your-project-id",
    "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
    "client_email": "your-service@project.iam.gserviceaccount.com",
    "client_id": "your-client-id",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token"
  }
}
```

Response:
```json
{
  "success": true,
  "message": "Processed 3 URLs: 3 successful, 0 failed",
  "results": [
    {
      "success": true,
      "message": "URL submitted successfully",
      "url": "https://example.com/page1"
    }
  ],
  "statistics": {
    "total": 3,
    "successful": 3,
    "failed": 0
  }
}
```

#### Check URL Status
```http
GET /api/v1/status/https://example.com/page
```

Response:
```json
{
  "url": "https://example.com/page",
  "status": "URL_UPDATED",
  "last_updated": "2025-09-14T10:30:00Z"
}
```

#### Cache Management

**Get Cache Statistics**
```http
GET /api/v1/cache/stats
```

Response:
```json
{
  "cached_services": 3,
  "has_default": false,
  "timestamp": "2025-09-14T10:30:00Z"
}
```

**Clear Cache**
```http
POST /api/v1/cache/clear
```

Response:
```json
{
  "success": true,
  "message": "Cache cleared successfully"
}
```

## 🐳 Docker Deployment

### Build Image
```bash
docker build -t google-indexing-api .
```

### Run Container
```bash
docker run -d \
  --name google-indexing-api \
  -p 8080:8080 \
  -e PORT=8080 \
  -e API_KEY=your-secure-api-key \
  -e GOOGLE_PROJECT_ID=your-project-id \
  -v /path/to/service-account.json:/root/service-account.json \
  google-indexing-api
```

### Docker Compose (Recommended)

```yaml
version: '3.8'
services:
  google-indexing-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - API_KEY=your-secure-api-key
      - GOOGLE_PROJECT_ID=your-project-id
      - GOOGLE_APPLICATION_CREDENTIALS=/root/service-account.json
      - LOG_LEVEL=info
    volumes:
      - ./service-account.json:/root/service-account.json:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## 🔒 Security

- **No API Key Required** - Akses langsung tanpa autentikasi
- CORS configuration
- Input validation untuk semua requests
- Service account credentials validation
- HTTPS enforcement (recommended untuk production)
- Rate limiting per client
- Service account caching dengan auto-cleanup

## 🚀 Service Account Only Benefits

### Keuntungan:
1. **Simplified Access** - Tidak perlu manajemen API key
2. **Direct Integration** - Service account langsung di request body
3. **Multi-tenant Ready** - Setiap request bisa menggunakan service account berbeda
4. **Performance** - Service account di-cache untuk menghindari re-authentication
5. **Scalability** - Support untuk banyak klien dengan service account berbeda

### Use Cases:
- **Public API** - Dapat diakses langsung tanpa registrasi
- **SaaS Platform** - Setiap customer menggunakan service account sendiri
- **Agency** - Mengelola multiple client websites
- **Enterprise** - Multiple departments dengan project Google terpisah
- **Development** - Testing tanpa setup API key

## 📊 Monitoring & Logging

- Structured logging dengan Logrus
- Request/response logging
- Error tracking
- Performance metrics
- Health check endpoint untuk monitoring

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Test API dengan curl (tanpa API key)
curl -X POST http://localhost:8080/api/v1/index \
  -H "Content-Type: application/json" \
  -d '{
    "url":"https://example.com",
    "service_account": {
      "type": "service_account",
      "project_id": "your-project-id",
      "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
      "client_email": "your-service@project.iam.gserviceaccount.com",
      "client_id": "your-client-id",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token"
    }
  }'
```

## 🚀 Production Deployment

### Ubuntu Server dengan Docker

1. Install Docker:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

2. Clone dan deploy:
```bash
git clone <repository-url>
cd google-indexing-api
cp .env.example .env
# Edit .env dengan konfigurasi production
docker-compose up -d
```

3. Setup reverse proxy (Nginx):
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

MIT License

## 🆘 Troubleshooting

### Common Issues

1. **Authentication Error**: Pastikan service account JSON file valid dan memiliki permission yang tepat
2. **API Key Error**: Periksa API key di environment variable
3. **Google API Quota**: Monitor quota usage di Google Cloud Console
4. **Port Already in Use**: Ubah PORT di environment variable

### Logs

Check application logs:
```bash
# Docker logs
docker logs google-indexing-api

# File logs (jika dikonfigurasi)
tail -f app.log
```

## 📞 Support

Untuk bantuan, silakan buka issue di repository atau hubungi tim development.
