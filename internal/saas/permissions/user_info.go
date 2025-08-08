package permissions

import (
    "context"
    "fmt"

    // "code.byted.org/lab-speech/saas_backend/clients/lark"
    "github.com/samber/lo"
)

// UserAdditional includes additional user info meaningful for display on web page.
type UserAdditional struct {
    EnName string      `json:"en_name,omitempty"`
    CnName string      `json:"cn_name,omitempty"`
    Avatar lark.Avatar `json:"avatar,omitempty"`
}

// FillUserAdditionals fills UserAdditional if possible inside ProjectPermissions.
func FillUserAdditionals(ctx context.Context, larkClient lark.UserInfoClient, pps ...*ProjectPermissions) error {
    if len(pps) == 0 || larkClient == nil {
        return fmt.Errorf("found invalid arg: pp %+v, client %+v", pps, larkClient)
    }
    emails := lo.Uniq(
        lo.Flatten(
            lo.FilterMap(pps, func(item *ProjectPermissions, _ int) ([]string, bool) {
                if lo.IsEmpty(item) {
                    return nil, false
                }

                return lo.FilterMap(item.PermissionOwners, func(p PermissionHolder, _ int) (string, bool) {
                    return p.User.Email, lo.IsNotEmpty(p.User.Email)
                }), true
            }),
        ),
    )
    if len(emails) == 0 {
        return nil // nothing to fill
    }
    userAdditionals, err := larkClient.BatchGetUsers(ctx, emails...)
    if err != nil {
        return fmt.Errorf("get user additional: %w", err)
    }
    userEmailToAdditional := lo.SliceToMap(userAdditionals, func(ui *lark.UserInfo) (string, *lark.UserInfo) {
        return ui.Email, ui
    })
    for _, pp := range pps {
        if pp == nil || len(pp.PermissionOwners) == 0 {
            continue
        }
        var newPermissionHolders PermissionHolders
        for _, po := range pp.PermissionOwners {
            newPO := PermissionHolder{
                User:        po.User,
                Target:      po.Target,
                Permissions: po.Permissions,
            }
            if additional, ok := userEmailToAdditional[po.User.Email]; ok {
                newPO.UserAdditional = &UserAdditional{
                    EnName: additional.EnName,
                    CnName: additional.Name,
                    Avatar: additional.Avatar,
                }
            }
            newPermissionHolders = append(newPermissionHolders, newPO)
        }
        pp.PermissionOwners = newPermissionHolders
        pp.Sort()
    }
    return nil
}
