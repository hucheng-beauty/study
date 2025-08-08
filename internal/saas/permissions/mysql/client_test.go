package mysql

import (
    "testing"

    // "code.byted.org/lab-speech/saas_backend/commons/more_testing"
    // "code.byted.org/lab-speech/saas_backend/permissions"

    "study/internal/saas/permissions"
)

func TestUserPermissions_ProjectPermissions(t *testing.T) {
    parcelID := 42
    adminDescriptor := "project-42_admin"
    tests := []struct {
        name                      string
        projectPermissions        ProjectPermissions            // when having this ProjectPermissions.
        expectedPermissionHolders permissions.PermissionHolders // then we expect this permissions.PermissionHolders.
    }{
        {
            name:                      "empty permissions",
            projectPermissions:        ProjectPermissions{},
            expectedPermissionHolders: nil,
        },
        {
            name: "single permission",
            projectPermissions: ProjectPermissions{
                {
                    Name:           "tester",
                    Email:          "tester@EMAIL",
                    PermissionType: "admin",
                },
            },
            expectedPermissionHolders: permissions.PermissionHolders{
                {
                    User: permissions.User{
                        Name:  "tester",
                        Email: "tester@EMAIL",
                    },
                    Target: "user",
                    Permissions: permissions.Permissions{
                        {
                            Descriptor: permissions.Descriptor{
                                Type: permissions.PermissionTypeOwner,
                                Name: adminDescriptor,
                            },
                            Temporality: "permanent",
                        },
                    },
                },
            },
        },
        {
            name: "multiple permissions",
            projectPermissions: ProjectPermissions{
                {
                    Name:           "tester",
                    Email:          "tester@EMAIL",
                    PermissionType: "admin",
                },
                {
                    Name:           "tester1",
                    Email:          "tester1@EMAIL",
                    PermissionType: "read",
                },
                {
                    Name:           "tester1",
                    Email:          "tester1@EMAIL",
                    PermissionType: "write",
                },
                {
                    Name:           "tester2",
                    Email:          "tester2@EMAIL",
                    PermissionType: "write",
                },
            },
            expectedPermissionHolders: permissions.PermissionHolders{
                {
                    User: permissions.User{
                        Name:  "tester",
                        Email: "tester@EMAIL",
                    },
                    Target: "user",
                    Permissions: permissions.Permissions{
                        {
                            Descriptor: permissions.Descriptor{
                                Type: permissions.PermissionTypeOwner,
                                Name: adminDescriptor,
                            },
                            Temporality: "permanent",
                        },
                    },
                },
                {
                    User: permissions.User{
                        Name:  "tester1",
                        Email: "tester1@EMAIL",
                    },
                    Target: "user",
                    Permissions: permissions.Permissions{
                        {
                            Descriptor: permissions.Descriptor{
                                Type: permissions.PermissionTypeRead,
                                Name: "Read-only",
                            },
                            Temporality: "temporary",
                        },
                        {
                            Descriptor: permissions.Descriptor{
                                Type: permissions.PermissionTypeWrite,
                                Name: "Operation",
                            },
                            Temporality: "temporary",
                        },
                    },
                },
                {
                    User: permissions.User{
                        Name:  "tester2",
                        Email: "tester2@EMAIL",
                    },
                    Target: "user",
                    Permissions: permissions.Permissions{
                        {
                            Descriptor: permissions.Descriptor{
                                Type: permissions.PermissionTypeWrite,
                                Name: "Operation",
                            },
                            Temporality: "temporary",
                        },
                    },
                },
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            permissionHolders := test.projectPermissions.PermissionHolders(parcelID)
            more_testing.ErrorIfNotEqual(t, "permission holders mismatch",
                test.expectedPermissionHolders, permissionHolders)
        })
    }
}

func TestUserAccess_UserAccess(t *testing.T) {
    tests := []struct {
        name                          string
        userAccess                    UserAccess             // when having this PermissionsUserAccess.
        expectedPermissionsUserAccess permissions.UserAccess // then we expect this permissions.PermissionsUserAccess.
    }{
        {
            name:                          "empty user access",
            userAccess:                    UserAccess{},
            expectedPermissionsUserAccess: nil,
        },
        {
            name: "single user access",
            userAccess: UserAccess{
                {
                    ProjectID:  32,
                    Permission: "admin",
                },
            },
            expectedPermissionsUserAccess: permissions.UserAccess{
                {
                    ProjectID:       32,
                    UserPermissions: []permissions.PermissionType{permissions.PermissionTypeOwner},
                },
            },
        },
        {
            name: "multiple user access",
            userAccess: UserAccess{
                {
                    ProjectID:  32,
                    Permission: "read",
                },
                {
                    ProjectID:  420,
                    Permission: "read",
                },
                {
                    ProjectID:  32,
                    Permission: "write",
                },
            },
            expectedPermissionsUserAccess: permissions.UserAccess{
                {
                    ProjectID: 32,
                    UserPermissions: []permissions.PermissionType{
                        permissions.PermissionTypeRead,
                        permissions.PermissionTypeWrite,
                    },
                },
                {
                    ProjectID:       420,
                    UserPermissions: []permissions.PermissionType{permissions.PermissionTypeRead},
                },
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            permissionsUserAccess := test.userAccess.PermissionsUserAccess()
            more_testing.ErrorIfNotEqual(t, "permissions user access mismatch",
                test.expectedPermissionsUserAccess, permissionsUserAccess)
        })
    }
}
