package builder

const (
	// Keywords used for built in variables
	LAYOUT  = "layout"
	PAGE    = "page"
	CONTENT = "content"
	SITE    = "site"

	// mode for all files created, bandaid due to permission problems
	DIR_MODE = 0777

	// Regex used to separate vars from content on a page.
	// these are awfully long names..
	VAR_CONTENT_DIVIDER = "^[-]{3,}$"
	NAME_VAL_DIVIDER    = ":"

	// Web serving constants
	WEB_ROOT         = "/"
	WEB_DEFAULT_ROOT = "index.html"
)
