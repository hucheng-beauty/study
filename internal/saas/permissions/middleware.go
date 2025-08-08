package permissions

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons/errutil"
)

var (
    RequireRead  AccessControl = PermissionOneOf{PermissionTypeOwner, PermissionTypeWrite, PermissionTypeRead}
    RequireWrite AccessControl = PermissionOneOf{PermissionTypeOwner, PermissionTypeWrite}
    RequireOwner AccessControl = PermissionOneOf{PermissionTypeOwner}
)

type AccessControl interface {
    HasRequiredPermissions(types []PermissionType) bool
}

type UserPermissionGetter interface {
    // UserAccess(ctx context.Context, user auth.User) (UserAccess, error)

    UserAccess(ctx context.Context, user any) (UserAccess, error)
}

type ProjectIDExtractor func(c *app.RequestContext) (int, error)

func MiddlewareAccessControl(
        projectIDExtractor ProjectIDExtractor, accessControl AccessControl,
        permissionManager UserPermissionGetter) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        projectID, err := projectIDExtractor(c)
        if err != nil {
            // errutil.HandleErrorHertz(ctx, c, err)
            return
        }
        user, err := auth.UserFromCtx(ctx)
        if err != nil {
            // err := fmt.Errorf("user from context: %w", err)
            // errutil.HandleErrorHertz(ctx, c, err)
            return
        }
        permissions, err := permissionManager.UserAccess(ctx, user)
        if err != nil {
            // err = errors.Wrap(err, "get resource permissions")
            // errutil.HandleErrorHertz(ctx, c, err)
            return
        }
        permissionsForProject := permissions.PermissionsToProject(projectID)
        hasRequiredPermissions := accessControl.HasRequiredPermissions(permissionsForProject)
        if !hasRequiredPermissions {
            // err := errutil.FromUserMessage("you don't have required permissions")
            // err = errutil.WithHTTPStatusCode(err, http.StatusForbidden)
            // errutil.HandleErrorHertz(ctx, c, err)
            return
        }
    }
}

type PermissionOneOf []PermissionType

func (requiredPermissions PermissionOneOf) HasRequiredPermissions(ownedPermissions []PermissionType) bool {
    for _, ownedPermission := range ownedPermissions {
        for _, requiredPermission := range requiredPermissions {
            if ownedPermission == requiredPermission {
                return true
            }
        }
    }
    return false
}
