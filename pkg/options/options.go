package options

import (
	"strings"

	"github.com/spf13/pflag"
)

func FlagNameFunc(fs *pflag.FlagSet, name string) string {
	return strings.Join([]string{fs.Name(), name}, ".")
}
