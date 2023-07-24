package databases

import (
	"backend/models"
	"backend/utils"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func GetTransactionsByUserId(UserId string) ([]models.Transaction, error) {
	dbClient := utils.GetFirestoreClient()
	ctx := context.Background()

	transactionLists := make([]models.Transaction, 0)

	itr := dbClient.Collection("user").
		Doc(UserId).
		Collection("transaction").
		OrderBy("Timestamp", firestore.Desc).
		Documents(ctx)
	
	for {
		transactionDoc, err := itr.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return transactionLists, err
		}

		// Get transaction
		transaction := models.Transaction{}
		transactionDoc.DataTo(&transaction)
		transaction.TransactionId = transactionDoc.Ref.ID

		// Append to transaction lists
		transactionLists = append(transactionLists, transaction)
	}

	return transactionLists, nil
}

func CreateTransaction(userId string, creatingTransaction models.CreatingTransaction) error {
	dbClient := utils.GetFirestoreClient()
	ctx := context.Background()

	_, _, err := dbClient.Collection("user").Doc(userId).Collection("transaction").Add(ctx, creatingTransaction)

	return err
}