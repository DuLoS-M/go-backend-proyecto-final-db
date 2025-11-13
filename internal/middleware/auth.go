package middleware

import (
	"net/http"
	"strings"

	"proyecto-bd-final/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token de autenticación requerido", nil)
			c.Abort()
			return
		}

		// Extraer el token del header "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Formato de token inválido", nil)
			c.Abort()
			return
		}

		// Validar el token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token inválido o expirado", err)
			c.Abort()
			return
		}

		// Guardar los claims en el contexto para uso posterior
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}
