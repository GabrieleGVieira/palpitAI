package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/gabrielevieira/palpitai/backend/internal/usecase"
)

func ListGroupsHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		groupList, err := groups.ListGroups(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Não foi possivel listar os grupos.")
			return
		}

		writeJSON(w, http.StatusOK, map[string][]dto.GroupListItemResponse{
			"groups": groupList,
		})
	}
}

func CreateGroupHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, displayName, err := userIDAndDisplayNameFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		var request dto.CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "JSON invalido.")
			return
		}

		normalizedRequest, err := usecase.NormalizeCreateGroupRequest(request)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		group, err := groups.CreateGroup(r.Context(), userID, displayName, normalizedRequest)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Não foi possivel criar o grupo.")
			return
		}

		writeJSON(w, http.StatusCreated, group)
	}
}

func UpdateGroupHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, err := userIDFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		var request dto.UpdateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "JSON invalido.")
			return
		}

		normalizedRequest, err := usecase.NormalizeUpdateGroupRequest(request)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		group, err := groups.UpdateGroup(r.Context(), ownerID, r.PathValue("groupID"), normalizedRequest)
		if err != nil {
			if apperrors.IsNotFound(err) {
				writeError(w, http.StatusNotFound, "Grupo não encontrado.")
				return
			}

			writeError(w, http.StatusInternalServerError, "Não foi possivel atualizar o grupo.")
			return
		}

		writeJSON(w, http.StatusOK, group)
	}
}

func JoinGroupHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, displayName, err := userIDAndDisplayNameFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		var request dto.JoinGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "JSON invalido.")
			return
		}

		inviteCode := usecase.NormalizeInviteCode(request.InviteCode)
		if inviteCode == "" {
			writeError(w, http.StatusBadRequest, "Informe o codigo do grupo.")
			return
		}

		response, err := groups.JoinGroup(r.Context(), userID, displayName, inviteCode)
		if err != nil {
			switch {
			case apperrors.IsNotFound(err):
				writeError(w, http.StatusNotFound, "Grupo não encontrado.")
			case apperrors.IsConflict(err):
				writeError(w, http.StatusConflict, "Este grupo atingiu o limite de participantes.")
			default:
				writeError(w, http.StatusInternalServerError, "Não foi possivel entrar no grupo.")
			}
			return
		}

		writeJSON(w, http.StatusOK, response)
	}
}

func ListJoinRequestsHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		requests, err := groups.ListJoinRequests(r.Context(), userID, r.PathValue("groupID"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Não foi possivel listar as solicitacoes.")
			return
		}

		writeJSON(w, http.StatusOK, map[string][]dto.JoinRequestResponse{
			"requests": requests,
		})
	}
}

func ApproveJoinRequestHandler(cfg config.Config, groups usecase.GroupUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, err := userIDFromRequest(r, cfg)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		if err := groups.ApproveJoinRequest(r.Context(), ownerID, r.PathValue("groupID"), r.PathValue("userID")); err != nil {
			switch {
			case apperrors.IsNotFound(err):
				writeError(w, http.StatusNotFound, "Solicitacao não encontrada.")
			case apperrors.IsConflict(err):
				writeError(w, http.StatusConflict, "Este grupo atingiu o limite de participantes.")
			default:
				writeError(w, http.StatusInternalServerError, "Não foi possivel aprovar a solicitacao.")
			}
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "approved",
		})
	}
}
