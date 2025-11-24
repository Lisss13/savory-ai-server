package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	Secret     string
	Expiration time.Duration
}

type JWTData struct {
	ID        uint
	Email     string
	CompanyID uint
	Role      string
}

func NewJWT(secret string, expiration time.Duration) *JWT {
	return &JWT{
		Secret:     secret,
		Expiration: expiration,
	}
}

func (j *JWT) GenerateToken(data JWTData) (string, time.Time, error) {

	mySigningKey := []byte(j.Secret)
	exp := time.Now().Add(j.Expiration)

	jt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      data.Email,
		"id":         data.ID,
		"company_id": data.CompanyID,
		"role":       data.Role,
	})
	s, err := jt.SignedString(mySigningKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return s, exp, nil
}

func (j *JWT) ParseToken(tokenString string) (bool, *JWTData) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	data := &JWTData{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		data.Email = fmt.Sprint(claims["email"])
		data.ID = uint(claims["id"].(float64))
		data.CompanyID = uint(claims["company_id"].(float64))
		if role, ok := claims["role"].(string); ok {
			data.Role = role
		} else {
			data.Role = "user"
		}
	}

	return token.Valid, data
}
