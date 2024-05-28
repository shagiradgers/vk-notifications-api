package dao

import (
	sq "github.com/Masterminds/squirrel"
	"vk-notifications-api/internal/storage"
)

type DAO interface {
	NewUserQuery() UserQuery
	NewNotificationQuery() NotificationQuery

	Close() error
}

type dao struct {
	db storage.Storage
}

func (d *dao) NewUserQuery() UserQuery {
	return NewUserQuery(d.db)
}

func (d *dao) NewNotificationQuery() NotificationQuery {
	return NewNotificationQuery(d.db)
}

func (d *dao) Close() error {
	return d.db.Close()
}

func NewDAO(s storage.Storage) DAO {
	return &dao{
		db: s,
	}
}

func qb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
