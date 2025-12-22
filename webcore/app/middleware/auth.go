package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/semanggilab/webcore-go/app/helper"
)

// // NewAuthFromConfig creates a new authentication middleware from the application configuration
// func NewAuthFromConfig(authConfig config.AuthConfig) fiber.Handler {
// 	// In a real implementation, you might want to add API key specific config
// 	// For now, we'll use the defaults

// 	return NewAuth(authConfig)
// }
//
// // NewAuth creates a new authentication middleware based on the configured type
// func NewAuth(config config.AuthConfig) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Get Authorization header
// 		authHeader := c.Get("Authorization")

// 		switch strings.ToLower(config.Type) {
// 		case "jwt":
// 			return validateJWT(c, authHeader, config.SecretKey)
// 		case "apikey":
// 			return validateAPIKey(c, authHeader, config.APIKeyHeader, config.APIKeyPrefix)
// 		case "none":
// 			// No authentication required
// 			return c.Next()
// 		default:
// 			return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 				HttpCode:  fiber.StatusUnauthorized,
// 				ErrorCode: 2,
// 				ErrorName: "UNAUTHORIZED",
// 				Message:   fmt.Sprintf("Unsupported authentication type: %s", config.Type),
// 			})
// 		}
// 	}
// }
//
// // validateJWT handles JWT token validation
// func validateJWT(c *fiber.Ctx, authHeader, secretKey string) error {
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 			HttpCode:  fiber.StatusUnauthorized,
// 			ErrorCode: 2,
// 			ErrorName: "UNAUTHORIZED",
// 			Message:   "Authorization header required",
// 		})
// 	}

// 	// Check if it's a Bearer token
// 	if !strings.HasPrefix(authHeader, "Bearer ") {
// 		return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 			HttpCode:  fiber.StatusUnauthorized,
// 			ErrorCode: 2,
// 			ErrorName: "UNAUTHORIZED",
// 			Message:   "Invalid authorization format. Expected 'Bearer <token>'",
// 		})
// 	}

// 	// Extract token
// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 	// Parse and validate token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
// 		// Validate the signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fiber.ErrUnauthorized
// 		}
// 		return []byte(secretKey), nil
// 	})

// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 			HttpCode:  fiber.StatusUnauthorized,
// 			ErrorCode: 2,
// 			ErrorName: "UNAUTHORIZED",
// 			Message:   "Invalid or expired token",
// 		})
// 	}

// 	// Extract claims
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// Store user info in context
// 		c.Locals("user_id", claims["user_id"])
// 		c.Locals("user_role", claims["role"])
// 		c.Locals("user_permissions", claims["permissions"])
// 		c.Locals("auth_type", "jwt")

// 		// Continue to next handler
// 		return c.Next()
// 	}

// 	return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 		HttpCode:  fiber.StatusUnauthorized,
// 		ErrorCode: 2,
// 		ErrorName: "UNAUTHORIZED",
// 		Message:   "Invalid token claims",
// 	})
// }

// // validateAPIKey handles API key validation
// func validateAPIKey(c *fiber.Ctx, authHeader, apiKeyHeader, apiKeyPrefix string) error {
// 	if apiKeyHeader == "" {
// 		apiKeyHeader = "X-API-Key" // Default header name
// 	}

// 	// Check Authorization header first (for backward compatibility)
// 	if authHeader != "" {
// 		if strings.HasPrefix(authHeader, "Bearer ") {
// 			// Handle Bearer token format
// 			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
// 			return validateAPIKeyValue(c, apiKey, apiKeyPrefix)
// 		} else if strings.HasPrefix(authHeader, "APIKey ") {
// 			// Handle APIKey prefix format
// 			apiKey := strings.TrimPrefix(authHeader, "APIKey ")
// 			return validateAPIKeyValue(c, apiKey, apiKeyPrefix)
// 		}
// 	}

// 	// Check dedicated API key header
// 	apiKey := c.Get(apiKeyHeader)
// 	if apiKey == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 			HttpCode:  fiber.StatusUnauthorized,
// 			ErrorCode: 2,
// 			ErrorName: "UNAUTHORIZED",
// 			Message:   fmt.Sprintf("API key required in Authorization header or %s header", apiKeyHeader),
// 		})
// 	}

// 	return validateAPIKeyValue(c, apiKey, apiKeyPrefix)
// }

// // validateAPIKeyValue validates the API key value
// func validateAPIKeyValue(c *fiber.Ctx, apiKey, apiKeyPrefix string) error {
// 	// Check prefix if specified
// 	if apiKeyPrefix != "" {
// 		if !strings.HasPrefix(apiKey, apiKeyPrefix) {
// 			return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
// 				HttpCode:  fiber.StatusUnauthorized,
// 				ErrorCode: 2,
// 				ErrorName: "UNAUTHORIZED",
// 				Message:   "Invalid API key format",
// 			})
// 		}
// 		apiKey = strings.TrimPrefix(apiKey, apiKeyPrefix)
// 	}

// 	// In a real implementation, you would validate the API key against your database
// 	// For now, we'll just store it in context and continue
// 	c.Locals("api_key", apiKey)
// 	c.Locals("auth_type", "apikey")

// 	// You could add additional validation here, such as:
// 	// - Check if API key exists in database
// 	// - Check if API key is active
// 	// - Extract user information associated with the API key

// 	return c.Next()
// }

// GetAuthType returns the authentication type from the context
func GetAuthType(c *fiber.Ctx) string {
	authType := c.Locals("auth_type")
	if authType == nil {
		return "unknown"
	}
	return authType.(string)
}

// GetUserID returns the user ID from the context
func GetUserID(c *fiber.Ctx) any {
	return c.Locals("user_id")
}

// GetUserRole returns the user role from the context
func GetUserRole(c *fiber.Ctx) any {
	return c.Locals("user_role")
}

// GetUserPermissions returns the user permissions from the context
func GetUserPermissions(c *fiber.Ctx) any {
	return c.Locals("user_permissions")
}

// GetAPIKey returns the API key from the context
func GetAPIKey(c *fiber.Ctx) string {
	apiKey := c.Locals("api_key")
	if apiKey == nil {
		return ""
	}
	return apiKey.(string)
}

// RoleRequired creates a middleware to check user roles
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authType := GetAuthType(c)

		// For API key authentication, you might want to implement different role checking logic
		if authType == "apikey" {
			// In a real implementation, you might want to check API key permissions/roles
			// For now, we'll allow access if API key is valid
			return c.Next()
		}

		// For JWT authentication, check the role claims
		userRole := GetUserRole(c)
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: 2,
				ErrorName: "UNAUTHORIZED",
				Message:   "User role not found in context",
			})
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(helper.APIError{
			HttpCode:  fiber.StatusUnauthorized,
			ErrorCode: 2,
			ErrorName: "UNAUTHORIZED",
			Message:   "Insufficient permissions",
		})
	}
}

// PermissionRequired creates a middleware to check user permissions
func PermissionRequired(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authType := GetAuthType(c)

		// For API key authentication, you might want to implement different permission checking logic
		if authType == "apikey" {
			// In a real implementation, you might want to check API key permissions
			// For now, we'll allow access if API key is valid
			return c.Next()
		}

		// For JWT authentication, check the permission claims
		userPermissions := GetUserPermissions(c)
		if userPermissions == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.APIError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: 2,
				ErrorName: "UNAUTHORIZED",
				Message:   "User permissions not found in context",
			})
		}

		permissions := userPermissions.([]any)
		for _, permission := range permissions {
			if permission == requiredPermission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(helper.APIError{
			HttpCode:  fiber.StatusUnauthorized,
			ErrorCode: 2,
			ErrorName: "UNAUTHORIZED",
			Message:   "Insufficient permissions",
		})
	}
}
