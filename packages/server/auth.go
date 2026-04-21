package server

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	sharedadapters "github.com/marlon-clemente/timyo-playground-backend/packages/adapters"
)

// JwtPayload represents the structure of the JWT access token from the external authentication service.
type JwtPayload struct {
	jwt.RegisteredClaims
}

// AuthMiddleware creates a Fiber middleware that validates JWT tokens.
// It extracts the "sub" claim and populates the "userId" local variable in the context.
//
// The middleware expects an "Authorization" header with the format: Bearer <token>.
func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		at, err := GetAccessToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid access token.",
			})
		}

		jtwAdapter := sharedadapters.NewJWTAdapter(secret)

		claims, err := jtwAdapter.ValidateToken(c.UserContext(), at)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired authentication token",
			})
		}

		c.Locals("userId", claims.Subject)
		c.Locals("role", claims.Role)
		c.Locals("accessToken", at)

		return c.Next()
	}
}

func GetAccessToken(c *fiber.Ctx) (string, error) {
	if token := c.Cookies("ati"); token != "" {
		return token, nil
	}

	authHeader := c.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", errors.New("missing or invalid authorization header")
	}

	return parts[1], nil
}
