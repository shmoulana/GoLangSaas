package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/pkg/utils"
	"github.com/sirupsen/logrus"
)

type LoggerService struct {
	ElasticClient *elasticsearch.Client
}

type logger struct {
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	Meta          interface{} `json:"meta"`
	Event         string      `json:"event"`
	Message       string      `json:"message"`
	CreatedAtUnix int64       `json:"created_at_unix"`
	CreatedAt     string      `json:"created_at"`
}

type loggerInterface interface {
	Name() string
	Type() string
}

func (s LoggerService) Info(data interface{}, msg string) {
	now := time.Now()

	var logData logger = logger{
		Status:        constant.LogStatusInfo,
		Type:          data.(loggerInterface).Type(),
		Event:         data.(loggerInterface).Name(),
		Meta:          nil,
		Message:       msg,
		CreatedAt:     utils.DateToString(now),
		CreatedAtUnix: now.Unix(),
	}

	logrus.Info(s.convertLogrus(logData))
	s.saveLog(logData)
	return
}

func (s LoggerService) Error(data interface{}, err error, msg string) {
	now := time.Now()

	var logData logger = logger{
		Status:        constant.LogStatusInfo,
		Type:          data.(loggerInterface).Type(),
		Event:         data.(loggerInterface).Name(),
		Meta:          err,
		Message:       msg,
		CreatedAt:     utils.DateToString(now),
		CreatedAtUnix: now.Unix(),
	}

	logrus.Error(s.convertLogrus(logData))
	s.saveLog(logData)
	return
}

// convert logger struct to readable message
func (s LoggerService) convertLogrus(val logger) string {
	return fmt.Sprintf("status=%s type=%s event=%s message=%s meta=%s", val.Status, val.Type, val.Event, val.Message, val.Meta)
}

// save log to elasticsearch
func (s LoggerService) saveLog(val logger) {
	res, err := s.ElasticClient.Ping()
	if err != nil {
		logrus.Error(err)
		return
	}

	if res.StatusCode != http.StatusOK {
		return
	}

	res, err = s.ElasticClient.Index("log", esutil.NewJSONReader(val))
	if err != nil {
		logrus.Error(err)
		return
	}

	return
}
