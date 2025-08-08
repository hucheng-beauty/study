package kani

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "testing"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/commons/more_testing"
)

type mockHTTPDoer struct {
    t       *testing.T
    method  string
    url     string
    headers map[string]string
    body    map[string]any

    respBody string
}

func (doer mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
    more_testing.ErrorIfNotEqual(doer.t, "invalid method", doer.method, req.Method)
    more_testing.ErrorIfNotEqual(doer.t, "invalid url", doer.url, req.URL.String())
    for k, v := range doer.headers {
        more_testing.ErrorIfNotEqual(doer.t, fmt.Sprintf("invalid header %s value", k), v, req.Header.Get(k))
    }
    if doer.body == nil {
        return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(doer.respBody))}, nil
    }
    b, err := ioutil.ReadAll(req.Body)
    more_testing.FailOnError(doer.t, "read request body", err)
    var m map[string]any
    err = json.Unmarshal(b, &m)
    more_testing.FailOnError(doer.t, "unmarshal request body", err)
    more_testing.ErrorIfNotEqual(doer.t, "invalid request body", doer.body, m)
    return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(doer.respBody))}, nil
}

func TestHTTPClient_CreateResource(t *testing.T) {
    resource := Resource{
        Name:               "Test",
        EnglishName:        "Test",
        Key:                "Test",
        Creator:            "Tester",
        TemporaryAuthorize: false,
        OwnerKeys:          []string{"Tester"},
        ProductName:        "AI-Lab智能语音",
        Confidentiality:    "L2",
    }
    httpDoer := mockHTTPDoer{
        t:       t,
        method:  http.MethodPost,
        url:     "https://ei.byted.org/ratak/resources/",
        headers: map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:123"))},
        body: map[string]interface{}{
            "name":                "Test",
            "en_name":             "Test",
            "key":                 "Test",
            "creator_key":         "Tester",
            "temporary_authorize": false,
            "owner_keys":          []interface{}{"Tester"},
            "product_name":        "AI-Lab智能语音",
            "confidentiality":     "L2",
        },
        respBody: "",
    }

    ctx := context.Background()
    client := HTTPClient{
        doer:       httpDoer,
        appID:      "user",
        secret:     "123",
        baseApiURL: "https://ei.byted.org/ratak",
    }
    err := client.CreateResource(ctx, resource)
    more_testing.FailOnError(t, "unexpected error", err)
}

func TestHTTPClient_AddPermissions(t *testing.T) {
    modifyRequest := ModifyPermissionRequest{
        ResourceKey:   "project-123",
        Actions:       []string{"admin", "read"},
        UserKeys:      []string{"user.123"},
        AuthorizeDays: 12,
    }
    httpDoer := mockHTTPDoer{
        t:      t,
        method: http.MethodPost,
        url:    "https://ei.byted.org/ratak/permissions/relationships/",
        headers: map[string]string{
            "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:123")),
            "x-app-version": "v2",
        },
        body: map[string]any{
            "resource_key":   "project-123",
            "actions":        []any{"admin", "read"},
            "user_keys":      []any{"user.123"},
            "authorize_days": 12.0,
        },
        respBody: "",
    }

    ctx := context.Background()
    client := HTTPClient{
        doer:       httpDoer,
        appID:      "user",
        secret:     "123",
        baseApiURL: "https://ei.byted.org/ratak",
    }
    err := client.AddPermissions(ctx, modifyRequest)
    more_testing.FailOnError(t, "unexpected error", err)
}

func TestHTTPClient_UserPermissions(t *testing.T) {
    httpDoer := mockHTTPDoer{
        t:      t,
        method: http.MethodGet,
        url:    "https://ee.byted.org/ratak/employees/tester/permissions/",
        headers: map[string]string{
            "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:123")),
        },
        body:     nil,
        respBody: "{\"AAA-123\": [\"admin\"]}",
    }
    expectedUserPermissions := []PermissionEntry{
        {
            ResourceKey:         "AAA-123",
            ResourcePermissions: []string{"admin"},
        },
    }

    ctx := context.Background()
    client := HTTPClient{
        doer:                             httpDoer,
        appID:                            "user",
        secret:                           "123",
        baseApiURL:                       "https://ee.byted.org/ratak/",
        getAllowedPermissionsURLTemplate: "https://ee.byted.org/ratak/employees/%s/permissions/",
    }
    perms, err := client.UserPermissions(ctx, auth.User{Username: "tester"})
    more_testing.FailOnError(t, "unexpected error", err)
    more_testing.ErrorIfNotEqual(t, "user resource permission mismatch", expectedUserPermissions, perms)
}

func TestHTTPClient_ResourcePermissions(t *testing.T) {
    httpDoer := mockHTTPDoer{
        t:      t,
        method: http.MethodGet,
        url:    "https://ei.byted.org/ratak/audit/resources/project-1291/?exclude_disabled_permissions=true",
        headers: map[string]string{
            "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:123")),
            "x-app-version": "v2",
        },
        body: nil,
        respBody: `
		{
			"project-1291":{
			  "actions":{
				"admin":[
				  {
					"auth_status":"enable",
					"create_time":"2022-07-25 17:15:15",
					"email":"jakub.daliga@bytedance.com",
					"expire_time":"",
					"name":"Jakub Daliga",
					"object_type":"user",
					"policy":{
					  
					},
					"temporality":"permanent",
					"user_status":"enable"
				  }
				]
			  },
			  "app_id":842,
			  "key":"project-1291",
			  "name":"project-1291",
			  "permissions":{
				"admin":{
				  "action":"admin",
				  "name":"project-1291_admin",
				  "status":"enable"
				}
			  },
			  "status":"enable"
			}
		  }
		`,
    }
    expectedPermissions := ResourcePermissions{
        Status:       "enable",
        AppID:        842,
        ResourceKey:  "project-1291",
        ResourceName: "project-1291",
        PermissionTypes: map[string]PermissionType{
            "admin": {Action: "admin", Name: "project-1291_admin"},
        },
        Permissions: map[string][]Permission{"admin": {
            {
                Email:       "jakub.daliga@bytedance.com",
                Name:        "Jakub Daliga",
                ObjectType:  "user",
                Temporality: "permanent",
            },
        }},
    }

    ctx := context.Background()
    client := HTTPClient{
        doer:       httpDoer,
        appID:      "user",
        secret:     "123",
        baseApiURL: "https://ei.byted.org/ratak",
    }
    perms, err := client.ResourcePermissions(ctx, "project-1291")
    more_testing.FailOnError(t, "unexpected error", err)
    more_testing.ErrorIfNotEqual(t, "resource permissions mismatch", expectedPermissions, perms)
}
