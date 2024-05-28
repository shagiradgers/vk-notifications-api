package dao

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"vk-notifications-api/internal/storage"
)

type UserQuery interface {
	CreateUser(
		ctx context.Context,
		vkID int64,
		role string,
		group string,
		firstname string,
		surname string,
		patronymic sql.NullString,
		mobilePhone string,
		status string,
	) (UserTable, error)
	GetUser(
		ctx context.Context,
		userID int64,
	) (UserTable, error)
	DeleteUser(
		ctx context.Context,
		userID int64,
	) error
	UpdateUser(
		ctx context.Context,
		user UserTable,
		fields ...string,
	) (UserTable, error)
	GetUserByFilter(
		ctx context.Context,
		user UserTable,
		limit uint64,
		offset uint64,
		fields ...string,
	) ([]UserTable, error)
}

type userQuery struct {
	db storage.Storage
}

func NewUserQuery(db storage.Storage) UserQuery {
	return &userQuery{db: db}
}

func (q *userQuery) CreateUser(
	ctx context.Context,
	vkID int64,
	role string,
	group string,
	firstname string,
	surname string,
	patronymic sql.NullString,
	mobilePhone string,
	status string,
) (UserTable, error) {
	var dest UserTable
	query := qb().
		Insert(userTableName).
		Columns(
			"vk_id",
			"role",
			"user_group",
			"firstname",
			"surname",
			"patronymic",
			"mobile_phone",
			"status",
		).
		Values(
			vkID,
			role,
			group,
			firstname,
			surname,
			patronymic,
			mobilePhone,
			status,
		).
		Suffix("RETURNING *")
	err := q.db.GetX(ctx, &dest, query)
	return dest, err
}

func (q *userQuery) GetUser(
	ctx context.Context,
	userID int64,
) (UserTable, error) {
	var dest UserTable
	query := qb().
		Select(dest.columns()...).
		From(userTableName).
		Where(sq.Eq{"id": userID})

	err := q.db.GetX(ctx, &dest, query)
	return dest, err
}

func (q *userQuery) DeleteUser(
	ctx context.Context,
	userID int64,
) error {
	query := qb().
		Delete(userTableName).
		Where(sq.Eq{"id": userID})
	return q.db.ExecX(ctx, query)
}

func (q *userQuery) UpdateUser(
	ctx context.Context,
	user UserTable,
	fields ...string,
) (UserTable, error) {
	var dest UserTable
	userMap := user.toMap()

	query := qb().
		Update(userTableName).
		Where(sq.Eq{"id": user.Id})

	for _, field := range fields {
		query = query.Set(field, userMap[field])
	}
	query = query.Suffix("RETURNING *")

	err := q.db.GetX(ctx, &dest, query)
	return dest, err
}

func (q *userQuery) GetUserByFilter(
	ctx context.Context,
	user UserTable,
	limit uint64,
	offset uint64,
	fields ...string,
) ([]UserTable, error) {
	var dest []UserTable
	userMap := user.toMap()

	query := qb().
		Select(UserTable{}.columns()...).
		From(userTableName)

	for _, field := range fields {
		query = query.Where(sq.Eq{field: userMap[field]})
	}
	query = query.
		Limit(limit).
		Offset(offset)

	err := q.db.SelectX(ctx, &dest, query)
	return dest, err
}
