package persistance

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

type AddressRepo struct {
	DB *sqlx.DB
}

type AddressRepository interface {
	FindUserAddress(ctx context.Context, userID int) (*[]entity.UserAddress, error)
	SaveAddress(ctx context.Context, address entity.UserAddress) error
	UpdateAddress(ctx context.Context, id int, address map[string]interface{}) error
	DeleteAddress(ctx context.Context, id int, userID int) error
}

func (storage AddressRepo) SaveAddress(ctx context.Context, address entity.UserAddress) error  {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
	Insert("user_address").
	Columns("user_id", "name", "phone", "address", "city", "postal_code").
	Values(address.UserId, address.Name, address.Phone, address.Address, address.City, address.PostalCode).
	ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage AddressRepo) UpdateAddress(ctx context.Context, id int, newAddress map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
	Update("user_address").SetMap(newAddress).
	Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage AddressRepo) FindUserAddress(ctx context.Context, userID int) (*[]entity.UserAddress, error) {
	sql, _, _ := sq.
	Select("name", "phone", "address", "city", "postal_code").
	From("user_address").
	Where("user_id=$1").ToSql()
	
	rows, err := storage.DB.QueryxContext(ctx, sql, userID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	addresses := []entity.UserAddress{}

	for rows.Next() {
		var address entity.UserAddress
		err := rows.StructScan(&address)
		if err != nil {
			log.Println(err)
		}
		addresses = append(addresses, address)
	}

	return &addresses, nil
}

func (storage AddressRepo) DeleteAddress(ctx context.Context, id int, userID int) error {

	sql, _, _ := sq.Delete("user_address").Where("id=$1 AND user_id=$2").ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, id, userID)
	log.Println(sql)

	if err != nil {
		return err
	}

	return nil
}