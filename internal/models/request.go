package models

type ServiceAccountCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

type IndexRequest struct {
	URL            string                     `json:"url" validate:"required,url" binding:"required"`
	ServiceAccount *ServiceAccountCredentials `json:"service_account" validate:"required" binding:"required"`
}

type BatchIndexRequest struct {
	URLs           []string                   `json:"urls" validate:"required,min=1,dive,url" binding:"required"`
	ServiceAccount *ServiceAccountCredentials `json:"service_account" validate:"required" binding:"required"`
}

type IndexResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}

type BatchIndexResponse struct {
	Success    bool                      `json:"success"`
	Message    string                    `json:"message"`
	Results    []IndexResponse           `json:"results,omitempty"`
	Statistics BatchIndexResponseStats   `json:"statistics,omitempty"`
}

type BatchIndexResponseStats struct {
	Total     int `json:"total"`
	Successful int `json:"successful"`
	Failed    int `json:"failed"`
}

type StatusResponse struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	LastUpdated string `json:"last_updated,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
