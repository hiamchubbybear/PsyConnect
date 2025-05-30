package middleware

import (
	"consultationservice/pkg/apiresponse"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleRequire(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerRoles := c.GetHeader("X-Roles")
		if headerRoles == "" {
			apiresponse.ErrorHandler(c, 401, "Missing roles header")
			c.Abort()
			return
		}
		roles := strings.Split(headerRoles, " ")
		role := roles[0]
		c.Set("roles", role)
		log.Printf("Print all roles %v", role)
		if required == "" {
			c.Next()
			return
		}
		expectedRole := fmt.Sprintf("role.%s:permission", required)
		if role == "role.admin:permission" || role == expectedRole {
			c.Next()
			return
		}
		apiresponse.ErrorHandler(c, 401, "Unauthenticated")
		c.Abort()
	}
}
