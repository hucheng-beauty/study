package permissions

import (
    "context"
    "testing"

    // "code.byted.org/lab-speech/saas_backend/clients/lark"
    jsoniter "github.com/json-iterator/go"
    "github.com/stretchr/testify/assert"
)

func TestFillUserAdditionals(t *testing.T) {
    exampleEmails := []string{
        "a@example.com",
        "b@example.com",
    }
    exampleAvatar := lark.Avatar{
        Avatar72:     "avatar_72",
        Avatar240:    "avatar_240",
        Avatar640:    "avatar_640",
        AvatarOrigin: "avatar_origin",
    }
    exampleDescriptor := Descriptor{Type: PermissionTypeOwner, Name: "owner_type"}
    type args struct {
        larkClient lark.UserInfoClient
        pps        []*ProjectPermissions
        expected   []*ProjectPermissions
    }
    tests := []struct {
        name    string
        args    args
        wantErr bool
    }{
        {
            name: "happy path",
            args: args{
                larkClient: &stubLarkClient{
                    stubBatchGetUsers: func(ctx context.Context, emails ...string) ([]*lark.UserInfo, error) {
                        assert.EqualValues(t, exampleEmails, emails)
                        return []*lark.UserInfo{
                            {
                                Email:  "a@example.com",
                                Name:   "cn_a",
                                EnName: "en_a",
                                Avatar: exampleAvatar,
                            },
                        }, nil
                    },
                },
                pps: []*ProjectPermissions{
                    {
                        PermissionTypes: []Descriptor{exampleDescriptor},
                        PermissionOwners: []PermissionHolder{
                            {
                                User:        User{Name: "a", Email: "a@example.com"},
                                Target:      "target_a",
                                Permissions: Permissions{{Temporality: "t", Descriptor: exampleDescriptor, Expires: 1}},
                            },
                            {
                                User:        User{Name: "c", Email: ""}, // empty email no additional info
                                Target:      "target_c",
                                Permissions: Permissions{{Temporality: "t", Descriptor: exampleDescriptor, Expires: 1}},
                            },
                        },
                    },
                    nil, // should return no err
                    {
                        PermissionTypes: Descriptors{exampleDescriptor},
                        // empty owner should return no err
                    },
                    {
                        PermissionTypes: []Descriptor{exampleDescriptor},
                        PermissionOwners: []PermissionHolder{
                            {
                                // no additional if no return from lark client
                                User:        User{Name: "b", Email: "b@example.com"},
                                Target:      "target_b",
                                Permissions: Permissions{{Temporality: "t2", Descriptor: exampleDescriptor, Expires: 2}},
                            },
                        },
                    },
                },
                expected: []*ProjectPermissions{
                    {
                        PermissionTypes: []Descriptor{exampleDescriptor},
                        PermissionOwners: []PermissionHolder{
                            {
                                User: User{
                                    Name: "a", Email: "a@example.com",
                                    UserAdditional: &UserAdditional{
                                        EnName: "en_a",
                                        CnName: "cn_a",
                                        Avatar: exampleAvatar,
                                    },
                                },
                                Target:      "target_a",
                                Permissions: Permissions{{Temporality: "t", Descriptor: exampleDescriptor, Expires: 1}},
                            },
                            {
                                User:        User{Name: "c", Email: ""}, // empty email no additional info
                                Target:      "target_c",
                                Permissions: Permissions{{Temporality: "t", Descriptor: exampleDescriptor, Expires: 1}},
                            },
                        },
                    },
                    nil, // should return no err
                    {
                        PermissionTypes: Descriptors{exampleDescriptor},
                    },
                    {
                        PermissionTypes: []Descriptor{exampleDescriptor},
                        PermissionOwners: []PermissionHolder{
                            {
                                User:        User{Name: "b", Email: "b@example.com"},
                                Target:      "target_b",
                                Permissions: Permissions{{Temporality: "t2", Descriptor: exampleDescriptor, Expires: 2}},
                            },
                        },
                    },
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := FillUserAdditionals(context.Background(), tt.args.larkClient, tt.args.pps...); (err != nil) != tt.wantErr {
                t.Errorf("FillUserAdditionals() error = %v, wantErr %v", err, tt.wantErr)
            }
            exStr, err := jsoniter.MarshalToString(tt.args.expected)
            assert.NoError(t, err)
            acStr, err := jsoniter.MarshalToString(tt.args.pps)
            assert.NoError(t, err)
            assert.JSONEq(t, exStr, acStr, acStr)
        })
    }
}

type stubLarkClient struct {
    stubBatchGetUsers func(ctx context.Context, emails ...string) ([]*lark.UserInfo, error)
}

func (s *stubLarkClient) BatchGetUsers(ctx context.Context, emails ...string) ([]*lark.UserInfo, error) {
    return s.stubBatchGetUsers(ctx, emails...)
}
