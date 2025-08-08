package mysql

import (
    "context"
    "testing"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons/more_testing"
    // "code.byted.org/lab-speech/saas_backend/permissions"

    "study/internal/saas/permissions"
)

type stubPermissionsClient struct {
    userPermittedProjectsFunc func(ctx context.Context, user auth.User) (UserAccess, error)
    projectPermissionsFunc    func(ctx context.Context, projectID int) (ProjectPermissions, error)
    deletePermissionFunc      func(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error
    addPermissionFunc         func(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error
}

func (spc stubPermissionsClient) UserPermittedProjects(ctx context.Context, user auth.User) (UserAccess, error) {
    return spc.userPermittedProjectsFunc(ctx, user)
}

func (spc stubPermissionsClient) ProjectPermissions(ctx context.Context, projectID int) (ProjectPermissions, error) {
    return spc.projectPermissionsFunc(ctx, projectID)
}

func (spc stubPermissionsClient) DeletePermission(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error {
    return spc.deletePermissionFunc(ctx, projectID, user, perm)
}

func (spc stubPermissionsClient) AddPermission(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error {
    return spc.addPermissionFunc(ctx, projectID, user, perm)
}

func (spc stubPermissionsClient) ProjectsPermissions(ctx context.Context, projectIDs []int) (map[int]*permissions.ProjectPermissions, error) {
    panic("TODO")
}

func Test_DBPermissionsGuard_UserAccess(t *testing.T) {
    // given
    expectedUser := auth.User{Username: "fake.username", Name: "Fake User"}
    permissionsClient := stubPermissionsClient{
        userPermittedProjectsFunc: func(ctx context.Context, user auth.User) (UserAccess, error) {
            more_testing.ErrorIfNotEqual(t, "invalid user permissions check", expectedUser, user)
            return UserAccess{
                {ProjectID: 2137, Permission: "fake-permission"},
                {ProjectID: 2137, Permission: "another-permission"},
                {ProjectID: 420, Permission: "another-project"},
            }, nil
        },
    }
    ctx := auth.CtxWithUser(context.Background(), expectedUser)
    // when
    guard := NewDBPermissionsGuard(permissionsClient, nil)
    userAccess, err := guard.UserAccess(ctx)
    // then
    more_testing.FailOnError(t, "", err)
    expectedUserAccess := permissions.UserAccess{
        {
            ProjectID:       2137,
            UserPermissions: []permissions.PermissionType{"fake-permission", "another-permission"},
        },
        {
            ProjectID:       420,
            UserPermissions: []permissions.PermissionType{"another-project"},
        },
    }
    more_testing.ErrorIfNotEqual(t, "invalid user access returned", expectedUserAccess, userAccess)
}

func Test_DBPermissionsGuard_ProjectPermissions(t *testing.T) {
    // given
    permissionsClient := stubPermissionsClient{
        projectPermissionsFunc: func(ctx context.Context, projectID int) (ProjectPermissions, error) {
            more_testing.ErrorIfNotEqual(t, "invalid projectID", 2137, projectID)
            return ProjectPermissions{
                {Username: "user-1", PermissionType: "read"},
                {Username: "user-2", PermissionType: "admin"},
                {Username: "user-1", PermissionType: "write"},
            }, nil
        },
    }
    guard := NewDBPermissionsGuard(permissionsClient, nil)
    // when
    projectPermissions, err := guard.ProjectPermissions(context.Background(), 2137)
    // then
    more_testing.FailOnError(t, "get project permissions", err)
    expectedPermissions := permissions.ProjectPermissions{
        PermissionTypes: permissions.Descriptors{
            {Type: permissions.PermissionTypeRead, Name: "Read-only"},
            {Type: permissions.PermissionTypeWrite, Name: "Operation"},
            {Type: permissions.PermissionTypeOwner, Name: "project-2137_admin"},
        },
        PermissionOwners: permissions.PermissionHolders{
            {
                User:   permissions.User{Name: "user-2", Email: "user-2@bytedance.com"},
                Target: "user",
                Permissions: permissions.Permissions{{
                    Descriptor:  permissions.Descriptor{Type: "admin", Name: "project-2137_admin"},
                    Temporality: "permanent",
                    Expires:     0,
                }},
            },
            {
                User:   permissions.User{Name: "user-1", Email: "user-1@bytedance.com"},
                Target: "user",
                Permissions: permissions.Permissions{
                    {
                        Descriptor:  permissions.Descriptor{Type: "write", Name: "Operation"},
                        Temporality: "temporary",
                        Expires:     0,
                    },
                    {
                        Descriptor:  permissions.Descriptor{Type: "read", Name: "Read-only"},
                        Temporality: "temporary",
                        Expires:     0,
                    },
                },
            },
        },
    }
    more_testing.ErrorIfNotEqual(t, "invalid project permissions", expectedPermissions, projectPermissions)
}

func Test_DBPermissionsGuard_FindUser(t *testing.T) {
    guard := NewDBPermissionsGuard(nil, nil)
    foundUsers, err := guard.FindUser(context.Background(), "some query")
    more_testing.FailOnError(t, "", err)
    expectedUsers := []auth.User{}
    more_testing.ErrorIfNotEqual(t, "invalid users found (expected none)", expectedUsers, foundUsers)
}

type changedPermission struct {
    projectID  int
    user       auth.User
    permission permissions.PermissionType
}

type staticOwnershipDecider struct {
    parcelOwnership permissions.ParcelOwnership
}

func (sod staticOwnershipDecider) ParcelOwnership(context.Context) (permissions.ParcelOwnership, error) {
    return sod.parcelOwnership, nil
}

func Test_DBPermissionsGuard_InitProjectPermissions(t *testing.T) {
    // given
    var addedPermissions []changedPermission
    permissionsClient := stubPermissionsClient{
        addPermissionFunc: func(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error {
            newPermission := changedPermission{projectID: projectID, user: user, permission: perm}
            addedPermissions = append(addedPermissions, newPermission)
            return nil
        },
    }
    ownershipDecider := staticOwnershipDecider{
        parcelOwnership: permissions.ParcelOwnership{Creator: "creator-user", Owners: []string{"owner-1", "owner-2"}},
    }
    guard := NewDBPermissionsGuard(permissionsClient, ownershipDecider)
    // when
    err := guard.InitProjectPermissions(context.Background(), 2137)
    // then
    more_testing.FailOnError(t, "", err)
    expectedPermissions := []changedPermission{
        {projectID: 2137, user: auth.User{Username: "owner-1"}, permission: permissions.PermissionTypeOwner},
        {projectID: 2137, user: auth.User{Username: "owner-2"}, permission: permissions.PermissionTypeOwner},
        {projectID: 2137, user: auth.User{Username: "creator-user"}, permission: permissions.PermissionTypeOwner},
    }
    more_testing.ErrorIfNotEqual(t, "invalid permissions added", expectedPermissions, addedPermissions)
}

func Test_DBPermissionsGuard_DeletePermissions(t *testing.T) {
    // given
    var deletedPermissions []changedPermission
    permissionsClient := stubPermissionsClient{
        deletePermissionFunc: func(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error {
            deletedPermission := changedPermission{projectID: projectID, user: user, permission: perm}
            deletedPermissions = append(deletedPermissions, deletedPermission)
            return nil
        },
    }
    guard := NewDBPermissionsGuard(permissionsClient, nil)
    // when
    err := guard.DeletePermissions(context.Background(), 2137, auth.User{Username: "fake-user"},
        []permissions.PermissionType{permissions.PermissionTypeWrite, permissions.PermissionTypeRead})
    // then
    more_testing.FailOnError(t, "", err)
    expectedChange := []changedPermission{
        {projectID: 2137, user: auth.User{Username: "fake-user"}, permission: permissions.PermissionTypeWrite},
        {projectID: 2137, user: auth.User{Username: "fake-user"}, permission: permissions.PermissionTypeRead},
    }
    more_testing.ErrorIfNotEqual(t, "invalid permissions deleted", expectedChange, deletedPermissions)
}

func Test_DBPermissionsGuard_AddPermissions(t *testing.T) {
    // given
    var addedPermissions []changedPermission
    permissionsClient := stubPermissionsClient{
        addPermissionFunc: func(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error {
            addedPermission := changedPermission{projectID: projectID, user: user, permission: perm}
            addedPermissions = append(addedPermissions, addedPermission)
            return nil
        },
    }
    guard := NewDBPermissionsGuard(permissionsClient, nil)
    // when
    err := guard.AddPermissions(context.Background(), 2137, auth.User{Username: "fake-user"},
        []permissions.PermissionType{permissions.PermissionTypeWrite, permissions.PermissionTypeRead}, 0)
    // then
    more_testing.FailOnError(t, "", err)
    expectedChange := []changedPermission{
        {projectID: 2137, user: auth.User{Username: "fake-user"}, permission: permissions.PermissionTypeWrite},
        {projectID: 2137, user: auth.User{Username: "fake-user"}, permission: permissions.PermissionTypeRead},
    }
    more_testing.ErrorIfNotEqual(t, "invalid permissions added", expectedChange, addedPermissions)
}
