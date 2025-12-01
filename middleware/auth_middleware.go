// package middleware

// import (
// 	"PrestasiMhs-API/utils"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// )

// // Protected melindungi route agar hanya bisa diakses user yang punya Token valid
// func Protected() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// 1. Ambil header Authorization
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token tidak ditemukan"})
// 		}

// 		// 2. Format harus "Bearer <token>"
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Format token salah"})
// 		}

// 		// 3. Validasi Token menggunakan Utils yang sudah kita update
// 		tokenString := parts[1]
// 		claims, err := utils.ValidateToken(tokenString)
// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token tidak valid atau kadaluarsa"})
// 		}

// 		// 4. Simpan UserID dan Role ke Context (agar bisa dipakai di Service nanti)
// 		c.Locals("user_id", claims.UserID)
// 		c.Locals("role", claims.Role)

// 		return c.Next()
// 	}
// }

// // RoleMiddleware membatasi akses berdasarkan Role (Contoh: hanya Admin)
// func RoleMiddleware(allowedRoles ...string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		userRole := c.Locals("role").(string)

// 		for _, role := range allowedRoles {
// 			if role == userRole {
// 				return c.Next()
// 			}
// 		}

// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Akses ditolak: Role Anda tidak diizinkan",
// 		})
// 	}
// }

package middleware

import (
	"PrestasiMhs-API/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token tidak ditemukan"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Format token salah"})
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token tidak valid atau kadaluarsa"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)

		for _, role := range allowedRoles {
			if role == userRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses ditolak: Role Anda tidak diizinkan",
		})
	}
}