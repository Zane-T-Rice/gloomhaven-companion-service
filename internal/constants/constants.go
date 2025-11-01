package constants

const (
	// Error Messages and Status Codes
	STATUS_CODE_FORBIDDEN   = 403
	STATUS_CODE_NOT_FOUND   = 404
	NOT_FOUND_ERROR_MESSAGE = "The requested resource was not found."
	FORBIDDEN_ERROR_MESSAGE = "Insufficient scope."

	// Scope Names
	SCOPE_READ_CAMPAIGNS  = "read:campaigns"
	SCOPE_WRITE_CAMPAIGNS = "write:campaigns"

	// Environment Variable Names
	AUDIENCE                         = "AUDIENCE"
	ISSUER                           = "ISSUER"
	GLOOMHAVEN_COMPANION_SERVICE_URL = "GLOOMHAVEN_COMPANION_SERVICE_URL"
	WEBSITE_DOMAIN                   = "WEBSITE_DOMAIN"
	API_GATEWAY_BASE_PATH            = "API_GATEWAY_BASE_PATH"
	LOCAL_SERVICE_PORT               = "LOCAL_SERVICE_PORT"
	LOCAL_DATABASE_ENDPOINT          = "LOCAL_DATABASE_ENDPOINT"

	// Secret Names
	AUDIENCE_SECRET_NAME              = "gloomhaven-companion-service-audience"
	ISSUER_SECRET_NAME                = "gloomhaven-companion-service-issuer"
	URL_SECRET_NAME                   = "gloomhaven-companion-service-url"
	WEBSITE_DOMAIN_SECRET_NAME        = "gloomhaven-companion-service-website-domain"
	API_GATEWAY_BASE_PATH_SECRET_NAME = "gloomhaven-companion-service-api-gateway-base-path"

	// Route Names
	CAMPAIGNS = "campaigns"

	// DynamoDB
	TABLE_NAME = "gloomhaven-companion-service"
	PARENT     = "parent"
	ENTITY     = "entity"
	SEPERATOR  = "#"
	ROOT       = SEPERATOR + "ROOT"
	CAMPAIGN   = SEPERATOR + "CAMPAIGN"
	PLAYER     = SEPERATOR + "PLAYER"
)

// This one is defined outside the const block to allow its address to be taken
var ENTITY_INDEX = "entity-index"
