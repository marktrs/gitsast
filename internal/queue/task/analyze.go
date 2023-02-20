package task

import (
	"time"

	"github.com/labstack/gommon/log"
)

type AnalyzeTask struct {
}

func NewAnalyzeTask() *AnalyzeTask {
	return &AnalyzeTask{}
}

func (t *AnalyzeTask) Start(reportID string) {
	log.Infof("starting analyzed task reportID=%s", reportID)
	time.Sleep(3 * time.Second)
	// 1. look up for rules

	// 2. look up for latest code

	// 3. scan for issue

	// 4. add issue to report (db)

	// 5. update report status (db)

	// 6. end task
	log.Info("analyzed task completed")
}
