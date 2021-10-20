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
	return handler.UserHandler{
		UserUseCase: usecase.UserUseCase{
			UserRepository: persistance.UserRepo{DB: i.DB},
		},
	}
}

func (i *Interactor) NewAddressHandler() handler.AddressHandler {
	return handler.AddressHandler{
		AddressUseCase: usecase.AddressUseCase{
			AddressRepository: persistance.AddressRepo{DB: i.DB},
		},
	}
}

func (i *Interactor) NewProductHandler() handler.ProductHandler {
	return handler.ProductHandler{
		ProductUseCase: usecase.ProductUseCase{
			ProductRepository: persistance.ProductRepo{DB: i.DB},
			ToppingRepository: persistance.ToppingRepo{DB: i.DB},
		},
	}
}

func (i *Interactor) NewCartHandler() handler.CartHandler {
	return handler.CartHandler{
		CartUseCase: usecase.CartUseCase{
			CartRepository: persistance.CartRepo{DB: i.DB},
		},
	}
}

func (i *Interactor) NewTransasctionHandler() handler.TransactionHandler {
	return handler.NewTransactionHandler(
		usecase.NewTransactionUseCase(
			persistance.NewTransactionFinder(i.DB),
			persistance.NewTransactionTx(i.DB),
			persistance.NewTransactionMutator(i.DB),
		))
}
