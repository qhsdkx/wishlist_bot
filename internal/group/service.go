package group

import (
	"fmt"
	"log/slog"
	"time"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/user"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
	log      *slog.Logger
}

func NewService(repo *Repository, userRepo *user.Repository, log *slog.Logger) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *Service) CreateGroupForBirthday(userID int64) (*Group, error) {
	existed, err := s.repo.FindByBirthdayUserID(userID)
	if err != nil {
		return nil, err
	}
	if existed != nil {
		return existed, nil
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.log.Error("GroupService.CreateGroupForBirthday", sl.Err(err))
		return nil, err
	}

	name := fmt.Sprintf("День рождения у %s %s", u.Name, u.Surname)

	g := &Group{
		Name:           name,
		BirthdayUserID: userID,
		Status:         GroupStatusUpcoming,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Create(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) GetGroupsForUser(userID int64) ([]Group, error) {
	return s.repo.FindAllForUser(userID)
}

func (s *Service) GetGroupByID(groupID int64) (*Group, error) {
	return s.repo.FindByID(groupID)
}

func (s *Service) FindByBirthdayUserID(userID int64) (*Group, error) {
	return s.repo.FindByBirthdayUserID(userID)
}

func (s *Service) JoinGroup(groupID, userID int64) error {
	g, err := s.repo.FindByID(groupID)
	if err != nil {
		return err
	}
	if g == nil {
		return fmt.Errorf("группа не найдена. Возможно, она уже устарела")
	}
	if g.BirthdayUserID == userID {
		return fmt.Errorf("вы не можете присоединиться к своей группе")
	}
	return s.repo.AddMember(groupID, userID)
}

func (s *Service) LeaveGroup(groupID, userID int64) error {
	return s.repo.RemoveMember(groupID, userID)
}

func (s *Service) GetGroupMembers(groupID int64) ([]GroupMember, error) {
	return s.repo.FindMembersByGroupID(groupID)
}

func (s *Service) IsMember(groupID, userID int64) bool {
	return s.repo.IsMember(groupID, userID)
}

func (s *Service) MarkGroupAsPassed(groupID int64) error {
	return s.repo.UpdateStatus(groupID, GroupStatusPassed)
}

func (s *Service) CleanupOldGroups() error {
	return s.repo.DeleteOldGroups(10)
}
