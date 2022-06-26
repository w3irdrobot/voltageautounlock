package main

import (
	"fmt"
	"os"
)

const (
	envAddress        = "ADDRESS"
	envNodeApi        = "VOLTAGE_NODE_API"
	envWebhookSecret  = "VOLTAGE_WEBHOOK_SECRET"
	envWalletPassword = "VOLTAGE_WALLET_PASSWORD"
)

type errMissingVariable struct {
	name string
}

func (e errMissingVariable) String() string {
	return fmt.Sprintf("missing environment variable %s", e.name)
}

func (e errMissingVariable) Error() string {
	return e.String()
}

type config struct {
	address        string
	nodeApi        string
	walletPassword string
	webhookSecret  string
}

func newConfig() (*config, error) {
	secret := os.Getenv(envWebhookSecret)
	if secret == "" {
		return nil, errMissingVariable{envWebhookSecret}
	}

	nodeApi := os.Getenv(envNodeApi)
	if nodeApi == "" {
		return nil, errMissingVariable{envNodeApi}
	}

	walletPassword := os.Getenv(envWalletPassword)
	if walletPassword == "" {
		return nil, errMissingVariable{envWalletPassword}
	}

	address := envVarWithDefault(envAddress, ":8080")

	return &config{
		address:        address,
		nodeApi:        nodeApi,
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
