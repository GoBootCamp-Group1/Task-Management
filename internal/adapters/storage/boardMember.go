package storage

import (
	"context"
	"errors"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type boardMemberRepo struct {
	db *gorm.DB
}

func NewBoardMemberRepo(db *gorm.DB) ports.BoardMemberRepo {
	return &boardMemberRepo{
		db: db,
	}
}

var (
	ErrBoardMemberNotFound = "Board member not found"
)

func (r *boardMemberRepo) Create(ctx context.Context, member *domains.BoardMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToBoardMemberEntity(member)
		if err := tx.WithContext(ctx).Table(entity.TableName()).Create(&entity).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		member.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *boardMemberRepo) GetByID(ctx context.Context, id uint) (*domains.BoardMember, error) {
	var m entities.BoardMember
	err := r.db.WithContext(ctx).Table(m.TableName()).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, ErrBoardMemberNotFound)
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.BoardMemberEntityToDomain(&m), nil
}

func (r *boardMemberRepo) Update(ctx context.Context, member *domains.BoardMember) error {
	var existingBoardMember entities.BoardMember
	if err := r.db.WithContext(ctx).Table(existingBoardMember.TableName()).Model(&entities.Board{}).Where("id = ?", member.ID).First(&existingBoardMember).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	existingBoardMember.RoleID = member.RoleID
	existingBoardMember.UserID = member.UserID
	existingBoardMember.BoardID = member.BoardID
	if err := r.db.WithContext(ctx).Save(&existingBoardMember).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *boardMemberRepo) Delete(ctx context.Context, id uint) error {
	var entity entities.BoardMember
	if err := r.db.WithContext(ctx).Table(entity.TableName()).Where("id = ?", id).Delete(&entity).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *boardMemberRepo) GetBoardMembers(ctx context.Context, boardID uint) ([]domains.BoardMember, error) {
	var boardMemberEntities []entities.BoardMember
	err := r.db.WithContext(ctx).Table(boardMemberEntities[0].TableName()).Where("board_id = ?", boardID).Find(&boardMemberEntities).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return mappers.BoardMemberEntitiesToDomain(boardMemberEntities), nil

}
func (r *boardMemberRepo) GetBoardMember(ctx context.Context, boardID, userID uint) (*domains.BoardMember, error) {
	var boardMember entities.BoardMember
	if err := r.db.WithContext(ctx).
		Table(boardMember.TableName()).
		Where("board_id = ? AND user_id = ?", boardID, userID).
		First(&boardMember).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if boardMember.ID == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, ErrBoardMemberNotFound)
	}

	return mappers.BoardMemberEntityToDomain(&boardMember), nil
}
