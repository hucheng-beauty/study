package kani

import (
    "context"
    "testing"

    "code.byted.org/lab-speech/saas_backend/auth"
    "code.byted.org/lab-speech/saas_backend/commons/more_testing"
    "code.byted.org/lab-speech/saas_backend/permissions"
)

type testClient struct {
    t *testing.T
}

func (t testClient) CreateResource(_ context.Context, resource Resource) error {
    expectedResource := Resource{
        Name:               "project-123",
        EnglishName:        "project-123",
        Key:                "project-123",
        Creator:            "fake-creator",
        TemporaryAuthorize: false,
        OwnerKeys:          []string{"fake-owner-1", "fake-owner-2"},
        Confidentiality:    "L2",
        ProductName:        "AI-Lab智能语音",
        Actions: map[string]string{
            string(permissions.PermissionTypeRead):  "Read-only",
            string(permissions.PermissionTypeWrite): "Operation",
        },
    }
    more_testing.ErrorIfNotEqual(t.t, "create resource mismatch", expectedResource, resource)
    return nil
}

func (t testClient) ResourcePermissions(_ context.Context, resourceKey string) (ResourcePermissions, error) {
    more_testing.ErrorIfNotEqual(t.t, "resource key mismatch", "project-123", resourceKey)
    return ResourcePermissions{
        Status:       "created",
        AppID:        123,
        ResourceKey:  "project-123",
        ResourceName: "project-123",
        PermissionTypes: map[string]PermissionType{
            "admin": {"admin", "Owner"},
            "read":  {"read", "Read-Only"},
        },
        Permissions: map[string][]Permission{
            "admin": {{
                Email:       "tester@bytedance.com",
                Name:        "tester",
                ObjectType:  "user",
                Temporality: "permanent",
            }},
            "read": {{
                Email:       "tester1@bytedance.com",
                Name:        "tester1",
                ObjectType:  "user",
                Temporality: "permanent",
            }},
        },
    }, nil
}

func (t testClient) UserPermissions(_ context.Context, _ auth.User) ([]PermissionEntry, error) {
    return []PermissionEntry{
        {
            ResourceKey:         "project-123",
            ResourcePermissions: []string{"admin", "write", "read"},
        },
        {
            ResourceKey:         "project-124",
            ResourcePermissions: []string{"read"},
        },
    }, nil
}

func (t testClient) AddPermissions(_ context.Context, reqBody ModifyPermissionRequest) error {
    expectedReqBody := ModifyPermissionRequest{
        ResourceKey:   "project-123",
        Actions:       []string{"read", "write"},
        UserKeys:      []string{"tester"},
        AuthorizeDays: DefaultAuthorizeDays,
    }
    more_testing.ErrorIfNotEqual(t.t, "modify request mismatch", expectedReqBody, reqBody)
    return nil
}

func (t testClient) DeletePermissions(_ context.Context, _ ModifyPermissionRequest) error {
    panic("not supported")
}

func (t testClient) FindUser(_ context.Context, _ string) ([]UserQueryResult, error) {
    panic("not supported")
}

type testOwnershipDecider struct{}

func (testOwnershipDecider) ParcelOwnership(context.Context) (permissions.ParcelOwnership, error) {
    return permissions.ParcelOwnership{Creator: "fake-creator", Owners: []string{"fake-owner-1", "fake-owner-2"}}, nil
}

func TestPermissionManager_InitProjectPermissions(t *testing.T) {
    projectID := 123
    manager := PermissionManager{client: testClient{t: t}, ownershipDecider: testOwnershipDecider{}}
    err := manager.InitProjectPermissions(context.Background(), projectID)
    more_testing.FailOnError(t, "unexpected error", err)
}

func TestPermissionManager_AddPermissions(t *testing.T) {
    projectID := 123
    user := auth.User{Username: "tester"}
    perms := []permissions.PermissionType{permissions.PermissionTypeRead, permissions.PermissionTypeWrite}
    tests := []struct {
        name      string
        days      int
        errorMess string
        errorFunc func(t *testing.T, errorMsg string, err error)
    }{
        {
            name:      "Correct permission adding",
            days:      DefaultAuthorizeDays,
            errorMess: "unexpected error",
            errorFunc: more_testing.FailOnError,
        },
        {
            name:      "Correct permission adding",
            days:      0,
            errorMess: "unexpected error",
            errorFunc: more_testing.FailOnError,
        },
        {
            name:      "Wrong authorize days",
            days:      99999999,
            errorMess: "expected error",
            errorFunc: more_testing.FailWithoutError,
        },
    }

    manager := PermissionManager{client: testClient{t: t}}
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            ctx := context.Background()
            err := manager.AddPermissions(ctx, projectID, user, perms, test.days)
            test.errorFunc(t, test.errorMess, err)
        })
    }
}

func TestPermissionManager_ProjectPermissions(t *testing.T) {
    projectID := 123
    expectedPermissions := permissions.ProjectPermissions{
        PermissionTypes: []permissions.Descriptor{
            {Name: "Read-Only", Type: permissions.PermissionTypeRead},
            {Name: "Owner", Type: permissions.PermissionTypeOwner},
        },
        PermissionOwners: permissions.PermissionHolders{
            {
                User:   permissions.User{Name: "tester", Email: "tester@bytedance.com"},
                Target: "user",
                Permissions: permissions.Permissions{
                    {
                        Descriptor: permissions.Descriptor{
                            Name: "Owner",
                            Type: permissions.PermissionTypeOwner,
                        },
                        Temporality: "permanent",
                    },
                },
            },
            {
                User:   permissions.User{Name: "tester1", Email: "tester1@bytedance.com"},
                Target: "user",
                Permissions: permissions.Permissions{
                    {
                        Descriptor: permissions.Descriptor{
                            Name: "Read-Only",
                            Type: permissions.PermissionTypeRead,
                        },
                        Temporality: "permanent",
                    },
                },
            },
        },
    }

    ctx := context.Background()
    manager := PermissionManager{client: testClient{t: t}}
    perms, err := manager.ProjectPermissions(ctx, projectID)
    more_testing.FailOnError(t, "unexpected error", err)
    more_testing.FailIfNotEqual(t, "project permission mismatch", expectedPermissions, perms)
}

func TestPermissionManager_UserPermissions(t *testing.T) {
    expectedUserAccess := permissions.UserAccess{
        {
            ProjectID: 123,
            UserPermissions: []permissions.PermissionType{
                permissions.PermissionTypeOwner,
                permissions.PermissionTypeWrite, permissions.PermissionTypeRead,
            },
        },
        {
            ProjectID:       124,
            UserPermissions: []permissions.PermissionType{permissions.PermissionTypeRead},
        },
    }

    ctx := auth.CtxWithUser(context.Background(), auth.User{})
    manager := PermissionManager{client: testClient{t: t}}
    userAccess, err := manager.UserAccess(ctx)
    more_testing.FailOnError(t, "unexpected error", err)
    more_testing.FailIfNotEqual(t, "user permissions mismatch", expectedUserAccess, userAccess)
}
