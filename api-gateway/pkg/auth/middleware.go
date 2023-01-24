package auth

import (
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

// auth middleware
func AuthJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, "permission denied")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteJSON(w, http.StatusBadRequest, "permission denied")
			return
		}
		uid, err := utils.GetUUID(r)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, "permission denied")
			return
		}

		if claims["id"] != uid.String() {
			utils.WriteJSON(w, http.StatusBadRequest, "permission denied")
			return
		}

		next(w, r)
	}
}