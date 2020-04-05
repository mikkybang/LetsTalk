package model

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/metaclips/FinalYearProject/values"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Admin) CheckAdminDetails(password string) error {
	result := db.Collection(values.AdminCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.StaffDetails.Email,
	})

	if err := result.Decode(&b); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(b.StaffDetails.Password, []byte(password)); err != nil {
		return err
	}
	return nil
}
