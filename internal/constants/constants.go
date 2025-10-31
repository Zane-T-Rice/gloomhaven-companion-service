package constants

const (
	// Error Messages and Status Codes
	STATUS_CODE_FORBIDDEN   = 403
	FORBIDDEN_ERROR_MESSAGE = "Insufficient scope."

	// Scope Names
	SCOPE_READ_CAMPAIGNS  = "read:campaigns"
	SCOPE_WRITE_CAMPAIGNS = "write:campaigns"

	// Environment Variable Names
	AUDIENCE                         = "AUDIENCE"
	ISSUER                           = "ISSUER"
	WEBSITE_DOMAIN                   = "WEBSITE_DOMAIN"
	LOCAL_SERVICE_PORT               = "LOCAL_SERVICE_PORT"
	API_GATEWAY_BASE_PATH            = "API_GATEWAY_BASE_PATH"
	LOCAL_DATABASE_ENDPOINT          = "LOCAL_DATABASE_ENDPOINT"
	GLOOMHAVEN_COMPANION_SERVICE_URL = "GLOOMHAVEN_COMPANION_SERVICE_URL"

	// Route Names
	CAMPAIGNS = "campaigns"

	// Secret Names
	AUDIENCE_SECRET_NAME              = "gloomhaven-companion-service-audience"
	ISSUER_SECRET_NAME                = "gloomhaven-companion-service-issuer"
	URL_SECRET_NAME                   = "gloomhaven-companion-service-url"
	WEBSITE_DOMAIN_SECRET_NAME        = "gloomhaven-companion-service-website-domain"
	API_GATEWAY_BASE_PATH_SECRET_NAME = "gloomhaven-companion-service-api-gateway-base-path"
)
