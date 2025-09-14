# Google Indexing API Application Plan

## 📋 Project Overview

Aplikasi untuk mengintegrasikan Google Indexing API menggunakan Golang dan framework Gin untuk mempercepat proses indexing URL di Google Search.

## 🎯 Objectives

- Membuat REST API untuk submit URL ke Google Indexing API
- Batch processing untuk multiple URLs
- Authentication dengan Service Account
- Monitoring dan logging

## 🏗️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Libraries**:
  - `google.golang.org/api/indexing/v3` - Google Indexing API
  - `github.com/gin-gonic/gin` - Web framework
  - `github.com/joho/godotenv` - Environment variables
  - `github.com/sirupsen/logrus` - Logging
  - `github.com/go-playground/validator/v10` - Validation

## 📁 Project Structure

```
google-indexing-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── indexing_handler.go
│   ├── services/
│   │   └── google_indexing_service.go
│   ├── middleware/
│   │   └── auth.go
│   └── models/
│       └── request.go
├── pkg/
│   └── utils/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## 🔧 Features

1. **Single URL Submission**

   - POST `/api/v1/index`
   - Request body: `{"url": "https://example.com/page"}`

2. **Batch URL Submission**

   - POST `/api/v1/index/batch`
   - Request body: `{"urls": ["url1", "url2", "url3"]}`

3. **URL Status Check**

   - GET `/api/v1/status/:url`

4. **Health Check**
   - GET `/api/health`

## 📝 Implementation Steps

### Phase 1: Setup & Configuration

- [x] Initialize Go module
- [x] Setup project structure
- [x] Create configuration management
- [x] Setup Google Service Account
- [x] Environment variables setup

### Phase 2: Core Development

- [x] Implement Google Indexing service
- [x] Create REST endpoints
- [x] Add request validation
- [x] Implement error handling
- [x] Add logging system

### Phase 3: Advanced Features

- [x] Batch processing with goroutines
- [x] Rate limiting
- [x] **Dynamic Service Account Support** - Service account melalui request body
- [x] **Service Account Caching** - Cache untuk performa optimal
- [x] Request queue system (melalui goroutines)
- [ ] Caching mechanism (untuk response data)
- [x] Metrics and monitoring

### Phase 4: Testing & Documentation

- [ ] Unit tests
- [ ] Integration tests
- [ ] API documentation (Swagger)
- [x] Deployment guide
- [x] Docker configuration

## 🔐 Security Considerations

- API key authentication
- Rate limiting per client
- Input validation
- CORS configuration
- HTTPS enforcement

## 📊 Monitoring & Logging

- Request/response logging
- Error tracking
- Performance metrics
- Google API quota monitoring

## 🚀 Deployment

- bikin script deploy sederhana di ubuntu docker. setingan ada di env
