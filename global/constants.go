package global

const (
	ApiPrefix     = "/api/v1"
	WebhookPrefix = "/webhook"

	// DefaultDecimalPrecision is used in different reporting to the FE.
	// Rounding floating point response would make responses more lightweight
	DefaultDecimalPrecision int = 2

	DateColumnPostfix = "_date"
	DateColumnFormat  = "DD Mon, YYYY"

	UnavailableVersion = "unavailable"

	RedisKeyPrefix = "cliplab:"
	RedisPostQueue = RedisKeyPrefix + "post_queue"
)
