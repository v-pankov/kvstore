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

func (md *MetaData) MarkDeleted() {
	md.DeletedAt = time.Now()
}

func (md *MetaData) MarkUpdated() {
	md.UpdatedAt = time.Now()
}

func (md *MetaData) MarkCreated() {
	md.CreatedAt = time.Now()
}
