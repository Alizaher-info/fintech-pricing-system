// Package jwt stellt zentrale JWT-Token-Verwaltung für alle Go-Services bereit
// Diese Bibliothek ermöglicht die Generierung und Validierung von JWT-Tokens
// für die Authentifizierung in microservice-basierten Systemen
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5" // KORREKTER Import der JWT-Bibliothek
)


type Manager struct {
	secretKey []byte
	issuer    string
}


type UserClaims struct {
	UserID               int64  `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Standard JWT-Claims (iss, exp, iat, etc.)
}

// NewManager erstellt eine neue JWT-Manager-Instanz
func NewManager(secret, issuer string) *Manager {
	return &Manager{
		secretKey: []byte(secret), // Konvertiert String zu Byte-Array für JWT-Lib
		issuer:    issuer,         // Speichert Aussteller für Token-Claims
	}
}


// HINWEIS: In der Produktion wird diese Funktion hauptsächlich für Tests verwendet,
// da die Token-Generierung normalerweise im PHP-Backend stattfindet

func (m *Manager) GenerateToken(userID int64, email, role string) (string, error) {
	// Erstelle Claims mit Benutzerdaten und Standard-JWT-Feldern
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,                                           // Wer hat den Token ausgestellt
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token läuft nach 24h ab
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Wann wurde Token erstellt
		},
	}

	// Erstelle neuen Token mit HS256-Algorithmus und den Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signiere den Token mit dem geheimen Schlüssel und gib ihn als String zurück
	return token.SignedString(m.secretKey)
}


// Verwendung: Wird von gRPC-Interceptors aufgerufen für jeden authentifizierten Request
func (m *Manager) ValidateToken(tokenString string) (*UserClaims, error) {
	// Parse Token mit Custom Claims und Validierungs-Callback
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Sicherheitscheck: Stelle sicher, dass HS256-Algorithmus verwendet wird
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Gib den geheimen Schlüssel für Signatur-Validierung zurück
		return m.secretKey, nil
	})

	// Prüfe ob beim Parsing Fehler aufgetreten sind
	if err != nil {
		return nil, err
	}

	// Extrahiere und validiere Claims
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil // Token ist gültig, gib Benutzerinformationen zurück
	}

	return nil, errors.New("invalid token")
}


