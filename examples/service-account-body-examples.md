# Contoh Request dengan Service Account di Body

## 1. Submit Single URL dengan Service Account

```bash
curl -X POST http://localhost:8080/api/v1/index \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "url": "https://example.com/page",
    "service_account": {
      "type": "service_account",
      "project_id": "your-project-id",
      "private_key_id": "your-private-key-id",
      "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_HERE\n-----END PRIVATE KEY-----\n",
      "client_email": "your-service@your-project.iam.gserviceaccount.com",
      "client_id": "your-client-id",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token",
      "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
      "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service%40your-project.iam.gserviceaccount.com"
    }
  }'
```

## 2. Submit Single URL dengan Default Service Account

```bash
curl -X POST http://localhost:8080/api/v1/index \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "url": "https://example.com/page",
    "use_default_account": true
  }'
```

## 3. Submit Batch URLs dengan Service Account

```bash
curl -X POST http://localhost:8080/api/v1/index/batch \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "urls": [
      "https://example.com/page1",
      "https://example.com/page2",
      "https://example.com/page3"
    ],
    "service_account": {
      "type": "service_account",
      "project_id": "your-project-id",
      "private_key_id": "your-private-key-id",
      "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_HERE\n-----END PRIVATE KEY-----\n",
      "client_email": "your-service@your-project.iam.gserviceaccount.com",
      "client_id": "your-client-id",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token",
      "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
      "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service%40your-project.iam.gserviceaccount.com"
    }
  }'
```

## 4. Get Cache Statistics

```bash
curl -X GET http://localhost:8080/api/v1/cache/stats \
  -H "Authorization: Bearer your-api-key"
```

Response:

```json
{
  "cached_services": 2,
  "has_default": true,
  "timestamp": "2025-09-14T10:30:00Z"
}
```

## 5. Clear Cache

```bash
curl -X POST http://localhost:8080/api/v1/cache/clear \
  -H "Authorization: Bearer your-api-key"
```

Response:

```json
{
  "success": true,
  "message": "Cache cleared successfully"
}
```

## Format Response

### Success Response (Single URL):

```json
{
  "success": true,
  "message": "URL submitted successfully",
  "url": "https://example.com/page"
}
```

### Success Response (Batch URLs):

```json
{
  "success": true,
  "message": "Processed 3 URLs: 3 successful, 0 failed",
  "results": [
    {
      "success": true,
      "message": "URL submitted successfully",
      "url": "https://example.com/page1"
    },
    {
      "success": true,
      "message": "URL submitted successfully",
      "url": "https://example.com/page2"
    },
    {
      "success": true,
      "message": "URL submitted successfully",
      "url": "https://example.com/page3"
    }
  ],
  "statistics": {
    "total": 3,
    "successful": 3,
    "failed": 0
  }
}
```

### Error Response:

```json
{
  "error": "Bad Request",
  "message": "Invalid service account: project_id is required",
  "code": 400
}
```

## Keuntungan Menggunakan Service Account di Body:

1. **Fleksibilitas**: Dapat menggunakan multiple service accounts tanpa restart
2. **Dynamic**: Setiap request bisa menggunakan service account yang berbeda
3. **Cache**: Service accounts di-cache untuk performa yang lebih baik
4. **Fallback**: Tetap bisa menggunakan default service account jika tidak ada di request
5. **Security**: Credentials hanya ada di memory sementara

## Security Notes:

- Pastikan menggunakan HTTPS di production
- Service account credentials di-cache untuk performa, clear cache secara berkala
- Validasi service account dilakukan sebelum digunakan
- Logs tidak mencatat private key untuk keamanan
