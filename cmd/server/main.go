package main

import (
	"flag"
	"github.com/newpurr/stock/pkg/utils/commonutils"
	_ "go.uber.org/automaxprocs"
	"os"
)

var (
	// Name is the name of the compiled software.
	// go build -ldflags "-X main.Version=x.y.z"
	Name string

	// Version is the version of the compiled software.
	Version string

	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
	Env   = commonutils.Env("ENV", "local")
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/"+Env, "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	newAppStarter(
		Name,
		Version,
		flagconf,
		id,
	).start()
}
