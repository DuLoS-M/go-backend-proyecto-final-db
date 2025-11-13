package middleware

import (
	"net/http"

	"proyecto-bd-final/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RequirePermission middleware que verifica si el usuario tiene el permiso requerido
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "No se encontraron roles de usuario", nil)
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al procesar roles", nil)
			c.Abort()
			return
		}

		// Verificar si el usuario tiene el permiso (esto se puede mejorar con una consulta a BD)
		// Por ahora, verificamos roles básicos
		hasPermission := false
		for _, role := range userRoles {
			if role == "admin" || role == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "No tienes permiso para acceder a este recurso", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware que verifica si el usuario tiene el rol requerido
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "No se encontraron roles de usuario", nil)
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al procesar roles", nil)
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range userRoles {
			if r == role || r == "admin" {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.ErrorResponse(c, http.StatusForbidden, "Rol insuficiente para acceder a este recurso", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
