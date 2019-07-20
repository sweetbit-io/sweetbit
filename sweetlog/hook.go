package sweetlog

import "github.com/sirupsen/logrus"

type SweetLog struct {
}

func New() *SweetLog {
	return &SweetLog{}
}

func (h *SweetLog) Fire(entry *logrus.Entry) error {
	// TODO(davidknezic) collect logs here
	return nil
}

func (h *SweetLog) Levels() []logrus.Level {
	return logrus.AllLevels
}
