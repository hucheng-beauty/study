package kani

import (
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "time"
)

const (
    DefaultAuthorizeDays = 30
    // MaxAuthorizeDays is the maximum authorization period supported by KANI (2 years - 730 days).
    MaxAuthorizeDays = 730
    // ParcelResourceKeyTemplate is template from which we generate KANI resource key for parcel.
    ParcelResourceKeyTemplate = "project-%d"
)

// ParcelResourceKey returns ResourceKey for parcel with provided id.
func ParcelResourceKey(parcelID int) string {
    return fmt.Sprintf(ParcelResourceKeyTemplate, parcelID)
}

// ParcelIDFromResourceKey returns parcel id from provided resourceKey.
func ParcelIDFromResourceKey(resourceKey string) (int, error) {
    splitKey := strings.Split(resourceKey, "-")
    if len(splitKey) != 2 {
        return 0, fmt.Errorf("[key=%s] invalid resourceKey to get ParcelID", resourceKey)
    }
    projectID, err := strconv.Atoi(splitKey[1])
    if err != nil {
        return 0, fmt.Errorf("[key=%s] convert resource key suffix to ParcelID: %w", resourceKey, err)
    }
    return projectID, nil
}

// Resource is post request body content as described in
// https://ei.byted.org/ratak/doc/#api-Resource-Resources-Post
type Resource struct {
    Name               string   `json:"name"`                // required
    EnglishName        string   `json:"en_name"`             // required
    Key                string   `json:"key"`                 // required
    Creator            string   `json:"creator_key"`         // required
    TemporaryAuthorize bool     `json:"temporary_authorize"` // required
    OwnerKeys          []string `json:"owner_keys"`          // required
    // ProductName: required, it is actually "Affiliated Business or Department"
    // https://ei.byted.org/ratak/doc/#api-Resource-Resources-Post
    ProductName     string `json:"product_name"`
    Confidentiality string `json:"confidentiality"` // required

    URL         *string           `json:"url,omitempty"`         // optional
    AdminURL    string            `json:"admin_url,omitempty"`   // optional
    Description *string           `json:"description,omitempty"` // optional
    Actions     map[string]string `json:"actions,omitempty"`     // optional
    ConfigKey   *string           `json:"config_key,omitempty"`  // optional
    ParentKeys  []string          `json:"parent_keys,omitempty"` // optional
    Tags        []string          `json:"tags,omitempty"`        // optional
}

type ResourcesToPermissions map[string]ResourcePermissions

// UnmarshalJSON is to unmarshal `ResourcesToPermissions`; it is customized
// because `expire_time` might be an empty string, and that will result in
// err;
func (rtp *ResourcesToPermissions) UnmarshalJSON(data []byte) error {
    tmp := map[string]tmpResourcePermissions{}
    err := json.Unmarshal(data, &tmp)
    if err != nil {
        return err
    }
    *rtp = ResourcesToPermissions{}
    for k, v := range tmp {
        permissions := map[string][]Permission{}
        for permKey, permValue := range v.Permissions {
            for _, perm := range permValue {
                var expireTime *time.Time
                if perm.ExpireTime != "" {
                    parsedTime, err := time.Parse("2006-01-02 15:04:05", perm.ExpireTime)
                    if err != nil {
                        return err
                    }
                    expireTime = &parsedTime
                }
                permissions[permKey] = append(permissions[permKey], Permission{
                    Email:       perm.Email,
                    Name:        perm.Name,
                    ObjectType:  perm.ObjectType,
                    Temporality: perm.Temporality,
                    ExpireTime:  expireTime,
                })
            }
        }
        (*rtp)[k] = ResourcePermissions{
            Status:          v.Status,
            AppID:           v.AppID,
            ResourceKey:     v.ResourceKey,
            ResourceName:    v.ResourceName,
            PermissionTypes: v.PermissionTypes,
            Permissions:     permissions,
        }
    }
    return nil
}

type ResourcePermissions struct {
    // if project is enabled - ignore
    Status string `json:"status"`
    // ID of KANI application
    AppID        int    `json:"app_id"`
    ResourceKey  string `json:"key"`
    ResourceName string `json:"name"`

    // List of permissions that exist for this project
    PermissionTypes map[string]PermissionType `json:"permissions"`
    // Map permission_type -> list of objects with given permission
    Permissions map[string][]Permission `json:"actions"`
}

type PermissionType struct {
    // Codename of permission (machine-readable)
    Action string `json:"action"`
    // Displayed name of permission
    Name string `json:"name"`
}

type Permission struct {
    // email of user with permission
    Email string `json:"email"`
    // name of user with permission
    Name string `json:"name"`
    // type of object that holds given permission - "user" for us
    ObjectType string `json:"object_type"`
    // if permission is temporary or permament (temporary/permament)
    Temporality string `json:"temporality"`
    // when permission will expire (user will lose access)
    ExpireTime *time.Time `json:"expire_time,omitempty"`
}

type UserQueryResult struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type ModifyPermissionRequest struct {
    ResourceKey   string   `json:"resource_key"`
    Actions       []string `json:"actions"`
    UserKeys      []string `json:"user_keys"`
    AuthorizeDays int      `json:"authorize_days"`
}

type tmpResourcePermissions struct {
    // if project is enabled - ignore
    Status string `json:"status"`
    // ID of KANI application
    AppID        int    `json:"app_id"`
    ResourceKey  string `json:"key"`
    ResourceName string `json:"name"`

    // List of permissions that exist for this project
    PermissionTypes map[string]PermissionType `json:"permissions"`
    // Map permission_type -> list of objects with given permission
    Permissions map[string][]tmpPermission `json:"actions"`
}

type tmpPermission struct {
    // email of user with permission
    Email string `json:"email"`
    // name of user with permission
    Name string `json:"name"`
    // type of object that holds given permission - "user" for us
    ObjectType string `json:"object_type"`
    // if permission is temporary or permament (temporary/permament)
    Temporality string `json:"temporality"`
    // when permission will expire (user will lose access)
    ExpireTime string `json:"expire_time,omitempty"`
}
