package database

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func InitSession() {
    Store = session.New(session.Config{
        CookieHTTPOnly: true,
        CookieSecure:   false, // Change this to true in production
        CookieSameSite: "Lax",
        Expiration:     24 * time.Hour,
        KeyLookup:      "cookie:session_id",
    })
}