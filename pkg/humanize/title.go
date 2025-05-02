package humanize

import "github.com/lonelyshoeh/api/pkg/utils"

// Title formats the input title into a consistent format.
// Current operations is just limiting the string to 60 characters
func Title(title string) string {
	const MaxLength = 60

	return utils.TruncateString(title, MaxLength)
}
