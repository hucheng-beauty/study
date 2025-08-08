package permissions

import (
    "context"
    "fmt"
    "net/http"
    "testing"

    // "code.byted.org/middleware/hertz/pkg/app"
    // "code.byted.org/middleware/hertz/pkg/app/server"
    // "code.byted.org/middleware/hertz_ext/v2/hertztest"
    //
    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons/errutil"
    // "code.byted.org/lab-speech/saas_backend/commons/more_testing"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
)

func mockProjectIDExtractor(_ *app.RequestContext) (int, error) {
    return 42, nil
}

type mockUserPermissionGetter struct {
    perms []PermissionType
}

func (m mockUserPermissionGetter) UserAccess(_ context.Context, _ auth.User) (UserAccess, error) {
    return UserAccess{
        UserAccessControl{
            ProjectID:       42,
            UserPermissions: m.perms,
        },
    }, nil
}

type mockUserPermissionGetterError struct{}

func (m mockUserPermissionGetterError) UserAccess(_ context.Context, _ auth.User) (UserAccess, error) {
    return nil, fmt.Errorf("terrible error occurred")
}

func TestMiddlewareAccessControl(t *testing.T) {
    tests := []struct {
        name             string
        code             int
        accessControl    AccessControl
        permissionGetter UserPermissionGetter
    }{
        {
            name:             "Correct authorize",
            code:             200,
            accessControl:    RequireOwner,
            permissionGetter: mockUserPermissionGetter{[]PermissionType{PermissionTypeOwner}},
        },
        {
            name:             "Greater permission present",
            code:             200,
            accessControl:    RequireRead,
            permissionGetter: mockUserPermissionGetter{[]PermissionType{PermissionTypeWrite}},
        },
        {
            name:             "Missing permission",
            code:             403,
            accessControl:    RequireOwner,
            permissionGetter: mockUserPermissionGetter{[]PermissionType{PermissionTypeRead}},
        },
        {
            name:             "Permission getter error",
            code:             500,
            accessControl:    RequireOwner,
            permissionGetter: mockUserPermissionGetterError{},
        },
    }

    errutil.FullErrorInHttpResponse = true
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            userInserter := app.HandlerFunc(func(ctx context.Context, c *app.RequestContext) {
                ctx = auth.CtxWithUser(ctx, auth.User{})
                c.Next(ctx)
            })
            accessDecider := MiddlewareAccessControl(mockProjectIDExtractor, test.accessControl, test.permissionGetter)
            router := server.Default()
            router.GET("/", userInserter, accessDecider)
            resp := hertztest.PerformRequest(router.Engine, http.MethodGet, "/", nil)
            more_testing.ErrorIfNotEqual(t, "status code mismatch", test.code, resp.Code)
        })
    }
}
