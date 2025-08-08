package permission

import (
    "fmt"
    "sort"
)

const (
    admin     level = 4
    readWrite level = 3
    write     level = 2
    read      level = 1
)

// level is only used for sorting permissions.
// We need sorted permissions for predictable E2E testing.
type level int

const (
    Read  Type = "read"
    Write Type = "write"
    Owner Type = "admin"
)

type Type string

func (pt Type) level() level {
    switch pt {
    case Owner:
        return admin
    case Write:
        return write
    case Read:
        return read

    default:
        return level(0)
    }
}

func (pt Type) Less(permission Type) bool { return pt.level() < permission.level() }

func (pt Type) Description(projectId int) Description {
    switch pt {
    case Owner:
        name := fmt.Sprintf("project-%d_admin", projectId)
        return Description{Action: pt, Name: name}
    case Read:
        return Description{Action: pt, Name: "ReadOnly"}
    case Write:
        return Description{Action: pt, Name: "Operation"}

    default:
        return Description{Action: pt, Name: "Unknown name"}
    }
}

type Description struct {
    Name   string `json:"name"`
    Action Type   `json:"action"`
}

type Descriptions []Description

func (pd Descriptions) Len() int { return len(pd) }

func (pd Descriptions) Less(i, j int) bool {
    if pd[i].Action.level() != pd[j].Action.level() {
        return pd[i].Action.level() > pd[j].Action.level()
    }
    return pd[i].Name < pd[j].Name
}

func (pd Descriptions) Swap(i, j int) {
    pd[i], pd[j] = pd[j], pd[i]
}

func (pd Descriptions) Sort() {
    sort.Slice(pd, func(i, j int) bool { return !pd.Less(i, j) })
}

type Permission struct {
    Description
    Temporality string `json:"temporality"`
    Expires     int64  `json:"expires,omitempty"`
}

type Permissions []Permission

func (ps Permissions) Len() int { return len(ps) }

func (ps Permissions) Less(i, j int) bool {
    if ps[i].Action.level() != ps[j].Action.level() {
        return ps[i].Action.level() > ps[j].Action.level()
    }
    return ps[i].Expires < ps[j].Expires
}

func (ps Permissions) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`

    // UserAdditional is empty unless manually injected by lark.UserInfoClient
    *UserAdditional `json:",inline"`
}

type Holder struct {
    User
    Target      string      `json:"permission_target"`
    Permissions Permissions `json:"permissions"`
}

func (ph Holder) level() level {
    rtLevel := level(0)

    for _, permission := range ph.Permissions {
        switch permission.Action {
        case "admin":
            return admin
        case "write":
            switch rtLevel {
            case read, readWrite:
                rtLevel = readWrite

            default:
                rtLevel = write
            }
        case "read":
            switch rtLevel {
            case write, readWrite:
                rtLevel = readWrite

            default:
                rtLevel = readWrite
            }

        default:
            return level(0)
        }
    }

    return rtLevel
}

type Holders []Holder

func (phs Holders) Len() int { return len(phs) }

func (phs Holders) Swap(i, j int) { phs[i], phs[j] = phs[j], phs[i] }

func (phs Holders) Less(i, j int) bool {
    if phs[i].level() != phs[j].level() {
        return phs[i].level() > phs[j].level()
    }
    return phs[i].Name < phs[j].Name
}

type UserAccessControl struct {
    ProjectID   int
    Permissions []Type
}

type UserAccess []UserAccessControl

func (uc UserAccess) ProjectIDs() []int {
    var projectIDs []int

    for _, allowedProject := range uc {
        projectIDs = append(projectIDs, allowedProject.ProjectID)
    }

    return projectIDs
}

func (uc UserAccess) Permissions(projectID int) []Type {
    for _, allowedProject := range uc {
        if allowedProject.ProjectID == projectID {
            return allowedProject.Permissions
        }
    }

    return nil
}

type Project struct {
    Types  Descriptions `json:"permission_types"`
    Owners Holders      `json:"user_permissions"`
}

func (pp Project) GetOwers(permissionType Type) []User {
    var users []User
    for _, owner := range pp.Owners {
        for _, permission := range owner.Permissions {
            if permission.Action == permissionType {
                users = append(users, User{
                    Email: owner.User.Email,
                    Name:  owner.User.Name,
                })

                break
            }
        }
    }
    return users
}

func (pp Project) Sort() {
    pp.Types.Sort()
    sort.Sort(pp.Owners)

    for i := range pp.Owners {
        sort.Sort(pp.Owners[i].Permissions)
    }
}
