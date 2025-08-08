package permission

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app"
)

var (
    NeededRead  AccessControl = NeededPermissions{Owner, Write, Read}
    NeededWrite AccessControl = NeededPermissions{Owner, Write}
    NeededOwner AccessControl = NeededPermissions{Owner}
)

type AccessControl interface {
    HasNeededPermissions(permissions []Type) bool
}

type NeededPermissions []Type

func (nps NeededPermissions) HasNeededPermissions(permissions []Type) bool {
    for _, permission := range permissions {
        for _, neededPermission := range nps {
            if permission == neededPermission {
                return true
            }
        }
    }
    return false
}

type UserPermissionGetter interface {
    // Permission(ctx context.Context, user auth.User) (UserAccess, error)

    Permissions(ctx context.Context, user any) (UserAccess, error)
}

type ProjectIdExtractor func(c *app.RequestContext) (int, error)

func MiddlewareAccessControl(
        accessControl AccessControl,
        getProjectId ProjectIdExtractor,
        getUserPermissions UserPermissionGetter) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // get user from context

        // 通过 user 获取其拥有的权限 todo get user from context
        permissions, err := getUserPermissions.Permissions(ctx, nil)
        if err != nil {
            return
        }

        // 获取该 user 在 projectId 下的权限
        projectId, _ := getProjectId(c)
        projectPermissions := permissions.Permissions(projectId)

        if !accessControl.HasNeededPermissions(projectPermissions) {
            return
        }
    }
}
