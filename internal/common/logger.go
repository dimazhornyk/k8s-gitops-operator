package common

import "github.com/sirupsen/logrus"

func NewLogger(level logrus.Level) (*logrus.Entry, error) {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
		},
	})

	logger.SetLevel(level)

	return logrus.NewEntry(logger), nil
}
