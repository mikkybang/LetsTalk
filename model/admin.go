package model

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/metaclips/LetsTalk/values"
)

func (b *Admin) CheckAdminDetails(password string) error {
	result := db.Collection(values.AdminCollectionName).FindOne(ctx, bson.M{
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

func (b *Admin) CreateAdmin() error {
	_, err := db.Collection(values.AdminCollectionName).InsertOne(ctx, b)
	return err
}

func (b User) UploadUser(r *http.Request) error {
	b.Name = strings.Title(b.Name)
	if names := strings.Split(b.Name, " "); len(names) > 1 {
		var err error
		b.Password, err = bcrypt.GenerateFromPassword([]byte(names[0]), values.DefaultCost)
		if err != nil {
			return err
		}
	}

	id := strings.Split(b.Email, "@")
	if len(id) > 1 {
		b.ID = id[0]
	}

	if b.Class == "student" {
		b.ParentEmail = r.FormValue("parentEmail")
		b.ParentNumber = r.FormValue("parentNumber")
	}

	values.MapEmailToName[b.Email] = b.Name
	_, err := db.Collection(values.UsersCollectionName).InsertOne(ctx, b)
	return err
}
