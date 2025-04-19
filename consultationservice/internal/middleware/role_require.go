package middleware

import (
	"consultationservice/pkg/apiresponse"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func RoleRequire(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedRole := fmt.Sprintf("role.%s:permission", required)

		headerRoles := c.GetHeader("X-Roles")
		if headerRoles == "" {
			apiresponse.ErrorHandler(c, 401, "Missing roles header")
			c.Abort()
			return
		}

		roles := strings.Split(headerRoles, " ")
		if len(roles) == 0 || roles[0] != expectedRole {
			apiresponse.ErrorHandler(c, 401, "Unauthenticated")
			c.Abort()
			return
		}

		c.Next()
	}
}
