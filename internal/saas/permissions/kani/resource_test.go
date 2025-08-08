package kani

import (
    "bytes"
    "encoding/json"
    "errors"
    "testing"
    "time"

    // "code.byted.org/lab-speech/saas_backend/commons/more_testing"
)

func TestMarshalResourcesToPermissions(t *testing.T) {
    realWorldData := `{"project-1291": {"actions": {"admin": [{"auth_status": "enable", "create_time": "2022-07-25 17:15:15", "email": "jakub.daliga@bytedance.com", "expire_time": "", "name": "Jakub Daliga", "object_type": "user", "policy": {}, "temporality": "permanent", "user_status": "enable"}]}, "app_id": 842, "key": "project-1291", "name": "project-1291", "permissions": {"admin": {"action": "admin", "name": "project-1291_admin", "status": "enable"}, "read": {"action": "read", "name": "Read-only", "status": "enable"}, "write": {"action": "write", "name": "Operation", "status": "enable"}}, "status": "enable"}}`
    resourcesToPermissions := &ResourcesToPermissions{}
    err := json.NewDecoder(bytes.NewBufferString(realWorldData)).Decode(&resourcesToPermissions)
    more_testing.FailOnError(t, "unexpected error", err)
    // fmt.Printf("permissions: %+v\n", (*resourcesToPermissions)["project-1291"].Permissions)
    more_testing.Equals("non-nil expired time", nil, (*resourcesToPermissions)["project-1291"].Permissions["admin"][0].ExpireTime)

    realWorldData2 := `{"project-1291": {"actions": {"admin": [{"auth_status": "disabled", "create_time": "2022-07-25 17:15:15", "email": "jakub.daliga@bytedance.com", "expire_time": "2022-07-25 17:15:15", "name": "Jakub Daliga", "object_type": "user", "policy": {}, "temporality": "permanent", "user_status": "enable"}]}, "app_id": 842, "key": "project-1291", "name": "project-1291", "permissions": {"admin": {"action": "admin", "name": "project-1291_admin", "status": "enable"}, "read": {"action": "read", "name": "Read-only", "status": "enable"}, "write": {"action": "write", "name": "Operation", "status": "enable"}}, "status": "enable"}}`
    err = json.NewDecoder(bytes.NewBufferString(realWorldData2)).Decode(&resourcesToPermissions)
    more_testing.FailOnError(t, "unexpected error", err)
    expectedTime, err := time.Parse("2006-01-02 15:04:05", "2022-07-25 17:15:15")
    more_testing.FailOnError(t, "unexpected error", err)
    more_testing.FailIfNotEqual(t, "wrong parse function", 15, expectedTime.Minute())

    // fmt.Printf("exp time: %+v\n", expectedTime)
    // fmt.Printf("marshaled as: %+v\n", resourcesToPermissions)
    if !expectedTime.Equal(*(*resourcesToPermissions)["project-1291"].Permissions["admin"][0].ExpireTime) {
        more_testing.FailOnError(t, "unequal date", errors.New("should be 2022-07-25 17:15:15"))
    }
}
