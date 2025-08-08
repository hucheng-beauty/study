package mysql

import (
    "context"
    "fmt"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons"
    // "code.byted.org/lab-speech/saas_backend/permissions"

    "study/internal/saas/permissions"
)

type DBPermissionsClient interface {
    UserPermittedProjects(ctx context.Context, user auth.User) (UserAccess, error)
    ProjectPermissions(ctx context.Context, projectID int) (ProjectPermissions, error)
    DeletePermission(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error
    AddPermission(ctx context.Context, projectID int, user auth.User, perm permissions.PermissionType) error
}

type DBPermissionsGuard struct {
    client           DBPermissionsClient
    ownershipDecider permissions.ParcelOwnershipDecider
}

func NewDBPermissionsGuard(client DBPermissionsClient,
        ownershipDecider permissions.ParcelOwnershipDecider) DBPermissionsGuard {
    return DBPermissionsGuard{client: client, ownershipDecider: ownershipDecider}
}

func (db DBPermissionsGuard) UserAccess(ctx context.Context) (permissions.UserAccess, error) {
    user, err := auth.UserFromCtx(ctx)
    if err != nil {
        return nil, fmt.Errorf("get user from context: %w", err)
    }
    userAccess, err := db.client.UserPermittedProjects(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("permitted projects for user %s: %w", user.Username, err)
    }
    return userAccess.PermissionsUserAccess(), nil
}

func (db DBPermissionsGuard) UserAccessWithName(ctx context.Context, userName string) (permissions.UserAccess, error) {
    userAccess, err := db.client.UserPermittedProjects(ctx, auth.User{Username: userName})
    if err != nil {
        return nil, fmt.Errorf("permitted projects for user %s: %w", userName, err)
    }
    return userAccess.PermissionsUserAccess(), nil
}

func (db DBPermissionsGuard) ProjectPermissions(ctx context.Context,
        projectID int) (permissions.ProjectPermissions, error) {
    descriptors := db.projectPermissionDescriptors(projectID)
    projectPermissions, err := db.client.ProjectPermissions(ctx, projectID)
    if err != nil {
        return permissions.ProjectPermissions{}, fmt.Errorf("project permissions: %w", err)
    }
    permissionHolders := projectPermissions.PermissionHolders(projectID)
    ppp := permissions.ProjectPermissions{
        PermissionTypes:  descriptors,
        PermissionOwners: permissionHolders,
    }
    ppp.Sort()
    return ppp, nil
}

func (db DBPermissionsGuard) ProjectsPermissions(ctx context.Context, projectIDs []int) (map[int]*permissions.ProjectPermissions, error) {
    res := map[int]*permissions.ProjectPermissions{}
    for _, pr := range projectIDs {
        perm, err := db.ProjectPermissions(ctx, pr)
        if err != nil {
            return nil, err
        }
        res[pr] = &perm
    }
    return res, nil
}

func (db DBPermissionsGuard) FindUser(_ context.Context, _ string) ([]auth.User, error) {
    return make([]auth.User, 0), nil
}

func (db DBPermissionsGuard) projectPermissionDescriptors(projectID int) permissions.Descriptors {
    return permissions.Descriptors{
        {
            Type: permissions.PermissionTypeOwner,
            Name: fmt.Sprintf("project-%d_admin", projectID),
        },
        {
            Type: permissions.PermissionTypeWrite,
            Name: "Operation",
        },
        {
            Type: permissions.PermissionTypeRead,
            Name: "Read-only",
        },
    }
}

func (db DBPermissionsGuard) InitProjectPermissions(ctx context.Context, projectID int) error {
    parcelOwnership, err := db.ownershipDecider.ParcelOwnership(ctx)
    if err != nil {
        return fmt.Errorf("get initial parcel ownership: %w", err)
    }
    owners := commons.UniqueStringSlice(append(parcelOwnership.Owners, parcelOwnership.Creator))
    return db.addPermissions(ctx, owners, projectID)
}

func (db DBPermissionsGuard) InitProjectPermissionsExplicit(ctx context.Context, projectID int,
        creator string, owners []string) error {
    usernames := commons.UniqueStringSlice(append(owners, creator))
    return db.addPermissions(ctx, usernames, projectID)
}

func (db DBPermissionsGuard) addPermissions(ctx context.Context, owners []string, projectID int) error {
    perms := []permissions.PermissionType{permissions.PermissionTypeOwner}
    for _, owner := range owners {
        // HACK(jakub.daliga): We are passing a handcrafted auth.User, because we know only Username is required.
        err := db.AddPermissions(ctx, projectID, auth.User{Username: owner}, perms, -1)
        if err != nil {
            return fmt.Errorf("add permission: %w", err)
        }
    }
    return nil
}

func (db DBPermissionsGuard) DeletePermissions(ctx context.Context, projectID int,
        user auth.User, permission []permissions.PermissionType) error {
    for _, perm := range permission {
        err := db.client.DeletePermission(ctx, projectID, user, perm)
        if err != nil {
            return fmt.Errorf("delete permission: %w", err)
        }
    }
    return nil
}

func (db DBPermissionsGuard) AddPermissions(ctx context.Context, projectID int, user auth.User,
        newPermissions []permissions.PermissionType, _ int) error {
    for _, perm := range newPermissions {
        err := db.client.AddPermission(ctx, projectID, user, perm)
        if err != nil {
            return fmt.Errorf("add permission: %w", err)
        }
    }
    return nil
}
