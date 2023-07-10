package thirdparty

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/yosepalexsander/waysbucks-api/config"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

func setupGlobalMidtransConfig() {
	midtrans.ServerKey = config.MIDTRANS_SERVER_KEY
	midtrans.ClientKey = config.MIDTRANS_CLIENT_KEY
	midtrans.Environment = midtrans.Sandbox
}

func CreateTransaction(t *entity.Transaction) *snap.Response {
	setupGlobalMidtransConfig()
	req := generateSnapReq(t)

	snapResp, _ := snap.CreateTransaction(req)
	return snapResp
}

func generateSnapReq(t *entity.Transaction) *snap.Request {
	custAddress := &midtrans.CustomerAddress{
		FName:       t.Name,
		LName:       "",
		Phone:       t.Phone,
		Address:     t.Address,
		City:        t.City,
		Postcode:    strconv.Itoa(t.PostalCode),
		CountryCode: "IDN",
	}

	var orderItems []midtrans.ItemDetails
	for _, order := range t.Orders {
		itemDetail := midtrans.ItemDetails{
			ID:    strconv.Itoa(order.Id),
			Name:  order.Name,
			Qty:   1,
			Price: int64(order.Price),
		}
		orderItems = append(orderItems, itemDetail)
	}

	// add service fee to item details because midtrans cannot put it automatically
	serviceFee := midtrans.ItemDetails{
		ID:    "FEE-" + t.Id,
		Name:  "Service Fee",
		Qty:   1,
		Price: int64(t.ServiceFee),
	}
	orderItems = append(orderItems, serviceFee)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.Id,
			GrossAmt: int64(t.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    t.Name,
			LName:    "",
			Email:    t.Email,
			Phone:    t.Phone,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		Items: &orderItems,
	}
	return req
}

func ParseTransactionResponse(reqBody io.ReadCloser) (*coreapi.TransactionStatusResponse, error) {
	transaction := new(coreapi.TransactionStatusResponse)

	if err := json.NewDecoder(reqBody).Decode(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}
