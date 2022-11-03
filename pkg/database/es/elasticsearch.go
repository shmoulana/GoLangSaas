package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/shmoulana/Redios/configs"
)

func NewClient(conf configs.Config) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: conf.ElasticUrls,
		Username:  conf.ElasticUser,
		Password:  conf.ElasticPassword,
	})
	if err != nil {
		return nil, err
	}

	return es, nil
}
