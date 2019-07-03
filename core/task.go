package core

import (
	"github.com/jakecoffman/cron"
)

var (
	Cron *cron.Cron
)

func InitCron() {
	Cron = cron.New()
	Cron.Start()
}

