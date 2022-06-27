package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	SecretHeader  = "VOLTAGE_SECRET"
	Status        = "status"
	WaitingUnlock = "waiting_unlock"
	LNDUnlockPath = "/v1/unlockwallet"
)

type unlockEvent struct {
	API     string        `json:"api"`
	Kind    string        `json:"type"`
	Details unlockDetails `json:"details"`
}

type unlockDetails struct {
	Status string `json:"status"`
}

type lndUnlock struct {
	WalletPassword string `json:"wallet_password"`
	StatelessInit  bool   `json:"stateless_init"`
}

type unlockHandler struct {
	ctx            context.Context
	secret         string
	nodeAPI        string
	walletPassword string
}

func newUnlockHandler(ctx context.Context, secret, nodeAPI, walletPassword string) *unlockHandler {
	return &unlockHandler{ctx, secret, nodeAPI, walletPassword}
}

func (h *unlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if value := r.Header.Get(SecretHeader); value != h.secret {
		log.Printf("secret '%s' from header didn't match expected secret\n", value)
		w.WriteHeader(http.StatusForbidden)

		return
	}

	var event unlockEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println("error reading the request body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.handleEvent(w, event)
}

func (h *unlockHandler) handleEvent(w http.ResponseWriter, event unlockEvent) {
	if event.API != h.nodeAPI {
		log.Printf("api '%s' does not match expected api\n", event.API)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if event.Kind != Status {
		log.Printf("event type '%s' does not match expected type\n", event.Kind)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if event.Details.Status != WaitingUnlock {
		log.Printf("event details status '%s' does not match expected status\n", event.Details.Status)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	body, err := json.Marshal(lndUnlock{
		WalletPassword: h.walletPassword,
		StatelessInit:  true,
	})
	if err != nil {
		log.Printf("error marshalling the request to lnd")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	endpoint := url.URL{
		Scheme: "https",
		Host:   h.nodeAPI + ":8080",
		Path:   LNDUnlockPath,
	}

	res, err := http.Post(endpoint.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error sending the request to lnd")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		body, _ := io.ReadAll(res.Body)
		log.Printf("error returned from unlocking LND: %s", body)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
