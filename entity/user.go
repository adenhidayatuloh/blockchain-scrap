package entity

import (
	"blockchain-scrap/pkg/errs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User merepresentasikan entitas pengguna dalam sistem
type User struct {
	ID                 uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email              string             `gorm:"unique;not null"`
	Password           string             `gorm:"not null"` // Password dalam bentuk hash
	BlockchainSearches []BlockchainSearch `gorm:"foreignKey:UserID"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// HashPassword mengenkripsi password pengguna menggunakan bcrypt
func (u *User) HashPassword() errs.MessageErr {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errs.NewInternalServerError("Gagal mengenkripsi password")
	}
	u.Password = string(hashedPassword)

	return nil
}

// ComparePassword membandingkan password yang diberikan dengan password yang tersimpan
func (u *User) ComparePassword(inputPassword string) errs.MessageErr {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))

	if err != nil {
		return errs.NewBadRequest("Password tidak valid!")
	}

	return nil
}

// CreateToken membuat JWT token untuk autentikasi
func (u *User) CreateToken() (string, errs.MessageErr) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    u.ID,
			"email": u.Email,
		})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Error saat menandatangani token:", err.Error())
		return "", errs.NewInternalServerError("Gagal menandatangani token JWT")
	}

	return signedToken, nil
}

// ParseToken memvalidasi dan mengurai token JWT
func (u *User) ParseToken(tokenString string) (*jwt.Token, errs.MessageErr) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, isValidMethod := t.Method.(*jwt.SigningMethodHMAC); !isValidMethod {
			return nil, errs.NewUnauthenticated("Metode token tidak valid")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, errs.NewUnauthenticated("Token tidak valid")
	}

	return token, nil
}

// ValidateToken memvalidasi token Bearer
func (u *User) ValidateToken(bearerToken string) errs.MessageErr {
	if !strings.HasPrefix(bearerToken, "Bearer") {
		return errs.NewUnauthenticated("Token harus menggunakan tipe Bearer")
	}

	tokenParts := strings.Fields(bearerToken)
	if len(tokenParts) != 2 {
		return errs.NewUnauthenticated("Format token tidak valid")
	}

	tokenString := tokenParts[1]
	token, err := u.ParseToken(tokenString)
	if err != nil {
		return err
	}

	claims, isValid := token.Claims.(jwt.MapClaims)
	if !isValid || !token.Valid {
		return errs.NewUnauthenticated("Token tidak valid")
	}

	return u.bindTokenToUserEntity(claims)
}

// bindTokenToUserEntity mengikat data dari token ke entitas user
func (u *User) bindTokenToUserEntity(claims jwt.MapClaims) errs.MessageErr {
	userID, hasID := claims["id"].(string)
	userEmail, hasEmail := claims["email"].(string)

	if !hasID {
		return errs.NewUnauthenticated("Token tidak mengandung ID")
	}

	if !hasEmail {
		return errs.NewUnauthenticated("Token tidak mengandung email")
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return errs.NewBadRequest("ID tidak valid")
	}

	u.ID = parsedUUID
	u.Email = userEmail

	return nil
}
