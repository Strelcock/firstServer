package link

import (
	"firstServer/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

type LinkRepositoryDeps struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (lr *LinkRepository) Create(link *Link) (*Link, error) {
	res := lr.Database.DB.Create(link)
	if res.Error != nil {
		return nil, res.Error
	}
	return link, nil
}

func (lr *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := lr.Database.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (lr *LinkRepository) Update(link *Link) (*Link, error) {
	result := lr.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (lr *LinkRepository) Delete(id uint) error {
	result := lr.Database.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (lr *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := lr.Database.DB.First(&link, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}
