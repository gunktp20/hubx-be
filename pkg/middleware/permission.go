package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gunktp20/digital-hubx-be/pkg/constant"
)

func PermissionCheck(c *fiber.Ctx) error {
	// Base permission prefix
	permPrefix := constant.BaseRole

	// Get roles from payload
	roles, ok := c.Locals("roles").([]string)
	if !ok {
		log.Println("No roles found in payload")
		return fiber.NewError(fiber.StatusUnauthorized, "No roles found in payload")
	}

	// Define allowed permissions using constants
	allowPerm := map[string]bool{
		constant.RolePermissionUserRead:        false,
		constant.RolePermissionAdminFullAccess: true,
		constant.RolePermissionSuper:           true,
	}

	// Define role hierarchy (optional)
	roleHierarchy := map[string][]string{
		constant.RolePermissionAdminFullAccess: {
			constant.RolePermissionUserRead, // Admin has access to User permissions
		},
		constant.RolePermissionSuper: {
			constant.RolePermissionAdminFullAccess, // Super includes Admin and User permissions
			constant.RolePermissionUserRead,
		},
	}

	// Iterate through all roles and check permissions
	for _, role := range roles {
		// Ensure the role starts with the base permission prefix
		if strings.HasPrefix(role, permPrefix) {
			// Extract the permission name (e.g., "User.Read" from "Digital.X.HUB.User.Read")
			perm := strings.TrimPrefix(role, permPrefix+".")

			// Check explicit permission
			if isAllowed, exists := allowPerm[perm]; exists && isAllowed {
				log.Printf("Permission granted: %s\n", perm)
				ctx := context.WithValue(c.UserContext(), constant.CtxPageName, perm)
				c.SetUserContext(ctx)
				return c.Next()
			}

			// Check permission hierarchy
			if higherPerms, ok := roleHierarchy[perm]; ok {
				for _, higherPerm := range higherPerms {
					if isAllowed, exists := allowPerm[higherPerm]; exists && isAllowed {
						log.Printf("Permission granted through hierarchy: %s\n", higherPerm)
						ctx := context.WithValue(c.UserContext(), constant.CtxPageName, higherPerm)
						c.SetUserContext(ctx)
						return c.Next()
					}
				}
			}
		}
	}

	// If no role matches, return unauthorized
	log.Println("Access denied: insufficient permissions")
	return fiber.NewError(fiber.StatusForbidden, "Access denied: insufficient permissions")
}
