package middleware

// func PermissionCheck(c *fiber.Ctx) error {
// 	perm := strings.Join([]string{
// 		constant.BaseRole,
// 	}, ".")
// 	allowPerm := map[string]bool{
// 		constant.RolePermissionSuper:  true,
// 		constant.RolePermissionMember: false,
// 		constant.RolePermissionAdmin:  true,
// 	}
// 	if err := checkPermAndSetRole(c, perm, allowPerm); err != nil {
// 		return err
// 	}
// 	ctx := context.WithValue(c.UserContext(), constant.CtxPageName, constant.RolePagePcmDbContractSummaryFg)
// 	c.SetUserContext(ctx)
// 	return c.Next()
// }

// func checkPermAndSetRole(c *fiber.Ctx, perm string, allowRole map[string]bool) error {
// 	roles, ok := c.Locals("roles").([]string)
// 	if !ok {
// 		return exception.NewWithStatus(http.StatusBadRequest, "invalid roles type.")
// 	}
// 	foundPerm, err := getAllAllowPerm(roles, perm, allowRole)
// 	if err != nil {
// 		return err
// 	}
// 	var roleVal string
// 	if findSuper(foundPerm, allowRole) {
// 		roleVal = constant.RolePermissionSuper
// 	} else {
// 		roleVal = getRoleValFromPermString(foundPerm[0])
// 	}
// 	ctx := context.WithValue(c.UserContext(), constant.CtxPcmPermission, roleVal)
// 	c.SetUserContext(ctx)
// 	return nil
// }

// func mockCtx(c *fiber.Ctx) {
// 	ctx := context.WithValue(c.UserContext(), constant.CtxPcmPermission, constant.RolePermissionMember)
// 	c.SetUserContext(ctx)
// }

// func getAllAllowPerm(roles []string, perm string, allowRole map[string]bool) ([]string, error) {
// 	foundPerm := []string{}
// 	for _, r := range roles {
// 		role, hasPrefix := strings.CutPrefix(r, perm+".")
// 		role = strings.Split(role, ".")[0]
// 		if hasPrefix && allowRole[role] {
// 			foundPerm = append(foundPerm, r)
// 		}
// 	}
// 	if len(foundPerm) == 0 {
// 		return nil, exception.NewWithStatus(http.StatusForbidden, "access denied")
// 	}
// 	return foundPerm, nil
// }

// func findSuper(roles []string, allowRole map[string]bool) bool {
// 	for _, p := range roles {
// 		rp := strings.Split(p, ".")
// 		p := rp[len(rp)-1]
// 		if p == constant.RolePermissionSuper {
// 			return true
// 		}
// 	}
// 	return false
// }

// func getRoleValFromPermString(val string) string {
// 	v := strings.Split(val, ".")
// 	return v[len(v)-1]
// }
