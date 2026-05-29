package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"strings"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

const inviteCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var (
	ErrGroupFull          = apperrors.NewConflict("group is full")
	ErrGroupNotFound      = apperrors.NewNotFound("group not found")
	ErrGroupOwnerRequired = apperrors.NewForbidden("group owner required")
)

type GroupUsecase struct {
	db Datastore
}

func NewGroupUsecase(db Datastore) GroupUsecase {
	return GroupUsecase{db: db}
}

func (uc GroupUsecase) ListGroups(ctx context.Context, userID string) ([]dto.GroupListItemResponse, error) {
	return ListGroups(ctx, uc.db, userID)
}

func (uc GroupUsecase) CreateGroup(ctx context.Context, userID string, displayName string, request dto.CreateGroupRequest) (dto.GroupResponse, error) {
	return CreateGroup(ctx, uc.db, userID, displayName, request)
}

func (uc GroupUsecase) UpdateGroup(ctx context.Context, ownerID string, groupID string, request dto.UpdateGroupRequest) (dto.GroupResponse, error) {
	return UpdateGroup(ctx, uc.db, ownerID, groupID, request)
}

func (uc GroupUsecase) JoinGroup(ctx context.Context, userID string, displayName string, inviteCode string) (dto.JoinGroupResponse, error) {
	return JoinGroup(ctx, uc.db, userID, displayName, inviteCode)
}

func (uc GroupUsecase) ListJoinRequests(ctx context.Context, ownerID string, groupID string) ([]dto.JoinRequestResponse, error) {
	return ListJoinRequests(ctx, uc.db, ownerID, groupID)
}

func (uc GroupUsecase) ListMembers(ctx context.Context, ownerID string, groupID string) ([]dto.GroupMemberResponse, error) {
	return ListMembers(ctx, uc.db, ownerID, groupID)
}

func (uc GroupUsecase) ApproveJoinRequest(ctx context.Context, ownerID string, groupID string, requesterID string) error {
	return ApproveJoinRequest(ctx, uc.db, ownerID, groupID, requesterID)
}

func (uc GroupUsecase) LeaveGroup(ctx context.Context, userID string, groupID string) error {
	return LeaveGroup(ctx, uc.db, userID, groupID)
}

func (uc GroupUsecase) RemoveMember(ctx context.Context, ownerID string, groupID string, memberID string) error {
	return RemoveMember(ctx, uc.db, ownerID, groupID, memberID)
}

func (uc GroupUsecase) TransferOwnership(ctx context.Context, ownerID string, groupID string, nextOwnerID string) error {
	return TransferOwnership(ctx, uc.db, ownerID, groupID, nextOwnerID)
}

func ListGroups(ctx context.Context, db Datastore, userID string) ([]dto.GroupListItemResponse, error) {
	return repositories.ListActiveUserGroups(ctx, db, userID)
}

func JoinGroup(ctx context.Context, db Datastore, userID string, displayName string, inviteCode string) (dto.JoinGroupResponse, error) {
	groupSummary, err := repositories.GroupInviteSummaryByCode(ctx, db, inviteCode)
	if errors.Is(err, repositories.ErrNotFound) {
		return dto.JoinGroupResponse{}, ErrGroupNotFound
	}
	if err != nil {
		return dto.JoinGroupResponse{}, err
	}

	currentStatus, err := repositories.GroupMemberStatus(ctx, db, groupSummary.ID, userID)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return dto.JoinGroupResponse{}, err
	}
	if currentStatus == "pending" {
		group, err := groupByID(ctx, db, groupSummary.ID, userID, "member", "pending")
		return dto.JoinGroupResponse{Group: group, MembershipStatus: "pending"}, err
	}
	if currentStatus == "active" {
		group, err := groupByID(ctx, db, groupSummary.ID, userID, "member", "active")
		return dto.JoinGroupResponse{Group: group, MembershipStatus: "active"}, err
	}

	if groupSummary.ParticipantLimit != nil && groupSummary.MemberCount >= *groupSummary.ParticipantLimit {
		return dto.JoinGroupResponse{}, ErrGroupFull
	}

	nextStatus := "active"
	if groupSummary.IsPrivate {
		nextStatus = "pending"
	}

	if err := repositories.InsertGroupMember(ctx, db, groupSummary.ID, userID, nextStatus, displayName); err != nil {
		return dto.JoinGroupResponse{}, err
	}

	group, err := groupByID(ctx, db, groupSummary.ID, userID, "member", nextStatus)
	if err != nil {
		return dto.JoinGroupResponse{}, err
	}

	return dto.JoinGroupResponse{Group: group, MembershipStatus: nextStatus}, nil
}

func ListJoinRequests(ctx context.Context, db Datastore, ownerID string, groupID string) ([]dto.JoinRequestResponse, error) {
	return repositories.ListPendingJoinRequests(ctx, db, ownerID, groupID)
}

func ListMembers(ctx context.Context, db Datastore, ownerID string, groupID string) ([]dto.GroupMemberResponse, error) {
	return repositories.ListActiveGroupMembers(ctx, db, ownerID, groupID)
}

func ApproveJoinRequest(ctx context.Context, db Datastore, ownerID string, groupID string, requesterID string) error {
	return withTx(ctx, db, func(tx repositories.Querier) error {
		capacity, err := repositories.OwnerGroupCapacity(ctx, tx, ownerID, groupID)
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrGroupNotFound
		}
		if err != nil {
			return err
		}

		if capacity.ParticipantLimit != nil && capacity.MemberCount >= *capacity.ParticipantLimit {
			return ErrGroupFull
		}

		err = repositories.ApprovePendingMember(ctx, tx, groupID, requesterID)
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrGroupNotFound
		}

		return err
	})
}

func LeaveGroup(ctx context.Context, db Datastore, userID string, groupID string) error {
	membership, err := repositories.GroupMembershipByUser(ctx, db, groupID, userID)
	if errors.Is(err, repositories.ErrNotFound) {
		return ErrGroupNotFound
	}
	if err != nil {
		return err
	}
	if membership.Role == "owner" {
		return ErrGroupOwnerRequired
	}
	if membership.Status != "active" {
		return ErrGroupNotFound
	}

	err = repositories.DeleteOwnGroupMembership(ctx, db, groupID, userID)
	if errors.Is(err, repositories.ErrNotFound) {
		return ErrGroupNotFound
	}

	return err
}

func RemoveMember(ctx context.Context, db Datastore, ownerID string, groupID string, memberID string) error {
	err := repositories.DeleteGroupMemberByOwner(ctx, db, ownerID, groupID, memberID)
	if errors.Is(err, repositories.ErrNotFound) {
		return ErrGroupNotFound
	}

	return err
}

func TransferOwnership(ctx context.Context, db Datastore, ownerID string, groupID string, nextOwnerID string) error {
	err := repositories.TransferGroupOwnership(ctx, db, ownerID, groupID, nextOwnerID)
	if errors.Is(err, repositories.ErrNotFound) {
		return ErrGroupNotFound
	}

	return err
}

func NormalizeCreateGroupRequest(request dto.CreateGroupRequest) (dto.CreateGroupRequest, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.MatchScope = strings.TrimSpace(request.MatchScope)

	if request.Name == "" {
		return request, apperrors.NewValidation("Informe o nome do grupo.")
	}

	if request.MatchScope != "all" && request.MatchScope != "selected" {
		return request, apperrors.NewValidation("Informe uma abrangencia de jogos valida.")
	}

	if request.MatchScope == "all" {
		request.SelectedTeams = []string{}
	}

	if request.MatchScope == "selected" {
		request.SelectedTeams = normalizeTeams(request.SelectedTeams)
		if len(request.SelectedTeams) == 0 {
			return request, apperrors.NewValidation("Selecione pelo menos uma selecao.")
		}
	}

	if request.HasUnlimitedParticipants {
		request.ParticipantLimit = nil
	} else if request.ParticipantLimit == nil || *request.ParticipantLimit < 2 {
		return request, apperrors.NewValidation("O limite precisa ser maior que 1.")
	}

	return request, nil
}

func NormalizeUpdateGroupRequest(request dto.UpdateGroupRequest) (dto.UpdateGroupRequest, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if request.Name == "" {
		return request, apperrors.NewValidation("Informe o nome do grupo.")
	}

	if request.HasUnlimitedParticipants {
		request.ParticipantLimit = nil
	} else if request.ParticipantLimit == nil || *request.ParticipantLimit < 2 {
		return request, apperrors.NewValidation("O limite precisa ser maior que 1.")
	}

	return request, nil
}

func NormalizeInviteCode(inviteCode string) string {
	inviteCode = strings.TrimSpace(inviteCode)
	inviteCode = strings.ToUpper(inviteCode)
	inviteCode = strings.ReplaceAll(inviteCode, " ", "")
	inviteCode = strings.ReplaceAll(inviteCode, "-", "")

	return inviteCode
}

func CreateGroup(ctx context.Context, db Datastore, userID string, displayName string, request dto.CreateGroupRequest) (dto.GroupResponse, error) {
	var group dto.GroupResponse
	err := withTx(ctx, db, func(tx repositories.Querier) error {
		for range 5 {
			inviteCode, err := generateInviteCode()
			if err != nil {
				return err
			}

			group, err = repositories.InsertGroupWithOwner(ctx, tx, userID, displayName, request, inviteCode)
			if err == nil {
				return nil
			}

			if repositories.IsUniqueViolation(err) {
				continue
			}

			return err
		}

		return errors.New("failed to generate unique invite code")
	})
	if err != nil {
		return dto.GroupResponse{}, err
	}

	return group, nil
}

func UpdateGroup(ctx context.Context, db Datastore, ownerID string, groupID string, request dto.UpdateGroupRequest) (dto.GroupResponse, error) {
	group, err := repositories.UpdateOwnedGroup(ctx, db, ownerID, groupID, request)
	if errors.Is(err, repositories.ErrNotFound) {
		return dto.GroupResponse{}, ErrGroupNotFound
	}
	if err != nil {
		return dto.GroupResponse{}, err
	}

	return group, nil
}

func groupByID(ctx context.Context, db Datastore, groupID string, userID string, role string, status string) (dto.GroupListItemResponse, error) {
	group, err := repositories.GroupListItemByID(ctx, db, groupID)
	if errors.Is(err, repositories.ErrNotFound) {
		return dto.GroupListItemResponse{}, ErrGroupNotFound
	}
	if err != nil {
		return dto.GroupListItemResponse{}, err
	}

	group.Role = role
	group.Status = status

	if group.OwnerID == userID {
		group.Role = "owner"
	}

	return group, nil
}

func normalizeTeams(teams []string) []string {
	seen := map[string]bool{}
	normalizedTeams := make([]string, 0, len(teams))

	for _, team := range teams {
		team = strings.TrimSpace(team)
		if team == "" || seen[team] {
			continue
		}

		seen[team] = true
		normalizedTeams = append(normalizedTeams, team)
	}

	return normalizedTeams
}

func generateInviteCode() (string, error) {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	for index, value := range buffer {
		buffer[index] = inviteCodeAlphabet[int(value)%len(inviteCodeAlphabet)]
	}

	return string(buffer), nil
}
