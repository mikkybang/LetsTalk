package model

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson"
)

var hashKey = []byte("f8e3d17d08d04ea3c27e61ddc0daf98d5ce911f69e50ad7e36da335354909f4e")

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte{69, 96, 123, 60, 87, 130, 59, 101, 151, 171, 191, 53, 108, 112, 170, 26, 163, 68, 160, 193, 103, 182, 108, 4, 150, 91, 83, 11, 118, 13, 179, 219}
var cookieHandler = securecookie.New(hashKey, blockKey)

func (b CookieDetail) CreateCookie(w http.ResponseWriter) error {
	exitTime := time.Now().Add(time.Hour * 2)
	b.Data["exitTime"] = exitTime.Local()
	b.Data["UUID"] = UUID.String()

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

	var ff map[string]interface{}
	result.Decode(&ff)

	// todo: fix this so there wont be a crash
	if ff["loginUUID"].(string) != b.Data["UUID"].(string) {
		return errors.New("invalid uuid")
	}

	// todo: also check for expiry time.
	return nil
}
