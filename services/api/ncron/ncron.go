package ncron

import (
	"log"
	"ncronus/services/types"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	CST *cron.Cron
	IST *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		CST: cron.New(cron.WithLocation(types.CST_TIMEZONE)),
		IST: cron.New(cron.WithLocation(types.IST_TIMEZONE)),
	}
}

func (c *Cron) StartCron() {
	c.CST.Start()
	c.IST.Start()
	log.Println("cron started")
}

func (c *Cron) StopCron() {
	<-c.CST.Stop().Done()
	<-c.IST.Stop().Done()
	log.Println("cron stopped")
}
