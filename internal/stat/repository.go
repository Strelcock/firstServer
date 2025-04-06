package stat

import (
	"firstServer/pkg/db"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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

func (sr *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	query := sr.DB.Table("stats").
		Select(selectQuery).
		Session(&gorm.Session{})

	query.Where("date BETWEEN ? and ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}
