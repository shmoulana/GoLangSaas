package router

const (
	PingPath         = "/v1/ping"
	TenantPath       = "/v1/tenants"
	TenantWithIdPath = "/v1/tenants/:id"
	UserSignInPath   = "/v1/user/sign-in"
	UserSignUpPath   = "/v1/user/sign-up"

	// Temporary just for testing
	EmailTestNowPath   = "/v1/email/test/now"
	EmailTestQueuePath = "/v1/email/test/queue"

	EmailNowPath   = "/v1/email/now"
	EmailQueuePath = "/v1/email/queue"
)
