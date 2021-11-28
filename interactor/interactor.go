package interactor

import (
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/handler"
	"github.com/yosepalexsander/waysbucks-api/persistance"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type Interactor struct {
	DB *sqlx.DB
}

type AppHandler struct {
	handler.UserHandler
	handler.AddressHandler
	handler.ProductHandler
	handler.CartHandler
	handler.TransactionHandler
}

func (i *Interactor) NewAppHandler() *AppHandler {
	appHandler := &AppHandler{}
	appHandler.UserHandler = i.NewUserHandler()
	appHandler.AddressHandler = i.NewAddressHandler()
	appHandler.ProductHandler = i.NewProductHandler()
	appHandler.CartHandler = i.NewCartHandler()
	appHandler.TransactionHandler = i.NewTransasctionHandler()
	return appHandler
}

func (i *Interactor) NewUserHandler() handler.UserHandler {
	return handler.NewUserHandler(usecase.NewUserUseCase(
		persistance.NewUserFinder(i.DB),
		persistance.NewUserMutator(i.DB),
	))
}

func (i *Interactor) NewAddressHandler() handler.AddressHandler {
	return handler.AddressHandler{
		AddressUseCase: usecase.NewAddressUseCase(
			persistance.NewAddressFinder(i.DB),
			persistance.NewAddressMutator(i.DB),
		),
	}
}

func (i *Interactor) NewProductHandler() handler.ProductHandler {
	return handler.NewProductHandler(usecase.NewProductUseCase(
		persistance.NewProductFinder(i.DB),
		persistance.NewProductMutator(i.DB),
		persistance.NewToppingRepo(i.DB),
	))
}

func (i *Interactor) NewCartHandler() handler.CartHandler {
	return handler.NewCartHandler(usecase.NewCartUseCase(persistance.NewCartRepo(i.DB)))
}

func (i *Interactor) NewTransasctionHandler() handler.TransactionHandler {
	return handler.NewTransactionHandler(
		usecase.NewTransactionUseCase(
			persistance.NewTransactionFinder(i.DB),
			persistance.NewTransactionTx(i.DB),
			persistance.NewTransactionMutator(i.DB),
		))
}
