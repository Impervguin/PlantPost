package main

import (
	"PlantSite/internal/infra/opentelemetry"

	"github.com/spf13/viper"
)

const (
	TracingPrefix            = "tracing"
	TracingEnabledKey        = "enabled"
	TracingJaegerEndpointKey = "endpoint"
	TracingServiceNameKey    = "service_name"
	TracingServiceVersionKey = "service_version"
	TracingEnvironmentKey    = "environment"
)

func GetTracingConfig() *opentelemetry.TracingConfig {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}

	return &opentelemetry.TracingConfig{
		Enabled:        viper.GetBool(Key(TracingPrefix, TracingEnabledKey)),
		JaegerEndpoint: viper.GetString(Key(TracingPrefix, TracingJaegerEndpointKey)),
		ServiceName:    viper.GetString(Key(TracingPrefix, TracingServiceNameKey)),
		ServiceVersion: viper.GetString(Key(TracingPrefix, TracingServiceVersionKey)),
		Environment:    viper.GetString(Key(TracingPrefix, TracingEnvironmentKey)),
	}
}
