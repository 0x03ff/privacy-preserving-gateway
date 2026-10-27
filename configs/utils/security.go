package configutils

import "regexp"

// Regex to allow only alphanumeric, dots, underscores, and dashes
var safeNameRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)


