package main

import (
	"fmt"
	"os"
)

const (
	envAddress = "ADDRESS"
	envNodeAPI = "VOLTAGE_NODE_API"
	//nolint:gosec
	envWebhookSecret  = "VOLTAGE_WEBHOOK_SECRET"
	envWalletPassword = "VOLTAGE_WALLET_PASSWORD"
)

type missingVariableError struct {
	name string
}

func (e missingVariableError) String() string {
	return fmt.Sprintf("missing environment variable %s", e.name)
}

func (e missingVariableError) Error() string {
	return e.String()
}

type config struct {
	address        string
	nodeAPI        string
	walletPassword string
	webhookSecret  string
}

func newConfig() (*config, error) {
	secret := os.Getenv(envWebhookSecret)
	if secret == "" {
		return nil, missingVariableError{envWebhookSecret}
	}

	nodeAPI := os.Getenv(envNodeAPI)
	if nodeAPI == "" {
		return nil, missingVariableError{envNodeAPI}
	}

	walletPassword := os.Getenv(envWalletPassword)
	if walletPassword == "" {
		return nil, missingVariableError{envWalletPassword}
	}

	address := envVarWithDefault(envAddress, ":8080")

	return &config{
		address:        address,
		nodeAPI:        nodeAPI,
		walletPassword: walletPassword,
		webhookSecret:  secret,
	}, nil
}

func envVarWithDefault(name, fallback string) string {
	if value, found := os.LookupEnv(name); found {
		return value
	}

	return fallback
}
