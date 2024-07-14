package cron_job

import (
	cronn "github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func CronJob() error {
	parser := cronn.NewParser(cronn.Minute | cronn.Hour | cronn.Dom | cronn.Month | cronn.Dow)

	parserOpts := cronn.WithParser(parser)

	job := cronn.New(parserOpts)
	_, err := job.AddFunc("*/2 * * * *", func() {
		myJob()
	})

	job.Start()
	if err != nil {
		return err
	}
	return nil
}

func myJob() {

	log.Info().Msg("Cron Job Running every 2 minutes")
}
