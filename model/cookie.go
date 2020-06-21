package model

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/metaclips/LetsTalk/values"
	"go.mongodb.org/mongo-driver/bson"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func (b CookieDetail) CreateCookie(w http.ResponseWriter) error {
	exitTime := time.Now().Add(time.Hour * 2)
	b.Data.ExitTime = exitTime.Local()
	b.Data.UUID = uuid.New().String()

	_, err := db.Collection(b.Collection).UpdateOne(context.TODO(), map[string]interface{}{"_id": b.Email},
		bson.M{"$set": bson.M{"loginUUID": b.Data.UUID, "expires": exitTime}})
	if err != nil {
		return err
	}

	encoded, err := cookieHandler.Encode(b.CookieName, b.Data)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     b.CookieName,
		Value:    encoded,
		Expires:  exitTime,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     b.Path,
	}

	http.SetCookie(w, cookie)
	return nil
}

func (b *CookieDetail) CheckCookie(r *http.Request, w http.ResponseWriter) error {
	cookie, err := r.Cookie(b.CookieName)
	if err != nil {
		return err
	}

	if err := cookieHandler.Decode(b.CookieName, cookie.Value, &b.Data); err != nil {
		// Reset cookies if cookie validation fails.
		http.SetCookie(w, &http.Cookie{
			Name:    b.CookieName,
			Expires: time.Now(),
		})

		return err
	}

	b.Email = b.Data.Email
	result := db.Collection(b.Collection).FindOne(context.TODO(), map[string]string{"_id": b.Email})

	if err := result.Err(); err != nil {
		return err
	}

	var data map[string]interface{}
	if err = result.Decode(&data); err != nil {
		return err
	}

	if data["loginUUID"] != b.Data.UUID {
		return values.ErrIncorrectUUID
	}

	// TODO: also check for expiry time.
	return nil
}
