package kratosutils

import (
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

func NewLogger() log.Logger {
	return log.NewStdLogger(os.Stdout)
}
