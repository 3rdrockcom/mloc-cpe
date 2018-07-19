package customer

import "github.com/epointpayment/mloc-cpe/app/codes"

var (
	// MsgInfoUpdated is a message given when customer information was updated
	MsgInfoUpdated = codes.New("MESSAGE_CUSTOMER_INFORMATION_UPDATED").
			WithMessage("Customer information has been updated successfully").
			InGroup("CUSTOMER").
			RegisterMessage()

	// MsgTransactionsInserted is a message given when customer transactions were inserted
	MsgTransactionsInserted = codes.New("MESSAGE_CUSTOMER_TRANSACTIONS_INSERTED").
				WithMessage("Customer transactions have been inserted successfully").
				InGroup("CUSTOMER").
				RegisterMessage()
)
