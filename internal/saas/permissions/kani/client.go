package kani

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"

    "study/internal/saas/permissions"

    "github.com/pkg/errors"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons"
    // "code.byted.org/lab-speech/saas_backend/commons/errutil"
    // "code.byted.org/lab-speech/saas_backend/permissions"
)

// Docs: https://bytedance.feishu.cn/space/doc/c7hcpsDmtldXT1XyYOSFla#
// API: https://ei.byted.org/ratak/doc/
type HTTPClient struct {
    doer commons.HTTPDoer

    appID                            string
    secret                           string
    baseApiURL                       string
    getAllowedPermissionsURLTemplate string
}

type PermissionsChange struct {
    Username      string   `json:"username"`
    Actions       []string `json:"actions"`
    AuthorizeDays int      `json:"authorize_days"`
}

func NewClient(doer commons.HTTPDoer, appID, secret, baseApiURL, getAllowedPermissionsURLTemplate string) Client {
    return HTTPClient{
        doer:                             doer,
        appID:                            appID,
        secret:                           secret,
        baseApiURL:                       baseApiURL,
        getAllowedPermissionsURLTemplate: getAllowedPermissionsURLTemplate,
    }
}

type PermissionEntry struct {
    ResourceKey         string
    ResourcePermissions []string
}

// UserAccessControl converts []PermissionEntry to []permissions.UserAccessControl.
func UserAccessControl(permissionEntries []PermissionEntry) ([]permissions.UserAccessControl, error) {
    var userAccess permissions.UserAccess
    for _, grantedPermissions := range permissionEntries {
        var grantedPermissionTypes []permissions.PermissionType
        for _, permission := range grantedPermissions.ResourcePermissions {
            grantedPermissionTypes = append(grantedPermissionTypes, permissions.PermissionType(permission))
        }
        projectID, err := ParcelIDFromResourceKey(grantedPermissions.ResourceKey)
        if err != nil {
            return nil, fmt.Errorf("project id from [resource key = %s]: %w",
                grantedPermissions.ResourceKey, err)
        }
        userAccess = append(userAccess, permissions.UserAccessControl{
            ProjectID:       projectID,
            UserPermissions: grantedPermissionTypes,
        })
    }
    return userAccess, nil
}

func (c HTTPClient) CreateResource(ctx context.Context, resource Resource) error {
    resourceBuf, err := commons.JSONReader(resource)
    if err != nil {
        return fmt.Errorf("json reader: %w", err)
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost,
        fmt.Sprintf("%s/resources/", c.baseApiURL), resourceBuf)
    if err != nil {
        return fmt.Errorf("create request with context: %w", err)
    }
    c.authorizeRequest(req)
    doer := commons.WithTTLogIDHTTPDoer(ctx, c.doer)
    resp, err := doer.Do(req)
    defer commons.CleanupResponse(resp)
    if err != nil {
        return errors.Wrap(err, "do request")
    }
    if resp.StatusCode/100 != 2 {
        return errutil.InvalidStatusCodeErr("create kani resource", resp)
    }
    return nil
}

func (c HTTPClient) ResourcePermissions(ctx context.Context, resourceKey string) (ResourcePermissions, error) {
    // NOTE(xinyu): /audit/resource/xxx/ is deprecated and now results in error in i18n;
    // CN has a whitelist mechanism; but for compatibility, I am changing this to the newest
    // api:
    // [Audit - 审计资源的权限授予明细](https://ei.byted.org/ratak/doc/#api-Audit-AuditResource-Get)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet,
        fmt.Sprintf("%s/audit/resources/%s/", c.baseApiURL, resourceKey), nil)
    if err != nil {
        return ResourcePermissions{}, fmt.Errorf("create request with context: %w", err)
    }
    req.Header.Set("x-app-version", "v2")

    reqParams := url.Values{}
    reqParams.Add("exclude_disabled_permissions", "true")
    req.URL.RawQuery = reqParams.Encode()

    c.authorizeRequest(req)

    doer := commons.WithTTLogIDHTTPDoer(ctx, c.doer)
    resp, err := doer.Do(req)
    defer commons.CleanupResponse(resp)
    if err != nil {
        return ResourcePermissions{}, fmt.Errorf("do request: %w", err)
    }
    if resp.StatusCode/100 != 2 {
        return ResourcePermissions{}, errutil.InvalidStatusCodeErr(
            "get resource permissions", resp,
        )
    }
    resourcesToPermissions := &ResourcesToPermissions{}
    err = json.NewDecoder(resp.Body).Decode(&resourcesToPermissions)
    if err != nil {
        return ResourcePermissions{}, fmt.Errorf("decode response: %w", err)
    }
    resourcePermission, rok := (*resourcesToPermissions)[resourceKey]
    if !rok {
        return ResourcePermissions{},
            fmt.Errorf(
                "required resource not in response list: resourceKey: %s: resp: %+v",
                resourceKey, resourcesToPermissions,
            )
    }
    return resourcePermission, nil
}

func (c HTTPClient) UserPermissions(ctx context.Context, user auth.User) ([]PermissionEntry, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet,
        fmt.Sprintf(c.getAllowedPermissionsURLTemplate, user.Username), nil)
    if err != nil {
        return nil, fmt.Errorf("create request with context: %w", err)
    }
    c.authorizeRequest(req)

    doer := commons.WithTTLogIDHTTPDoer(ctx, c.doer)
    resp, err := doer.Do(req)
    defer commons.CleanupResponse(resp)
    if err != nil {
        return nil, errors.Wrap(err, "do request")
    }
    if resp.StatusCode != 200 {
        return nil, errutil.InvalidStatusCodeErr("get allowed resources", resp)
    }
    rawPermissions := map[string][]string{}
    err = json.NewDecoder(resp.Body).Decode(&rawPermissions)
    if err != nil {
        return nil, errors.Wrap(err, "decode response")
    }
    var permissions []PermissionEntry
    for resourceKey, grantedPermissions := range rawPermissions {
        resourceKey, grantedPermissions := resourceKey, grantedPermissions
        permissions = append(permissions, PermissionEntry{
            ResourceKey:         resourceKey,
            ResourcePermissions: grantedPermissions,
        })
    }
    return permissions, nil
}

func (c HTTPClient) FindUser(ctx context.Context, query string) ([]UserQueryResult, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost,
        fmt.Sprintf("%s/users/search/?query=%s", c.baseApiURL, url.QueryEscape(query)), nil)
    if err != nil {
        return nil, fmt.Errorf("create request with context: %w", err)
    }
    req.Header.Set("x-app-version", "v2")
    c.authorizeRequest(req)
    doer := commons.WithTTLogIDHTTPDoer(ctx, c.doer)
    resp, err := doer.Do(req)
    defer commons.CleanupResponse(resp)
    if err != nil {
        return nil, errors.Wrap(err, "do request")
    }
    if resp.StatusCode/100 != 2 {
        return nil, errutil.InvalidStatusCodeErr("user search", resp)
    }
    var queryResult []UserQueryResult
    err = json.NewDecoder(resp.Body).Decode(&queryResult)
    if err != nil {
        return nil, errors.Wrap(err, "decode response body")
    }
    return queryResult, nil
}

func (c HTTPClient) AddPermissions(ctx context.Context, reqBody ModifyPermissionRequest) error {
    buf, err := c.encodeRequest(reqBody)
    if err != nil {
        return fmt.Errorf("encode request: %w", err)
    }
    err = c.modifyPermissions(ctx, buf, http.MethodPost)
    if err != nil {
        return fmt.Errorf("modify permissions: %w", err)
    }
    return nil
}

func (c HTTPClient) DeletePermissions(ctx context.Context, reqBody ModifyPermissionRequest) error {
    buf, err := c.encodeRequest(reqBody)
    if err != nil {
        return fmt.Errorf("encode request: %w", err)
    }
    err = c.modifyPermissions(ctx, buf, http.MethodDelete)
    if err != nil {
        return fmt.Errorf("modify permissions: %w", err)
    }
    return nil
}

func (c HTTPClient) encodeRequest(reqBody ModifyPermissionRequest) (*bytes.Buffer, error) {
    buf := bytes.NewBuffer(nil)
    err := json.NewEncoder(buf).Encode(reqBody)
    if err != nil {
        return nil, fmt.Errorf("encode request body: %w", err)
    }
    return buf, nil
}

func (c HTTPClient) modifyPermissions(ctx context.Context, buf *bytes.Buffer, requestType string) error {
    req, err := http.NewRequestWithContext(ctx, requestType,
        fmt.Sprintf("%s/permissions/relationships/", c.baseApiURL), buf)
    if err != nil {
        return fmt.Errorf("create request with context: %w", err)
    }
    req.Header.Set("x-app-version", "v2")
    c.authorizeRequest(req)
    doer := commons.WithTTLogIDHTTPDoer(ctx, c.doer)
    resp, err := doer.Do(req)
    defer commons.CleanupResponse(resp)
    if err != nil {
        return fmt.Errorf("do request: %w", err)
    }
    if resp.StatusCode/100 != 2 {
        return errutil.InvalidStatusCodeErr("modify resource permissions", resp)
    }
    return nil
}

func (c HTTPClient) authorizeRequest(r *http.Request) {
    r.SetBasicAuth(c.appID, c.secret)
}
