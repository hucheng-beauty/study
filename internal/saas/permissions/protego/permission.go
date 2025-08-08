package protego

import (
    "context"
    "errors"
    "fmt"
    "time"

    "study/internal/saas/permissions"
    "study/internal/saas/permissions/kani"

    // "code.byted.org/gopkg/logs/v2/log"
    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/clients/protego"
    // "code.byted.org/lab-speech/saas_backend/permissions"
    // "code.byted.org/lab-speech/saas_backend/permissions/kani"
    // "code.byted.org/security/authorization/lib/model/apis/kitex_gen/security/authorization/common"
    // "code.byted.org/security/authorization/svc/http/wrapper"
    "github.com/samber/lo"
)

// ErrUsername is returned when User.ErrUsername field is empty as input.
var ErrUsername = errors.New("username in auth.User is empty")

type Client interface {
    UserAccess(ctx context.Context) (permissions.UserAccess, error)
    AddPermissions(ctx context.Context, projectID int, user auth.User,
            newPermissions []permissions.PermissionType, days int) error
    DeletePermissions(ctx context.Context, projectID int, user auth.User,
            permission []permissions.PermissionType) error
    ProjectPermissions(ctx context.Context,
            projectID int) (permissions.ProjectPermissions, error)
    // ProjectsPermissions queries project permissions in batch.
    ProjectsPermissions(ctx context.Context, projectIDs []int) (map[int]*permissions.ProjectPermissions, error)
    // InitProjectPermissions uses `ParcelOwnershipDecider` to get the creator and owner
    // and create the corresponding resource point.
    InitProjectPermissions(ctx context.Context, projectID int) error
    // InitProjectPermissionsExplicit get creator and owners key from argument instead of
    // the `ParcelOwnershipDecider` like `InitProjectPermissions`.
    InitProjectPermissionsExplicit(ctx context.Context, projectID int, creator string, owners []string) error
    // FindUser finds user from his identity registered in protego; both UserName and Name
    // are returned as his unique name in company with suffix (eg. zhangteng.dance).
    FindUser(ctx context.Context,
            queryUserSearch string) ([]auth.User, error)
    UserAccessWithName(ctx context.Context, userName string) (permissions.UserAccess, error)
}

type httpClient struct {
    protegoClient    protego.Client
    ownershipDecider permissions.ParcelOwnershipDecider
}

func NewHTTPPermissionManager(client protego.Client, ownershipDecider permissions.ParcelOwnershipDecider) Client {
    return &httpClient{
        protegoClient:    client,
        ownershipDecider: ownershipDecider,
    }
}

func (h *httpClient) ProjectPermissions(ctx context.Context, projectID int) (permissions.ProjectPermissions, error) {
    perms, err := h.ProjectsPermissions(ctx, []int{projectID})
    if err != nil {
        return permissions.ProjectPermissions{}, err
    }
    perm, ok := perms[projectID]
    if !ok {
        return permissions.ProjectPermissions{}, fmt.Errorf("found no perm info for project: %d", projectID)
    }
    return *perm, nil
}

func (h *httpClient) ProjectsPermissions(ctx context.Context,
        projectIDs []int) (map[int]*permissions.ProjectPermissions, error) {
    body := &protego.PolicyOperationRequest{}
    for _, projectID := range projectIDs {
        resourceKey := kani.ParcelResourceKey(projectID)
        body.Policies = append(body.Policies, &wrapper.Policy{

            IdentityType: lo.ToPtr[int64](0),
            RoleType:     lo.ToPtr[int64](1),
            ResourceType: lo.ToPtr[int64](2),
            Resource: &wrapper.PolicyEntity{
                Ns:         lo.ToPtr(h.protegoClient.Namespace()),
                LocationV2: h.protegoClient.LocationV2(),
                PathV2:     []*wrapper.Pair{{Key: "key", Val: resourceKey}},
            },
        })
    }
    resp, err := h.protegoClient.PolicyRead(ctx, body)
    if err != nil {
        return nil, err
    }
    return h.projectsPermissions(resp), nil
}

// projectPermissions converts protego policy to the format inferred from kani.
func (h *httpClient) projectsPermissions(resp *protego.PolicyResponse) map[int]*permissions.ProjectPermissions {
    if resp == nil || len(resp.Data.Policies) == 0 {
        return nil
    }
    perms := map[int]*permissions.ProjectPermissions{}
    typeToPerms := map[int]permissions.Descriptors{}
    userToPerms := map[int]map[permissions.User]permissions.Permissions{}
    for _, policy := range resp.Data.Policies {
        permAction := ""
        permName := ""
        for _, path := range policy.Role.PathV2 {
            if path.Key == "action" {
                permAction = path.Val
            }
        }
        permNameAny, permNameFound := policy.Role.Attributes["name"]
        if permNameFound {
            permName, permNameFound = permNameAny.(string)
        }
        var parcelID int
        for _, path := range policy.Resource.PathV2 {
            if path.Key == "key" {
                if !permNameFound {
                    permName = path.Val + "_" + permAction
                }
                parcelID, _ = kani.ParcelIDFromResourceKey(path.Val)
            }
        }
        permDescriptor := permissions.Descriptor{
            Type: permissions.PermissionType(permAction),
            Name: permName,
        }
        typeToPerms[parcelID] = append(typeToPerms[parcelID], permDescriptor)
        // find resource / project
        // find users
        // check if policy is enabled and not expired; note that expire time from protego
        // is always in nano-second
        if expireTime := lo.FromPtr(policy.ExpireTime); policy.IsEnable != nil && *policy.IsEnable &&
                (expireTime == 0 || time.Now().Nanosecond() < int(expireTime)) {
            name, _ := policy.Identity.Attributes["name"].(string)
            email, _ := policy.Identity.Attributes["email"].(string)
            user := permissions.User{
                Name:  name,
                Email: email,
            }
            if userToPerms[parcelID] == nil {
                userToPerms[parcelID] = make(map[permissions.User]permissions.Permissions)
            }
            userToPerms[parcelID][user] = append(userToPerms[parcelID][user], permissions.Permission{
                Descriptor:  permDescriptor,
                Temporality: "permanent",
                Expires:     time.Unix(0, expireTime).UnixMilli(),
            })
        }
    }
    for projectID, userPerms := range userToPerms {
        perm := &permissions.ProjectPermissions{
            PermissionTypes: lo.Uniq(typeToPerms[projectID]),
        }
        for user, perms := range userPerms {
            perm.PermissionOwners = append(perm.PermissionOwners, permissions.PermissionHolder{
                User:        user,
                Target:      "user",
                Permissions: perms,
            })
        }
        perm.Sort()
        perms[projectID] = perm
    }
    return perms
}

func (h *httpClient) UserAccess(ctx context.Context) (permissions.UserAccess, error) {
    user, err := auth.UserFromCtx(ctx)
    if err != nil {
        return nil, fmt.Errorf("get user from context: %w", err)
    }
    return h.userAccessWithName(ctx, user.Username)
}

func (h *httpClient) UserAccessWithName(ctx context.Context, userName string) (permissions.UserAccess, error) {
    return h.userAccessWithName(ctx, userName)
}

func (h *httpClient) userAccessWithName(ctx context.Context, userName string) (permissions.UserAccess, error) {
    body := &protego.BatchAllowedEntityRequest{
        IsStrongConsistency: lo.ToPtr(true),
        IsAudit:             lo.ToPtr(true),
        ExactMatch:          lo.ToPtr(true),
        Tasks: []*wrapper.Task{
            {
                Principal: &wrapper.TaskEntity{
                    Ns:     "user",
                    PathV2: []*wrapper.Pair{{Key: "id", Val: userName}},
                },
                Object: &wrapper.TaskEntity{
                    Ns:         h.protegoClient.Namespace(),
                    LocationV2: h.protegoClient.LocationV2(),
                    Attributes: map[string]interface{}{"country": []string{}}, // show resources from all countries
                },
                SkipResourceOverlapFetch:    lo.ToPtr(true),
                ResourceReferenceNamespaces: []string{""},
                IdentityReferenceNamespaces: []string{
                    "group_" + h.protegoClient.Namespace(), "group_role_" + h.protegoClient.Namespace(), "group_ldap"},
            },
        },
    }
    resp, err := h.protegoClient.BatchAllowedEntity(ctx, body)
    if err != nil {
        return nil, err
    }
    if len(resp.Data.BatchResponses) == 0 ||
            len(resp.Data.BatchResponses[0].Workload.AllowedTuples) == 0 {
        return permissions.UserAccess{}, nil
    }
    return h.userAccess(ctx, resp.Data.BatchResponses[0].Workload.AllowedTuples), nil
}

func (h *httpClient) userAccess(ctx context.Context, policies []*wrapper.Tuple) permissions.UserAccess {
    userAccessRaw := lo.FilterMap(policies, func(pol *wrapper.Tuple, _ int) (permissions.UserAccessControl, bool) {
        if pol == nil {
            return permissions.UserAccessControl{}, false
        }
        // check if expired
        if pol.IsEnable != nil && !*pol.IsEnable {
            return permissions.UserAccessControl{}, false
        }
        if pol.ExpireTime != nil &&
                *pol.ExpireTime != 0 &&
                time.Unix(0, *pol.ExpireTime).Before(time.Now()) {
            return permissions.UserAccessControl{}, false
        }
        resourcePath := ""
        for _, p := range pol.Resource.PathV2 {
            if p.Key == "key" {
                resourcePath = p.Val
                break
            }
        }
        projectID, err := kani.ParcelIDFromResourceKey(resourcePath)
        if err != nil {
            log.V2.Error().With(ctx).Str("cannot find project from resource:", resourcePath, "policy:").Obj(pol).Emit()
            return permissions.UserAccessControl{}, false
        }
        rolePath := ""
        for _, p := range pol.Role.PathV2 {
            if p.Key == "action" {
                rolePath = p.Val
                break
            }
        }
        return permissions.UserAccessControl{
            ProjectID:       projectID,
            UserPermissions: []permissions.PermissionType{permissions.PermissionType(rolePath)},
        }, true
    })
    projectToPerm := map[int]permissions.UserAccessControl{}
    for _, perm := range userAccessRaw {
        _, ok := projectToPerm[perm.ProjectID]
        if ok {
            projectToPerm[perm.ProjectID] = permissions.UserAccessControl{
                ProjectID:       perm.ProjectID,
                UserPermissions: append(projectToPerm[perm.ProjectID].UserPermissions, perm.UserPermissions[0]),
            }
        } else {
            projectToPerm[perm.ProjectID] = permissions.UserAccessControl{
                ProjectID:       perm.ProjectID,
                UserPermissions: []permissions.PermissionType{perm.UserPermissions[0]},
            }
        }
    }
    return lo.Values(projectToPerm)
}

func (h *httpClient) AddPermissions(ctx context.Context, projectID int, user auth.User,
        newPermissions []permissions.PermissionType, days int) error {
    err := validateUsername(&user)
    if err != nil {
        return err
    }
    resourceKey := kani.ParcelResourceKey(projectID)
    policies := lo.Map(newPermissions, func(permType permissions.PermissionType, _ int) *wrapper.Policy {
        return h.policyWithUser(resourceKey, permType, user.Username,
            lo.ToPtr(time.Now().AddDate(0, 0, days).UnixNano()))
    })
    body := &protego.PolicyOperationRequest{
        Relation:      lo.ToPtr(common.HierarchyRelation_Equivalent),
        UpdateOnExist: lo.ToPtr(true),
        Policies:      policies,
    }
    _, err = h.protegoClient.PolicyCreate(ctx, body)
    if err != nil {
        return err
    }
    return nil
}

func (h *httpClient) DeletePermissions(ctx context.Context, projectID int, user auth.User,
        permission []permissions.PermissionType) error {
    err := validateUsername(&user)
    if err != nil {
        return err
    }
    resourceKey := kani.ParcelResourceKey(projectID)
    policies := lo.Map(permission, func(permType permissions.PermissionType, _ int) *wrapper.Policy {
        return h.policyWithUser(resourceKey, permType, user.Username, nil)
    })
    body := &protego.PolicyOperationRequest{
        Policies:   policies,
        BasePolicy: &wrapper.Policy{IsEnable: lo.ToPtr(false)},
    }
    _, err = h.protegoClient.PolicyUpdate(ctx, body)
    return err
}

func (h *httpClient) policyWithUser(resourceKey string,
        permType permissions.PermissionType, name string, expireTime *int64) *wrapper.Policy {
    return &wrapper.Policy{
        IsEnable:   lo.ToPtr(true),
        ExpireTime: expireTime,
        Identity: &wrapper.PolicyEntity{
            Ns:     lo.ToPtr("user"),
            PathV2: []*wrapper.Pair{{Key: "id", Val: name}},
        },
        IdentityType: lo.ToPtr[int64](0),
        Role:         h.permissionTypeToResource(permType, resourceKey),
        RoleType:     lo.ToPtr[int64](1),
        Resource: &wrapper.PolicyEntity{
            Ns:         lo.ToPtr(h.protegoClient.Namespace()),
            LocationV2: h.protegoClient.LocationV2(),
            PathV2:     []*wrapper.Pair{{Key: "key", Val: resourceKey}},
        },
        ResourceType: lo.ToPtr[int64](2),
    }
}

func (h *httpClient) permissionTypeToResource(permType permissions.PermissionType,
        resourceKey string) *wrapper.PolicyEntity {
    if permType == permissions.PermissionTypeOwner {
        return &wrapper.PolicyEntity{
            Ns:         lo.ToPtr("reftype_" + h.protegoClient.Namespace()),
            LocationV2: h.protegoClient.LocationV2(),
            PathV2: []*wrapper.Pair{
                {Key: "action", Val: string(permissions.PermissionTypeOwner)},
            },
        }
    }
    return &wrapper.PolicyEntity{
        LocationV2: h.protegoClient.LocationV2(),
        Ns:         lo.ToPtr(h.protegoClient.Namespace()),
        PathV2: []*wrapper.Pair{
            {Key: "resource", Val: resourceKey},
            {Key: "action", Val: string(permType)},
        },
    }
}

func (h *httpClient) InitProjectPermissionsExplicit(ctx context.Context, projectID int,
        creator string, owners []string) error {
    return h.initProjectPermissionsExplicit(ctx, projectID, creator, owners)
}

func (h *httpClient) InitProjectPermissions(ctx context.Context, projectID int) error {
    parcelOwnership, err := h.ownershipDecider.ParcelOwnership(ctx)
    if err != nil {
        return fmt.Errorf("get initial parcel ownership: %w", err)
    }
    return h.initProjectPermissionsExplicit(ctx, projectID, parcelOwnership.Creator, parcelOwnership.Owners)
}

func (h *httpClient) initProjectPermissionsExplicit(ctx context.Context, projectID int,
        creator string, owners []string) error {
    // creatorID, err := h.userID(ctx, creator)
    // if err != nil {
    //	return fmt.Errorf("acquiring creator id from %s: %w", creator, err)
    // }
    creatorID := int64(0)
    resourceKey := kani.ParcelResourceKey(projectID)
    body := &protego.TransactionRequest{
        Queries: []*wrapper.Query{
            // create resource itself
            {
                Operation:            common.ApiOperation_Create,
                EntityType:           lo.ToPtr(common.EntityType_Resource_Object),
                UpdateOnExist:        lo.ToPtr(true),
                ReturnCreatedResults: lo.ToPtr(true),
                EntitiesWithRefs: []*wrapper.PolicyEntityWithReference{{
                    Entity: &wrapper.PolicyEntity{
                        IsEnable:   lo.ToPtr(true),
                        Ns:         lo.ToPtr(h.protegoClient.Namespace()),
                        LocationV2: h.protegoClient.LocationV2(),
                        PathV2:     []*wrapper.Pair{{Key: "key", Val: resourceKey}},
                        // NOTE (xinyu): magic struct from protego compatible with kani, do not change
                        Attributes: map[string]interface{}{
                            "admin_url":         "https://ee.byted.org/kani/",
                            "audit_periodicity": 360,
                            "audit_status":      0,
                            "auth_block_type":   0,
                            "authorizable":      1,
                            "creator":           creatorID,
                            "description":       resourceKey,
                            "en_name":           resourceKey,
                            "name":              resourceKey,
                            "region":            "",
                            "region_type":       2,
                            "resource_url":      "https://ee.byted.org/kani",
                            "security_level":    2,
                            // kani的业务线, a.k.a. 资源所属业务/部门
                            "product": map[string]any{
                                // kani_product_112549_6gv5ym1p meaning "AI-LAB 智能语音", hard-written for now
                                "default": "kani_product_112549_6gv5ym1p",
                            },
                            // NOTE (xinyu): meaning of these optional fields is in:
                            // https://bytedance.feishu.cn/wiki/wikcn5zA6CMThQR5usCZsxUgXYd
                            // "tags": ["tag123"]
                            // "custom_tags": ["tag123"]
                            // "general_tags": ["tag123"]
                            // "line_of_business": ["tag_123"]
                        },
                    },
                    // NOTE (xinyu): no need to refer to admin group; only need to build ref to
                    // groups entities inside group_kani_xxx using the format below.
                    // References: []References{{
                    // 	UpdateOnExist:        true,
                    // 	ReturnCreatedResults: true,
                    // 	Operation:            0,
                    // 	RefType:              -1,
                    // 	RefEntities: []RefEntities{{
                    // 		Entity:     *h.permissionTypeToResource(permissions.PermissionTypeOwner, resourceKey),
                    // 		EntityType: 0,
                    // 		RefInfo: RefInfo{
                    // 			ExpireTime: 0,
                    // 			IsEnable:   true,
                    // 			Condition:  []Condition{map[string]any{}},
                    // 		},
                    // 	}},
                    // }},
                }},
            },
            // create action for results
            {
                Operation:            common.ApiOperation_Create,
                EntityType:           lo.ToPtr(common.EntityType_Role_Action),
                UpdateOnExist:        lo.ToPtr(true),
                ReturnCreatedResults: lo.ToPtr(true),
                Entities: []*wrapper.PolicyEntity{
                    h.permissionTypeToResourceOnCreate(permissions.PermissionTypeRead,
                        resourceKey, "Operation", creatorID),
                    h.permissionTypeToResourceOnCreate(permissions.PermissionTypeWrite,
                        resourceKey, "Read-only", creatorID),
                },
            },
            // create owner of resource
            {
                Relation:             lo.ToPtr(common.HierarchyRelation_Equivalent),
                UpdateOnExist:        lo.ToPtr(true),
                ReturnCreatedResults: lo.ToPtr(true),
                Policies: lo.Map(owners, func(owner string, _ int) *wrapper.Policy {
                    return h.policyWithUser(resourceKey, permissions.PermissionTypeOwner,
                        owner, lo.ToPtr(int64(0)) /*onwer never expires*/)
                }),
            },
        },
    }
    _, err := h.protegoClient.Transaction(ctx, body)
    if err != nil {
        return err
    }
    return nil
}

func (h *httpClient) permissionTypeToResourceOnCreate(permType permissions.PermissionType,
        resourceKey, resourceName string, creatorID int64) *wrapper.PolicyEntity {
    re := h.permissionTypeToResource(permType, resourceKey)
    re.IsEnable = lo.ToPtr(true)
    re.Attributes = map[string]any{
        "creator": creatorID,
        "name":    resourceName,
    }
    return re
}

func (h *httpClient) userID(ctx context.Context, name string) (int64, error) {
    entities, err := h.userEntity(ctx, name)
    if err != nil {
        return 0, err
    }
    if len(entities.Data.Entities) == 0 {
        return 0, fmt.Errorf("entity data has 0 length: %+v", entities)
    }
    id := entities.Data.Entities[0].ID

    if id == nil || *id == 0 {
        return 0, fmt.Errorf("entity data has id 0 or nil, possibly default value: %+v", entities)
    }

    return *id, nil
}

func (h *httpClient) FindUser(ctx context.Context, userName string) ([]auth.User, error) {
    entities, err := h.userEntity(ctx, userName)
    if err != nil {
        return nil, err
    }
    var users []auth.User
    for _, en := range entities.Data.Entities {
        var email, name string
        email, _ = en.Attributes["email"].(string)
        // userName, _ = en.Attributes["name"].(string)
        for _, p := range en.PathV2 {
            if p.Key == "id" {
                name = p.Val
                break
            }
        }
        users = append(users, auth.User{
            Email:    email,
            Username: name,
            Name:     name,
        })
    }
    return users, nil
}

func (h *httpClient) userEntity(ctx context.Context, name string) (*protego.PolicyEntityResponse, error) {
    body := &protego.PolicyEntityRequest{
        EntityType: 0,
        Entities: []*wrapper.PolicyEntity{
            {Ns: lo.ToPtr("user"), PathV2: []*wrapper.Pair{{Key: "id", Val: name}}},
        },
    }
    return h.protegoClient.PolicyEntityRead(ctx, body)
}

func validateUsername(user *auth.User) error {
    if user == nil || len(user.Username) == 0 {
        return fmt.Errorf("user name empty: %+v: %w", user, ErrUsername)
    }
    return nil
}
