package ender

import (
	"openreplay/backend/pkg/env"
)

type Config struct {
	GroupEnder       string
	TopicTrigger     string
	LoggerTimeout    int
	TopicRawWeb      string
	TopicRawIOS      string
	ProducerTimeout  int
	PartitionsNumber int
}

func New() *Config {
	return &Config{
		GroupEnder:       env.String("GROUP_ENDER"),
		TopicTrigger:     env.String("TOPIC_TRIGGER"),
		LoggerTimeout:    env.Int("LOG_QUEUE_STATS_INTERVAL_SEC"),
		TopicRawWeb:      env.String("TOPIC_RAW_WEB"),
		TopicRawIOS:      env.String("TOPIC_RAW_IOS"),
		ProducerTimeout:  2000,
		PartitionsNumber: env.Int("PARTITIONS_NUMBER"),
	}
}
