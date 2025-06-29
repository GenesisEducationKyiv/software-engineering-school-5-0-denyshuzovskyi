package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
)

type SubscriptionService interface {
	Subscribe(context.Context, dto.SubscriptionRequest) error
	Confirm(context.Context, string) error
	Unsubscribe(context.Context, string) error
}

type SubscriptionHandler struct {
	subscriptionService SubscriptionService
	log                 *slog.Logger
}

func NewSubscriptionHandler(subscriptionService SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		log:                 log,
	}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		h.log.Error("error parsing form", "error", err)
		return
	}

	var subscriptionReq dto.SubscriptionRequest
	subscriptionReq.Email = r.FormValue("email")
	subscriptionReq.City = r.FormValue("city")
	subscriptionReq.Frequency = r.FormValue("frequency")

	if err = h.subscriptionService.Subscribe(r.Context(), subscriptionReq); err != nil {
		if errors.Is(err, commonerrors.ErrValidationFailed) {
			http.Error(w, "invalid input", http.StatusBadRequest)
			h.log.Error("error validating data", "error", err)
			return
		}
		if errors.Is(err, commonerrors.ErrLocationNotFound) {
			http.Error(w, "invalid input", http.StatusBadRequest)
			h.log.Error("couldn't validate city", "error", err)
			return
		}
		if errors.Is(err, commonerrors.ErrSubscriptionAlreadyExists) {
			http.Error(w, "email already subscribed", http.StatusConflict)
			h.log.Error("subscription already exists", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error making subscription", "error", err)
		return
	}
}

func (h *SubscriptionHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	if err := h.subscriptionService.Confirm(r.Context(), token); err != nil {
		if errors.Is(err, commonerrors.ErrInvalidToken) {
			http.Error(w, "invalid token", http.StatusBadRequest)
			h.log.Error("invalid token", "error", err)
			return
		} else if errors.Is(err, commonerrors.ErrTokenNotFound) {
			http.Error(w, "token not found", http.StatusNotFound)
			h.log.Error("token not found", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error confirming subscription", "error", err)
		return
	}
}

func (h *SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	if err := h.subscriptionService.Unsubscribe(r.Context(), token); err != nil {
		if errors.Is(err, commonerrors.ErrInvalidToken) {
			http.Error(w, "invalid token", http.StatusBadRequest)
			h.log.Error("invalid token", "error", err)
			return
		} else if errors.Is(err, commonerrors.ErrTokenNotFound) {
			http.Error(w, "token not found", http.StatusNotFound)
			h.log.Error("token not found", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error confirming subscription", "error", err)
		return
	}
}
