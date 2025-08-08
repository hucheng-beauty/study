package kani

import (
    "context"
    "fmt"
    "net/http"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons/errutil"
    // "code.byted.org/lab-speech/saas_backend/permissions"

    "study/internal/saas/permissions"
)

const permissionTarget string = "user"

type Client interface {
    CreateResource(ctx context.Context, resource Resource) error
    ResourcePermissions(ctx context.Context, resourceKey string) (ResourcePermissions, error)
    UserPermissions(ctx context.Context, user auth.User) ([]PermissionEntry, error)
    AddPermissions(ctx context.Context, reqBody ModifyPermissionRequest) error
    DeletePermissions(ctx context.Context, reqBody ModifyPermissionRequest) error
    FindUser(ctx context.Context, query string) ([]UserQueryResult, error)
}

type PermissionManager struct {
    client           Client
    ownershipDecider permissions.ParcelOwnershipDecider
}

func NewPermissionManager(kaniClient Client, ownershipDecider permissions.ParcelOwnershipDecider) PermissionManager {
    return PermissionManager{client: kaniClient, ownershipDecider: ownershipDecider}
}

func (p PermissionManager) UserAccess(ctx context.Context) (permissions.UserAccess, error) {
    user, err := auth.UserFromCtx(ctx)
    if err != nil {
        return nil, fmt.Errorf("get user from context: %w", err)
    }
    permissionEntries, err := p.client.UserPermissions(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("user raw permissions: %w", err)
    }
    userAccess, err := UserAccessControl(permissionEntries)
    if err != nil {
        return nil, fmt.Errorf("convert raw permissions to user access: %w", err)
    }
    return userAccess, nil
}

func (p PermissionManager) AddPermissions(ctx context.Context, projectID int, user auth.User,
        newPermissions []permissions.PermissionType, days int) error {
    resourceKey := ParcelResourceKey(projectID)
    var actions []string
    for _, newPermission := range newPermissions {
        actions = append(actions, string(newPermission))
    }
    reqBody := ModifyPermissionRequest{
        ResourceKey:   resourceKey,
        Actions:       actions,
        UserKeys:      []string{user.Username},
        AuthorizeDays: days,
    }
    // TODO(jakub.daliga) check if kani allows us to send authorize days with delete(it's not listed)
    if reqBody.AuthorizeDays == 0 {
        reqBody.AuthorizeDays = DefaultAuthorizeDays
    }
    if reqBody.AuthorizeDays > MaxAuthorizeDays {
        errMsg := fmt.Sprintf("request exceeds maximum permission duration (%v days)", MaxAuthorizeDays)
        return errutil.WithHTTPStatusCode(
            errutil.FromUserMessage(errMsg),
            http.StatusBadRequest,
        )
    }
    err := p.client.AddPermissions(ctx, reqBody)
    if err != nil {
        return fmt.Errorf("add permissions: %w", err)
    }
    return nil
}

func (p PermissionManager) DeletePermissions(ctx context.Context, projectID int, user auth.User,
        permission []permissions.PermissionType) error {
    resourceKey := ParcelResourceKey(projectID)
    var actions []string
    for _, newPermission := range permission {
        actions = append(actions, string(newPermission))
    }
    reqBody := ModifyPermissionRequest{
        ResourceKey: resourceKey,
        Actions:     actions,
        UserKeys:    []string{user.Username},
    }
    err := p.client.DeletePermissions(ctx, reqBody)
    if err != nil {
        return fmt.Errorf("perform modify permission request: %w", err)
    }
    return nil
}

func (p PermissionManager) ProjectPermissions(ctx context.Context,
        projectID int) (permissions.ProjectPermissions, error) {
    resourceKey := ParcelResourceKey(projectID)
    resourcePermissions, err := p.client.ResourcePermissions(ctx, resourceKey)
    if err != nil {
        return permissions.ProjectPermissions{}, fmt.Errorf("get resource permissions: %w", err)
    }
    return p.convertResourcePermissionToProjectPermissions(resourcePermissions), nil
}

func (p PermissionManager) convertResourcePermissionToProjectPermissions(rp ResourcePermissions) permissions.ProjectPermissions {
    ppp := permissions.ProjectPermissions{}
    ppp.PermissionTypes = p.descriptors(rp)
    ppp.PermissionOwners = p.permissionHolders(rp)
    ppp.Sort()
    return ppp
}

func (p PermissionManager) descriptors(rp ResourcePermissions) permissions.Descriptors {
    var descriptors permissions.Descriptors
    for _, permissionType := range rp.PermissionTypes {
        descriptors = append(descriptors,
            permissions.Descriptor{
                Name: permissionType.Name,
                Type: permissions.PermissionType(permissionType.Action),
            })
    }
    return descriptors
}

func (p PermissionManager) permissionHolders(rp ResourcePermissions) permissions.PermissionHolders {
    var permissionHolders permissions.PermissionHolders
    for action, actionPermissions := range rp.Permissions {
    actionLoop:
        for _, actionPermission := range actionPermissions {
            var expires int64
            if actionPermission.ExpireTime != nil {
                expires = actionPermission.ExpireTime.UnixNano() / 1000000
            }
            permission := permissions.Permission{
                Descriptor: permissions.Descriptor{
                    Type: permissions.PermissionType(action),
                    Name: rp.PermissionTypes[action].Name,
                },
                Temporality: actionPermission.Temporality,
                Expires:     expires,
            }
            for i, existingPermission := range permissionHolders {
                if actionPermission.Email != existingPermission.User.Email {
                    continue
                }
                // this actionPermission is for our user
                existingPermission.Permissions = append(existingPermission.Permissions, permission)
                permissionHolders[i] = existingPermission
                continue actionLoop
            }
            // user for this permission not found
            permissionHolders = append(permissionHolders, permissions.PermissionHolder{
                User: permissions.User{
                    Email: actionPermission.Email,
                    Name:  actionPermission.Name,
                },
                Target:      permissionTarget,
                Permissions: []permissions.Permission{permission},
            })
        }
    }
    return permissionHolders
}

// TODO(jakub.daliga): This should be bound to a separate entity - this is SaaS
//
//	business logic merged with KANI client logic.
func (p PermissionManager) InitProjectPermissions(ctx context.Context, projectID int) error {
    parcelOwnership, err := p.ownershipDecider.ParcelOwnership(ctx)
    if err != nil {
        return fmt.Errorf("get initial parcel ownership: %w", err)
    }

    resourceKey := ParcelResourceKey(projectID)
    resource := Resource{
        // required fields
        Name:               resourceKey,
        EnglishName:        resourceKey,
        Key:                resourceKey,
        TemporaryAuthorize: false,
        Creator:            parcelOwnership.Creator,
        OwnerKeys:          parcelOwnership.Owners,
        // NOTE(xinyu): The default value of "cnfidentiality" field is "normal";
        // but this value is deprecated and will result in error in i18n; "L2"
        // is equivalent to "normal" in terms of valid days of permissions given
        // to client; its value is hard-coded for now and is subject to change;
        // for more info, see:
        // https://ei.byted.org/ratak/doc/#api-Resource-Resources-Post
        // https://bytedance.feishu.cn/wiki/wikcn1wX1aY4PHSIBa8IZymWVGi#IDxH1Y
        // https://bytedance.feishu.cn/docs/doccniv4iKqoMQfnVTAkNUA1QMc
        Confidentiality: "L2",
        // NOTE(xinyu): The value is the "affiliated business" that this resource
        // belongs to, by definition; I filled with ai-lab-speech for now; it only
        // affects the "permission requiring" process when the confidentiality is
        // set to "L4", which is not suitable for our case; but this field is required;
        // for more info, see:
        // https://ei.byted.org/ratak/doc/#api-Resource-Resources-Post
        // https://bytedance.feishu.cn/docs/doccniv4iKqoMQfnVTAkNUA1QMc
        ProductName: "AI-Lab智能语音",

        // optional fields
        Actions: map[string]string{
            string(permissions.PermissionTypeRead):  "Read-only",
            string(permissions.PermissionTypeWrite): "Operation",
        },
    }
    err = p.client.CreateResource(ctx, resource)
    if err != nil {
        return fmt.Errorf("create resource: %w", err)
    }
    return nil
}

func (p PermissionManager) FindUser(ctx context.Context,
        queryUserSearch string) ([]auth.User, error) {
    matchedUsers, err := p.client.FindUser(ctx, queryUserSearch)
    if err != nil {
        return nil, fmt.Errorf("user search: %w", err)
    }
    var authUsers []auth.User
    for _, user := range matchedUsers {
        authUsers = append(authUsers, auth.User{Email: user.Email, Name: user.Name})
    }
    return authUsers, nil
}
