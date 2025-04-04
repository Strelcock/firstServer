package stat

import (
	"firstServer/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (sr *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	sr.Db.Find(&stat, "link_id = ? and date = ?", linkId, datatypes.Date(time.Now()))
	if stat.ID == 0 {
		sr.Db.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks++
		sr.Db.Save(&stat)
	}
}
