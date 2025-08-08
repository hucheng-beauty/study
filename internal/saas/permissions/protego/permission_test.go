package protego

import (
    "context"
    "net/http"
    "os"
    "reflect"
    "testing"
    "time"

    // "code.byted.org/gopkg/ctxvalues"
    // "code.byted.org/gopkg/logid"
    // "code.byted.org/gopkg/logs/v2/log"
    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/clients/protego"
    // "code.byted.org/lab-speech/saas_backend/commons"
    // "code.byted.org/lab-speech/saas_backend/gardening"
    // "code.byted.org/lab-speech/saas_backend/permissions"
    // "code.byted.org/lab-speech/saas_backend/permissions/kani"
    jsoniter "github.com/json-iterator/go"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TestProjectPermissions_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    projectID := 615
    var doer commons.HTTPDoer = http.DefaultClient
    doer = commons.WithHTTPDumpDoer(doer)
    doer = commons.WithTTEnvFromReqCtxHTTPDoer(doer)

    client := testPermissionManager()
    resp, err := client.ProjectPermissions(ctx, projectID)
    assert.NoError(t, err)
    log.V2.Warn().With(ctx).Str("protego resp:\n").Obj(resp).Emit()

    kaniClient := kani.NewClient(doer, "842", "C1094D65425A40C889C25C8CC484DAD3",
        "https://ei.byted.org/ratak", "https://ee.byted.org/ratak/employees/%v/permissions/")
    concierge := kani.NewPermissionManager(kaniClient, gardening.ContextUserOwnershipDecider{})
    perms, err := concierge.ProjectPermissions(ctx, projectID)
    assert.NoError(t, err)
    log.V2.Warn().With(ctx).Str("kani resp:\n").Obj(perms).Emit()
}

func TestProjectsPermissions_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    projectIDs := []int{5776, 5220}
    var doer commons.HTTPDoer = http.DefaultClient
    doer = commons.WithHTTPDumpDoer(doer)
    doer = commons.WithTTEnvFromReqCtxHTTPDoer(doer)

    client := testPermissionManager()
    resp, err := client.(*httpClient).ProjectsPermissions(ctx, projectIDs)
    assert.NoError(t, err)
    log.V2.Warn().With(ctx).Str("protego resp:\n").Obj(resp).Emit()
    // kaniClient := kani.NewClient(doer, "842", "C1094D65425A40C889C25C8CC484DAD3",
    // 	"https://ei.byted.org/ratak", "https://ee.byted.org/ratak/employees/%v/permissions/")
    // concierge := kani.NewPermissionManager(kaniClient, gardening.ContextUserOwnershipDecider{})
    // perms, err := concierge.ProjectPermissions(ctx, projectID)
    // assert.NoError(t, err)
    // log.V2.Warn().With(ctx).Str("kani resp:\n").Obj(perms).Emit()
}

func testPermissionManager() Client {
    cl := protego.NewHTTPClient(protego.HTTPClientConfig{Namespace: "kani_841"}, testDoer())
    // cl.OverwriteLocationV2([]*wrapper.Pair{{Key: "r", Val: "cn"}})
    return NewHTTPPermissionManager(cl, gardening.ContextUserOwnershipDecider{})
}

func testDoer() commons.HTTPDoer {
    var doer commons.HTTPDoer = http.DefaultClient
    doer = commons.WithHTTPDumpDoer(doer)
    return commons.WithTTEnvFromReqCtxHTTPDoer(doer)
}

func TestAddPermissions_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    projectID := 615

    client := testPermissionManager()
    err := client.AddPermissions(ctx, projectID, auth.User{Username: "yuanyichen"},
        []permissions.PermissionType{permissions.PermissionTypeOwner}, 365)
    assert.NoError(t, err)
}

func TestDeletePermissions_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    projectID := 1330

    client := testPermissionManager()
    err := client.DeletePermissions(ctx, projectID, auth.User{Name: "huangxinyu.eng"},
        []permissions.PermissionType{permissions.PermissionTypeOwner,
            permissions.PermissionTypeWrite})
    assert.NoError(t, err)
}

func TestFindUser_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    client := testPermissionManager()

    users, err := client.FindUser(ctx, "huangxinyu.eng")
    assert.NoError(t, err)
    log.V2.Warn().Str("users:").Obj(users).Emit()
}

func TestUserAccess_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    userName := "gaohongliang.ghl"
    ctx = auth.CtxWithUser(ctx, auth.User{Name: userName, Username: userName})
    client := testPermissionManager()

    users, err := client.UserAccess(ctx)
    assert.NoError(t, err)
    log.V2.Warn().Str("users:").Obj(users).Emit()

    kaniClient := kani.NewClient(testDoer(), "842", "C1094D65425A40C889C25C8CC484DAD3",
        "https://ei.byted.org/ratak", "https://ee.byted.org/ratak/employees/%v/permissions/")
    concierge := kani.NewPermissionManager(kaniClient, gardening.ContextUserOwnershipDecider{})
    perms, err := concierge.UserAccess(ctx)
    assert.NoError(t, err)
    log.V2.Warn().With(ctx).Str("kani users:").Obj(perms).Emit()

    assert.ElementsMatch(t, perms, users)
}

func TestInitProject_E2E(t *testing.T) {
    t.Skip("comment this to run end 2 end test")
    defer log.Flush()
    userName := "huangxinyu.eng"
    projectID := 99999

    ctx := context.Background()
    ctx = ctxvalues.SetLogID(ctx, logid.GenLogID())
    ctx = auth.CtxWithUser(ctx, auth.User{Name: userName, Username: userName})
    client := testPermissionManager()

    err := client.InitProjectPermissions(ctx, projectID)
    assert.NoError(t, err)
}

func TestProjectPermissionsConversion(t *testing.T) {
    projectID := 1130
    cl := &httpClient{}
    var resp *protego.PolicyResponse
    err := json.UnmarshalFromString(policyRespStr, &resp)
    assert.NoError(t, err)
    conv := cl.projectsPermissions(resp)
    convStr, err := json.MarshalToString(conv[projectID])
    assert.NoError(t, err)
    assert.JSONEq(t, kaniConverted, convStr, "conv = %s", convStr)
}

func TestNanoTimeConvert(t *testing.T) {
    nano := time.Now().Nanosecond()
    nanoT := time.Unix(0, int64(nano-1))
    assert.True(t, nanoT.Before(time.Now()))
}

var policyRespStr string = `{"status":0,"data":{"operation":1,"policies":[{"id":96921731,"version":1649498891988765000,"identity":{"id":121163,"version":1685358503619264000,"ns":"user","location":"","path":"id:michal.siwinski","metadata":[{}],"attributes":{"name":"Michał Siwiński","email":"michal.siwinski@bytedance.com","employee_type":1,"is_eu_white":1,"is_gso_white":0,"en_name":"Michał Siwiński","en_work_country":"Poland","ldap_role":"正式","dimission_stage":-1,"en_ldap_title":"RD","ldap_title":"技术","en_ldap_role":"Regular","is_tg":0,"large_department":"AI-Lab","work_country":"波兰"},"pathV2":[{"key":"id","val":"michal.siwinski"}],"locationV2":[],"is_enable":false,"create_time":1642676140330740000},"identity_type":0,"role":{"id":69367348,"version":1673846119877771000,"ns":"kani_842","location":"r:cn","path":"resource:project-1130/action:read","metadata":[{}],"attributes":{"creator":199280,"is_auto_approval":0,"name":"Read-only"},"pathV2":[{"key":"resource","val":"project-1130"},{"key":"action","val":"read"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1648788642446556000},"role_type":1,"resource":{"id":10841918,"version":1663814588680421000,"ns":"kani_842","location":"r:cn","path":"key:project-1130","metadata":[{}],"attributes":{"admin_url":"","region":"","region_type":0,"security_level":2,"auth_block_type":0,"en_name":"","audit_status":0,"resource_url":"","creator":199280,"description":"","name":"project-1130","status":1,"audit_periodicity":360,"authorizable":0},"pathV2":[{"key":"key","val":"project-1130"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1649325694850133000},"resource_type":2,"expire_time":1651661777000000000,"condition":[{}],"is_enable":false,"is_break_glass":false,"information":{},"create_time":1649498891988765000},{"id":96921732,"version":1649498891988765000,"identity":{"id":121163,"version":1685358503619264000,"ns":"user","location":"","path":"id:michal.siwinski","metadata":[{}],"attributes":{"employee_type":1,"is_tg":0,"dimission_stage":-1,"email":"michal.siwinski@bytedance.com","en_ldap_title":"RD","en_name":"Michał Siwiński","en_ldap_role":"Regular","en_work_country":"Poland","is_eu_white":1,"large_department":"AI-Lab","ldap_role":"正式","ldap_title":"技术","work_country":"波兰","is_gso_white":0,"name":"Michał Siwiński"},"pathV2":[{"key":"id","val":"michal.siwinski"}],"locationV2":[],"is_enable":false,"create_time":1642676140330740000},"identity_type":0,"role":{"id":69367350,"version":1673846119880125000,"ns":"kani_842","location":"r:cn","path":"resource:project-1130/action:write","metadata":[{}],"attributes":{"is_auto_approval":0,"name":"Operation","creator":199280},"pathV2":[{"key":"resource","val":"project-1130"},{"key":"action","val":"write"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1648788642449363000},"role_type":1,"resource":{"id":10841918,"version":1663814588680421000,"ns":"kani_842","location":"r:cn","path":"key:project-1130","metadata":[{}],"attributes":{"authorizable":0,"region_type":0,"security_level":2,"auth_block_type":0,"description":"","name":"project-1130","resource_url":"","admin_url":"","audit_status":0,"creator":199280,"en_name":"","audit_periodicity":360,"status":1,"region":""},"pathV2":[{"key":"key","val":"project-1130"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1649325694850133000},"resource_type":2,"expire_time":1651661777000000000,"condition":[{}],"is_enable":false,"is_break_glass":false,"information":{},"create_time":1649498891988765000},{"id":97935171,"version":1649522420058047000,"identity":{"id":199280,"version":1685382901263269000,"ns":"user","location":"","path":"id:jakub.daliga","metadata":[{}],"attributes":{"email":"jakub.daliga@bytedance.com","employee_type":1,"en_ldap_title":"RD","en_name":"Jakub Daliga","is_eu_white":1,"work_country":"波兰","en_work_country":"Poland","is_gso_white":0,"large_department":"AI-Lab","ldap_role":"正式","name":"Jakub Daliga","dimission_stage":-1,"en_ldap_role":"Regular","is_tg":0,"ldap_title":"技术"},"pathV2":[{"key":"id","val":"jakub.daliga"}],"locationV2":[],"is_enable":false,"create_time":1642700373646106000},"identity_type":0,"role":{"id":1798815,"version":1666594958997705000,"ns":"reftype_kani_842","location":"r:cn","path":"action:admin","metadata":[{}],"pathV2":[{"key":"action","val":"admin"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1643530116475578000},"role_type":1,"resource":{"id":10841918,"version":1663814588680421000,"ns":"kani_842","location":"r:cn","path":"key:project-1130","metadata":[{}],"attributes":{"audit_status":0,"auth_block_type":0,"en_name":"","security_level":2,"creator":199280,"description":"","region":"","region_type":0,"resource_url":"","status":1,"admin_url":"","audit_periodicity":360,"authorizable":0,"name":"project-1130"},"pathV2":[{"key":"key","val":"project-1130"}],"locationV2":[{"key":"r","val":"cn"}],"is_enable":true,"create_time":1649325694850133000},"resource_type":2,"expire_time":0,"condition":[{}],"is_enable":true,"is_break_glass":false,"information":{},"create_time":1649522420058047000}],"base_policy":null}}`

var kaniConverted string = `{"permission_types":[{"action":"read","name":"Read-only"},{"action":"write","name":"Operation"},{"action":"admin","name":"project-1130_admin"}],"user_permissions":[{"name":"Jakub Daliga","email":"jakub.daliga@bytedance.com","permission_target":"user","permissions":[{"action":"admin","name":"project-1130_admin","temporality":"permanent"}]}]}`

func TestProjectsPermissionsConversion(t *testing.T) {
    file, err := os.ReadFile("./testdata/policy_read_response.json")
    require.NoError(t, err)
    var resp *protego.PolicyResponse
    err = json.Unmarshal(file, &resp)
    require.NoError(t, err)

    var sorted map[int]*permissions.ProjectPermissions
    file, err = os.ReadFile("./testdata/project_permissions_response.json")
    require.NoError(t, err)
    err = json.Unmarshal(file, &sorted)
    require.NoError(t, err)

    cl := &httpClient{}
    assert.NoError(t, err)
    conv := cl.projectsPermissions(resp)
    convStr, err := json.MarshalToString(conv)
    assert.NoError(t, err)
    sortedStr, err := json.MarshalToString(sorted)
    assert.NoError(t, err)
    assert.JSONEq(t, sortedStr, convStr, "conv = %s", convStr)
}

func Test_httpClient_UserAccess(t *testing.T) {
    type fields struct {
        protegoClient    protego.Client
        ownershipDecider permissions.ParcelOwnershipDecider
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    permissions.UserAccess
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := &httpClient{
                protegoClient:    tt.fields.protegoClient,
                ownershipDecider: tt.fields.ownershipDecider,
            }
            got, err := h.UserAccess(tt.args.ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("httpClient.UserAccess() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("httpClient.UserAccess() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestHttpClientuserAccess(t *testing.T) {
    type fields struct {
        protegoClient    protego.Client
        ownershipDecider permissions.ParcelOwnershipDecider
    }
    type args struct {
        ctx                     context.Context
        batchEntityRespJSONPath string
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   permissions.UserAccess
    }{
        {
            name: "multiple member for one person",
            args: args{
                ctx:                     context.Background(),
                batchEntityRespJSONPath: "./testdata/batch_entity_response.json",
            },
            want: permissions.UserAccess{
                {
                    ProjectID:       1337,
                    UserPermissions: []permissions.PermissionType{permissions.PermissionTypeRead, permissions.PermissionTypeWrite},
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := &httpClient{
                protegoClient:    tt.fields.protegoClient,
                ownershipDecider: tt.fields.ownershipDecider,
            }
            file, err := os.ReadFile(tt.args.batchEntityRespJSONPath)
            require.NoError(t, err)
            var resp *protego.BatchAllowedEntityResponse
            err = json.Unmarshal(file, &resp)
            require.NoError(t, err)

            got := h.userAccess(tt.args.ctx, resp.Data.BatchResponses[0].Workload.AllowedTuples)
            assert.EqualValues(t, tt.want[0].ProjectID, got[0].ProjectID)
            assert.ElementsMatch(t, tt.want[0].UserPermissions, got[0].UserPermissions)
            assert.EqualValues(t, 1, len(got))
        })
    }
}
