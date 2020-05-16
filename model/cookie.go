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

// todo: fix this
var hashKey = []byte("f8e3d17d08d04ea3c27e61ddc0daf98d5ce911f69e50ad7e36da335354909f4e")
var blockKey = []byte{69, 96, 123, 60, 87, 130, 59, 101, 151, 171, 191, 53, 108, 112, 170, 26, 163, 68, 160, 193, 103, 182, 108, 4, 150, 91, 83, 11, 118, 13, 179, 219}
var cookieHandler = securecookie.New(hashKey, blockKey)

func (b CookieDetail) CreateCookie(w http.ResponseWriter) error {
	exitTime := time.Now().Add(time.Hour * 2)
	b.Data["exitTime"] = exitTime.Local()
	b.Data["UUID"] = uuid.New().String()

	_, err := db.Collection(b.Collection).UpdateOne(context.TODO(), map[string]interface{}{"_id": b.Email},
		bson.M{"$set": bson.M{"loginUUID": b.Data["UUID"], "expires": exitTime}})
	if err != nil {
		return err
	}

	if encoded, err := cookieHandler.Encode(b.CookieName, b.Data); err == nil {
		cookie := &http.Cookie{
			Name:    b.CookieName,
			Value:   encoded,
			Expires: exitTime,
			Path:    b.Path,
		}

		http.SetCookie(w, cookie)
	} else {
		return err
	}

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

	email, ok := b.Data["Email"].(string)
	if !ok {
		email = ""
	}
	b.Email = email
	result := db.Collection(b.Collection).FindOne(context.TODO(), map[string]interface{}{"_id": email})
	if err := result.Err(); err != nil {
		return err
	}

	var data map[string]interface{}
	if err = result.Decode(&data); err != nil {
		return err
	}

	if data["loginUUID"] != b.Data["UUID"] {
		return values.ErrIncorrectUUID
	}

	// todo: also check for expiry time.
	return nil
}
