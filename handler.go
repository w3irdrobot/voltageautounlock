package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	SecretHeader  = "Voltage-Secret"
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

type handler struct {
	ctx            context.Context
	secret         string
	nodeAPI        string
	walletPassword string
}

func newHandler(ctx context.Context, secret, nodeAPI, walletPassword string) *handler {
	return &handler{ctx, secret, nodeAPI, walletPassword}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	h.handleEvent(r.Context(), w, event)
}

func (h *handler) handleEvent(ctx context.Context, w http.ResponseWriter, event unlockEvent) {
	if event.API != h.nodeAPI {
		http.Error(w, fmt.Sprintf("api '%s' does not match expected api", event.API), http.StatusBadRequest)

		return
	}

	if event.Kind != Status {
		http.Error(w, fmt.Sprintf("event type '%s' does not match expected type", event.Kind), http.StatusBadRequest)

		return
	}

	if event.Details.Status != WaitingUnlock {
		http.Error(
			w,
			fmt.Sprintf("event details status '%s' does not match expected status", event.Details.Status),
			http.StatusBadRequest)

		return
	}

	body, err := json.Marshal(lndUnlock{
		WalletPassword: base64.StdEncoding.EncodeToString([]byte(h.walletPassword)),
		StatelessInit:  true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshalling the request to lnd: %s", err), http.StatusInternalServerError)

		return
	}

	endpoint := url.URL{Scheme: "https", Host: h.nodeAPI + ":8080", Path: LNDUnlockPath}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the request to lnd: %s", err), http.StatusInternalServerError)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("error sending the request to lnd: %s", err), http.StatusInternalServerError)

		return
	}
	defer res.Body.Close()

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		body, _ := io.ReadAll(res.Body)
		http.Error(w, fmt.Sprintf("error returned from unlocking LND: %s", body), http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
