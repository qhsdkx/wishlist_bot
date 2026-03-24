package group

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"wishlist-bot/internal/logger/sl"
)

type Repository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewRepository(db *sql.DB, log *slog.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (r *Repository) Create(g *Group) error {
	const op = "GroupRepository.Create"

	query := `INSERT INTO birthday_groups (name, birthday_user_id, status, created_at)
			  VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, g.Name, g.BirthdayUserID, g.Status, g.CreatedAt).Scan(&g.ID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error creating group: %w", err)
	}
	return nil
}

func (r *Repository) UpdateStatus(groupID int64, status string) error {
	const op = "GroupRepository.UpdateStatus"

	query := `UPDATE birthday_groups SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, groupID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error updating status: %w", err)
	}
	return nil
}

func (r *Repository) FindByID(id int64) (*Group, error) {
	const op = "GroupRepository.FindByID"

	query := `SELECT id, name, birthday_user_id, status, created_at
			  FROM birthday_groups WHERE id = $1`
	g := &Group{}
	err := r.db.QueryRow(query, id).Scan(&g.ID, &g.Name, &g.BirthdayUserID, &g.Status, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Error(op, sl.Err(err))
		return nil, fmt.Errorf("error finding group: %w", err)
	}
	return g, nil
}

func (r *Repository) FindByBirthdayUserID(userID int64) (*Group, error) {
	const op = "GroupRepository.FindByBirthdayUserID"

	query := `SELECT id, name, birthday_user_id, status, created_at
			  FROM birthday_groups WHERE birthday_user_id = $1`
	g := &Group{}
	err := r.db.QueryRow(query, userID).Scan(&g.ID, &g.Name, &g.BirthdayUserID, &g.Status, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Error(op, sl.Err(err))
		return nil, fmt.Errorf("error finding group: %w", err)
	}
	return g, nil
}

func (r *Repository) FindAllForUser(userID int64) ([]Group, error) {
	const op = "GroupRepository.FindAllForUser"

	query := `SELECT DISTINCT g.id, g.name, g.birthday_user_id, g.status, g.created_at
			  FROM birthday_groups g
			  LEFT JOIN birthday_group_members gm ON g.id = gm.group_id AND gm.user_id = $1
			  WHERE g.birthday_user_id = $1 OR gm.user_id IS NOT NULL
			  ORDER BY g.created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return nil, fmt.Errorf("error finding groups: %w", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.BirthdayUserID, &g.Status, &g.CreatedAt); err != nil {
			r.log.Error(op, sl.Err(err))
			return nil, fmt.Errorf("error scanning group: %w", err)
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *Repository) AddMember(groupID, userID int64) error {
	const op = "GroupRepository.AddMember"

	query := `INSERT INTO birthday_group_members (group_id, user_id, joined_at)
			  VALUES ($1, $2, $3)
			  ON CONFLICT (group_id, user_id) DO NOTHING`
	_, err := r.db.Exec(query, groupID, userID, time.Now())
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error adding member: %w", err)
	}
	return nil
}

func (r *Repository) RemoveMember(groupID, userID int64) error {
	const op = "GroupRepository.RemoveMember"

	query := `DELETE FROM birthday_group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, groupID, userID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error removing member: %w", err)
	}
	return nil
}

func (r *Repository) FindMembersByGroupID(groupID int64) ([]GroupMember, error) {
	const op = "GroupRepository.FindMembersByGroupID"

	query := `SELECT id, group_id, user_id, joined_at
			  FROM birthday_group_members WHERE group_id = $1`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return nil, fmt.Errorf("error finding members: %w", err)
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var m GroupMember
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.JoinedAt); err != nil {
			r.log.Error(op, sl.Err(err))
			return nil, fmt.Errorf("error scanning member: %w", err)
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *Repository) DeleteOldGroups(daysOld int) error {
	const op = "GroupRepository.DeleteOldGroups"

	query := `DELETE FROM birthday_groups WHERE created_at < NOW() - INTERVAL '1 day' * $1`
	_, err := r.db.Exec(query, daysOld)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error deleting old groups: %w", err)
	}
	return nil
}

func (r *Repository) IsMember(groupID, userID int64) bool {
	const op = "GroupRepository.IsMember"

	query := `SELECT EXISTS(SELECT 1 FROM birthday_group_members WHERE group_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, groupID, userID).Scan(&exists)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return false
	}
	return exists
}
