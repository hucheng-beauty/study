package permissions

import (
    "fmt"
    "sort"

    // "code.byted.org/lab-speech/saas_backend/auth"
)

const (
    PermissionTypeRead  PermissionType = "read"
    PermissionTypeWrite PermissionType = "write"
    PermissionTypeOwner PermissionType = "admin"

    permissionLevelAdmin permissionLevel = 4
    permissionLevelRW    permissionLevel = 3
    permissionLevelW     permissionLevel = 2
    permissionLevelR     permissionLevel = 1
)

type (
    PermissionType    string
    Descriptors       []Descriptor
    Permissions       []Permission
    PermissionHolders []PermissionHolder
    UserAccess        []UserAccessControl
)

func (u UserAccess) UserAllowedProjectIDs() []int {
    var allowedProjectIDs []int
    for _, allowedProject := range u {
        allowedProjectIDs = append(allowedProjectIDs, allowedProject.ProjectID)
    }
    return allowedProjectIDs
}

func (u UserAccess) PermissionsToProject(projectID int) []PermissionType {
    for _, allowedProject := range u {
        if allowedProject.ProjectID == projectID {
            return allowedProject.UserPermissions
        }
    }
    return nil
}

type UserAccessControl struct {
    ProjectID       int
    UserPermissions []PermissionType
}

func (pp ProjectPermissions) GetUsersWithPermissionType(permissionType PermissionType) []auth.User {
    var users []auth.User
    for _, permissionOwner := range pp.PermissionOwners {
        for _, permission := range permissionOwner.Permissions {
            if permission.Type == permissionType {
                users = append(users, auth.User{Email: permissionOwner.User.Email, Name: permissionOwner.User.Name})
                break
            }
        }
    }
    return users
}

func (pp ProjectPermissions) Sort() {
    pp.PermissionTypes.Sort()
    sort.Sort(pp.PermissionOwners)
    for i := range pp.PermissionOwners {
        sort.Sort(pp.PermissionOwners[i].Permissions)
    }
}

type ProjectPermissions struct {
    PermissionTypes  Descriptors       `json:"permission_types"`
    PermissionOwners PermissionHolders `json:"user_permissions"`
}

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    // UserAdditional is empty unless manually injected by lark.UserInfoClient
    *UserAdditional `json:",inline"`
}

type PermissionHolder struct {
    User
    Target      string      `json:"permission_target"`
    Permissions Permissions `json:"permissions"`
}

func (s PermissionHolders) Len() int { return len(s) }

func (s PermissionHolders) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s PermissionHolders) Less(i, j int) bool {
    if s[i].permissionLevel() != s[j].permissionLevel() {
        return s[i].permissionLevel() > s[j].permissionLevel()
    }
    return s[i].Name < s[j].Name
}

type Descriptor struct {
    Type PermissionType `json:"action"`
    Name string         `json:"name"`
}

func (p Descriptors) Len() int { return len(p) }

func (p Descriptors) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p Descriptors) Less(i, j int) bool {
    if p[i].Type.permissionLevel() != p[j].Type.permissionLevel() {
        return p[i].Type.permissionLevel() > p[j].Type.permissionLevel()
    }
    return p[i].Name < p[j].Name
}

func (p Descriptors) Sort() {
    sort.Slice(p, func(i, j int) bool { return !p.Less(i, j) })
}

// Less returns `true` if permission `p` is of lower level than permission `p2`.
// TODO: Reduce the amount of types associated with permissions.
func (p PermissionType) Less(p2 PermissionType) bool {
    return p.permissionLevel() < p2.permissionLevel()
}

func (p PermissionType) permissionLevel() permissionLevel {
    switch p {
    case PermissionTypeOwner:
        return permissionLevelAdmin
    case PermissionTypeWrite:
        return permissionLevelW
    case PermissionTypeRead:
        return permissionLevelR
    default:
        return 0
    }
}

func (p PermissionType) Descriptor(parcelID int) Descriptor {
    switch p {
    case PermissionTypeOwner:
        return Descriptor{
            Type: p,
            Name: fmt.Sprintf("project-%d_admin", parcelID),
        }
    case PermissionTypeRead:
        return Descriptor{
            Type: p,
            Name: "Read-only",
        }
    case PermissionTypeWrite:
        return Descriptor{
            Type: p,
            Name: "Operation",
        }
    default:
        return Descriptor{
            Type: p,
            Name: "Unknown name",
        }
    }
}

type Permission struct {
    Descriptor
    Temporality string `json:"temporality"`
    Expires     int64  `json:"expires,omitempty"`
}

func (p Permissions) Len() int { return len(p) }

func (p Permissions) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p Permissions) Less(i, j int) bool {
    if p[i].Type.permissionLevel() != p[j].Type.permissionLevel() {
        return p[i].Type.permissionLevel() > p[j].Type.permissionLevel()
    }
    return p[i].Expires < p[j].Expires
}

// permissionLevel is only used for sorting permissions.
// We need sorted permissions for predictable E2E testing.
type permissionLevel int

func (p PermissionHolder) permissionLevel() permissionLevel {
    level := permissionLevel(0)
    for _, perm := range p.Permissions {
        switch perm.Type {
        case "admin":
            return permissionLevelAdmin
        case "write":
            switch level {
            case permissionLevelR, permissionLevelRW:
                level = permissionLevelRW
            default:
                level = permissionLevelW
            }
        case "read":
            switch level {
            case permissionLevelW, permissionLevelRW:
                level = permissionLevelRW
            default:
                level = permissionLevelR
            }
        default:
            return 0
        }
    }
    return level
}
