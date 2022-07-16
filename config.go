package main

import (
	"fmt"
	"os"
)

const (
	envAddress = "ADDRESS"
	envPort    = "PORT"
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

	address := getAddressFromEnv()

	return &config{
		address:        address,
		nodeAPI:        nodeAPI,
		walletPassword: walletPassword,
		webhookSecret:  secret,
	}, nil
}

func getAddressFromEnv() string {
	if address := os.Getenv(envAddress); address != "" {
		return address
	}

	if port := os.Getenv(envPort); port != "" {
		return ":" + port
	}

	return ":8080"
}
