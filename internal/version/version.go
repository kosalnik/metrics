package version

import (
	"fmt"
	"io"

	"github.com/kosalnik/metrics/internal/log"
)

type VersionInfo struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const hello = `Build version: %s
Build date: %s
Build commit: %s
`

func (b VersionInfo) Print(w io.Writer) {
	p := []any{
		val(b.BuildVersion),
		val(b.BuildDate),
		val(b.BuildCommit),
	}
	if _, err := fmt.Fprintf(w, hello, p...); err != nil {
		log.Error().Err(err).Msg("Fail to print the build version")
	}
}

func val(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}
