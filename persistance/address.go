package persistance

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type addressRepo struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) repository.AddressRepository {
	return &addressRepo{db}
}

func (storage *addressRepo) SaveAddress(ctx context.Context, userID string, address entity.Address) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
		Insert("user_address").
		Columns("id", "user_id", "name", "phone", "address", "city", "postal_code").
		Values(address.Id, userID, address.Name, address.Phone, address.Address, address.City, address.PostalCode).
		ToSql()
	_, err := storage.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *addressRepo) FindAllUserAddresses(ctx context.Context, userID string) ([]entity.Address, error) {
	sql, _, _ := sq.
		Select("id", "name", "phone", "address", "city", "postal_code").
		From("user_address").
		Where("user_id=$1").ToSql()

	addresses := []entity.Address{}

	rows, err := storage.db.QueryxContext(ctx, sql, userID)
	for rows.Next() {
		address := entity.Address{}
		err = rows.StructScan(&address)
		addresses = append(addresses, address)
	}

	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (storage *addressRepo) FindAddress(ctx context.Context, id string) (*entity.Address, error) {
	sql, _, _ := sq.
		Select("id", "user_id", "name", "phone", "address", "city", "postal_code").
		From("user_address").Where("id=$1").ToSql()

	var address entity.Address
	err := storage.db.QueryRowxContext(ctx, sql, id).StructScan(&address)

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (storage *addressRepo) UpdateAddress(ctx context.Context, id string, newAddress map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
		Update("user_address").SetMap(newAddress).
		Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *addressRepo) DeleteAddress(ctx context.Context, id string, userID string) error {
	sql, _, _ := sq.Delete("user_address").Where("id=$1 AND user_id=$2").ToSql()
	result, err := storage.db.ExecContext(ctx, sql, id, userID)

	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
