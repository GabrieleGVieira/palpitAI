package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/gabrielevieira/palpitai/backend/internal/config"
)

type AccountDeletionUsecase interface {
	DeleteAccount(ctx context.Context, userID string) error
}

func DeleteAccountHandler(cfg config.Config, accounts AccountDeletionUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r, cfg)
		if err != nil {
			fmt.Printf("account deletion failed: %v\n", err)
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		if err := accounts.DeleteAccount(r.Context(), userID); err != nil {
			if apperrors.IsConflict(err) {
				writeError(w, http.StatusConflict, "Transfira a propriedade dos grupos que você administra antes de excluir sua conta.")
				return
			}

			fmt.Printf("account deletion failed: %v\n", err)
			slog.Error("account deletion failed", "error", err)
			writeError(w, http.StatusInternalServerError, "Não foi possível excluir a conta agora.")
			return
		}

		if err := deleteSupabaseAuthUser(r, cfg, userID); err != nil {
			fmt.Printf("account deletion failed: %v\n", err)
			slog.Error("supabase auth user deletion failed", "error", err)
		}

		slog.Info("account deletion processed")
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Conta marcada para exclusão e dados pessoais anonimizados.",
		})
	}
}

func deleteSupabaseAuthUser(r *http.Request, cfg config.Config, userID string) error {
	if strings.TrimSpace(cfg.SupabaseURL) == "" || strings.TrimSpace(cfg.SupabaseServiceRoleKey) == "" {
		return nil
	}

	endpoint, err := url.JoinPath(cfg.SupabaseURL, "/auth/v1/admin/users", userID)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(r.Context(), http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)
	request.Header.Set("apikey", cfg.SupabaseServiceRoleKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New("supabase auth delete failed")
	}

	return nil
}
