package main

import (
	"fmt"
	"github.com/franchb/cli"
	"github.com/franchb/go-aws-s3-static-site-push/actions/aws/s3/push"
	"github.com/franchb/go-aws-s3-static-site-push/log"
	"os"
)

type rootT struct {
	cli.Helper
	LogLevel string `cli:"log-level" usage:"Replaces the prefix with a user one."`
}

// Validate implements cli.Validator interface
func (r *rootT) Validate(ctx *cli.Context) error {
	return nil
}

func main() {
	code := cli.Run(new(rootT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)
		logLevel, _ := log.ParseLevel(argv.LogLevel)
		log.InitLogs(logLevel)
		action := push.NewS3PushAction()
		if err := action.Open(); err != nil {
			return fmt.Errorf("failed to open %s handler\nError=%s",
				ctx.Color().Cyan(action.Name()), ctx.Color().RedBg(err))
		}
		if err := action.Do(); err != nil {
			return fmt.Errorf("failed to run action %s \nError=%s",
				ctx.Color().Cyan(action.Name()), ctx.Color().RedBg(err))
		}
		if err := action.Close(); err != nil {
			return fmt.Errorf("failed to run action %s \nError=%s",
				ctx.Color().Cyan(action.Name()), ctx.Color().RedBg(err))
		}
		return nil
	})
	os.Exit(code)
}
