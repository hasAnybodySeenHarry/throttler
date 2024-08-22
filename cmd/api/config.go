package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	port     int
	env      string
	redis    rd
	grpcAddr string
}

type rd struct {
	addr     string
	password string
}

func loadConfig(cfg *config) {
	flag.StringVar(&cfg.grpcAddr, "grpcAddr", os.Getenv("GRPC-ADDR"), "The address of the gRPC server")

	flag.IntVar(&cfg.port, "port", getEnvInt("PORT", 8080), "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "The environment of the server")

	flag.StringVar(&cfg.redis.addr, "addr", os.Getenv("ADDR"), "The environment of the server")
	flag.StringVar(&cfg.redis.password, "password", os.Getenv("PASSWORD"), "The environment of the server")

	flag.Parse()
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Invalid value for environment variable %s: %s\n", key, valueStr)
		return defaultValue
	}

	return value
}
