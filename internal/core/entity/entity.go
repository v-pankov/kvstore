package entity

import "time"

type MetaData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (md *MetaData) IsDeleted() bool {
	return !(md.DeletedAt.IsZero() || md.DeletedAt.Before(md.UpdatedAt))
}
