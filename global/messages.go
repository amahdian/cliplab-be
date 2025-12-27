package global

type Messages string

const (
	InvalidUrl      Messages = "Invalid URL. please send me a valid URL."
	NotSupportedUrl Messages = "%s is not supported yet. We'll add support for this social media soon."
)
