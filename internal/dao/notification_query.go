package dao

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"time"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/storage"
)

type NotificationQuery interface {
	GetNotification(ctx context.Context, ID int64) (NotificationTable, error)
	AddNotification(
		ctx context.Context,
		senderID int64,
		receiverIDs []int64,
		message string,
		mediaContent sql.NullString,
		date time.Time,
	) (NotificationTable, error)
	UpdateNotificationStatus(ctx context.Context, ID int64, status string) error
	GetNotificationsByIDs(
		ctx context.Context,
		IDs []int64,
		limit int64,
		offset int64,
	) ([]NotificationTable, error)
}

type notificationQuery struct {
	db storage.Storage
}

func NewNotificationQuery(db storage.Storage) NotificationQuery {
	return &notificationQuery{
		db: db,
	}
}

func (q *notificationQuery) GetNotification(
	ctx context.Context,
	ID int64,
) (NotificationTable, error) {
	var dest NotificationTable

	query := qb().
		Select(dest.columns()...).
		From(notificationTableName).
		Where(sq.Eq{"id": ID})
	err := q.db.GetX(ctx, &dest, query)
	return dest, err
}

func (q *notificationQuery) AddNotification(
	ctx context.Context,
	senderID int64,
	receiverIDs []int64,
	message string,
	mediaContent sql.NullString,
	date time.Time,
) (NotificationTable, error) {
	var dest NotificationTable
	query := qb().
		Insert(notificationTableName).
		Columns(
			"sender_id",
			"receiver_ids",
			"message",
			"media_content",
			"date",
			"status",
		).
		Values(
			senderID,
			pq.Array(receiverIDs),
			message,
			mediaContent,
			date,
			desc.NotificationStatus_CREATED.String(),
		).
		Suffix("RETURNING *")

	err := q.db.GetX(ctx, &dest, query)
	return dest, err
}

func (q *notificationQuery) UpdateNotificationStatus(
	ctx context.Context,
	ID int64,
	status string,
) error {
	query := qb().
		Update(notificationTableName).
		Set("status", status).
		Where(sq.Eq{"id": ID})

	return q.db.ExecX(ctx, query)
}

func (q *notificationQuery) GetNotificationsByIDs(
	ctx context.Context,
	IDs []int64,
	limit int64,
	offset int64,
) ([]NotificationTable, error) {
	var dest []NotificationTable

	query := qb().
		Select(NotificationTable{}.columns()...).
		From(notificationTableName).
		Where("id = any(?)", pq.Array(IDs)).
		Limit(uint64(limit)).
		Offset(uint64(offset))
	err := q.db.SelectX(ctx, &dest, query)
	return dest, err
}
