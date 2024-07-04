package version

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type Build struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const hello = `Build version: %s
Build date: %s
Build commit: %s
`

func (b Build) Print(w io.Writer) {
	p := []any{
		val(b.BuildVersion),
		val(b.BuildDate),
		val(b.BuildCommit),
	}
	if _, err := fmt.Fprintf(w, hello, p...); err != nil {
		logrus.WithError(err).Error("Fail to print the build version")
	}
}

func val(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}