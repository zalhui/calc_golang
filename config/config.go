package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TimeAddition       time.Duration
	TimeSubtraction    time.Duration
	TimeMultiplication time.Duration
	TimeDivision       time.Duration
	ComputingPower     int
	JWTSecret          string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		TimeAddition:       getEnvDuration("TIME_ADDITION_MS", 1000),
		TimeSubtraction:    getEnvDuration("TIME_SUBTRACTION_MS", 1000),
		TimeMultiplication: getEnvDuration("TIME_MULTIPLICATIONS_MS", 1000),
		TimeDivision:       getEnvDuration("TIME_DIVISIONS_MS", 1000),
		ComputingPower:     getEnvInt("COMPUTING_POWER", 1),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return time.Duration(value) * time.Millisecond
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
