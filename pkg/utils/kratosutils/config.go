package kratosutils

import (
	"github.com/go-kratos/kratos/v2/config"
)

func NewConfLoader(cs config.Source) (confLoader config.Config, cleanup func(), err error) {
	c := config.New(
		config.WithSource(cs),
	)
	cleanup = func() {
		_ = c.Close()
	}

	err = c.Load()

	return c, cleanup, err
}
