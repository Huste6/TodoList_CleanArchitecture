package common

import "time"

type SQLModel struct {
	Id        int        `json:"-" gorm:"column:id;"`
	FakeID    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"`
}

func (SQLModel *SQLModel) Mask(dbType DbType) {
	uid := NewUID(uint32(SQLModel.Id), int(dbType), 1)
	SQLModel.FakeID = &uid
}
