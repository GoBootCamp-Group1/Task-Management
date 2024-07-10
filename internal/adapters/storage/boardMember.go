package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
	"gorm.io/gorm"
)

type boardMemberRepo struct {
	db *gorm.DB
}

var (
	ErrBoardMemberAlreadyExists = errors.New("board member already exists")
)

func NewBoardMemberRepo(db *gorm.DB) port.BoardMemberRepo {
	return &boardMemberRepo{
		db: db,
	}
}

func (r *boardMemberRepo) Create(ctx context.Context, member *domain.BoardMember) error {
	var existingBoardMember entities.BoardMember
	err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("board_id = ? AND user_id = ?", member.BoardID, member.UserID).First(&existingBoardMember).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingBoardMember.ID != 0 {
		return ErrBoardMemberAlreadyExists
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToBoardMemberEntity(member)
		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return err
		}
		member.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *boardMemberRepo) GetByID(ctx context.Context, id uint) (*domain.BoardMember, error) {
	var m entities.BoardMember
	err := r.db.WithContext(ctx).Model(&entities.BoardMember{}).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mappers.BoardMemberEntityToDomain(&m), nil
}

func (r *boardMemberRepo) Update(ctx context.Context, member *domain.BoardMember) error {

	var existingBoardMember *entities.BoardMember
	if err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("id = ?", member.ID).First(&existingBoardMember).Error; err != nil {
		return err
	}
	existingBoardMember.RoleID = member.RoleID
	existingBoardMember.UserID = member.UserID
	existingBoardMember.BoardID = member.BoardID
	return r.db.WithContext(ctx).Save(&existingBoardMember).Error
}

func (r *boardMemberRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.BoardMember{}, id).Error
}

func (r *boardMemberRepo) GetBoardMembers(ctx context.Context, boardID uint) ([]domain.BoardMember, error) {
	var boardMemberEntities []entities.BoardMember
	err := r.db.WithContext(ctx).Where("board_id = ?", boardID).Find(&boardMemberEntities).Error
	if err != nil {
		return nil, err
	}

	return mappers.BoardMemberEntitiesToDomain(boardMemberEntities), nil

}
