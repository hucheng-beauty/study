package permission

import (
    "context"
    "fmt"

    "github.com/samber/lo"
)

// UserAdditional includes additional user info meaningful for display on web page.
type UserAdditional struct {
    EnName string      `json:"en_name,omitempty"`
    CnName string      `json:"cn_name,omitempty"`
    Avatar lark.Avatar `json:"avatar,omitempty"`
}

// FillUserAdditionals fills UserAdditional if possible inside ProjectPermissions.
func FillUserAdditionals1(ctx context.Context,
        larkClient lark.UserInfoClient, project ...*Project) error {
    return nil
}

func FillUserAdditional(ctx context.Context, larkClient lark.UserInfoClient, projects ...*Project) error {
    if len(projects) == 0 || larkClient == nil {
        return fmt.Errorf("found invalid arg: projects %+v, client %+v", projects, larkClient)
    }

    emails := lo.Uniq(
        lo.Flatten(
            lo.FilterMap(projects, func(item *Project, _ int) ([]string, bool) {
                if lo.IsEmpty(item) {
                    return nil, false
                }

                return lo.FilterMap(item.Owners,
                    func(h Holder, _ int) (string, bool) {
                        return h.User.Email, lo.IsNotEmpty(h.User.Email)
                    },
                ), true
            })))
    if len(emails) == 0 {
        return nil // nothing to fill
    }
    users, err := larkClient.BatchGetUsers(ctx, emails...)
    if err != nil {
        return fmt.Errorf("get user additional: %w", err)
    }

    userEmailToAdditional := lo.SliceToMap(users,
        func(ui *lark.UserInfo) (string, *lark.UserInfo) {
            return ui.Email, ui
        },
    )

    for _, project := range projects {
        if project == nil || len(project.Owners) == 0 {
            continue
        }

        var holders Holders
        for _, po := range project.Owners {
            holder := Holder{
                User:        po.User,
                Target:      po.Target,
                Permissions: po.Permissions,
            }

            if additional, ok := userEmailToAdditional[po.User.Email]; ok {
                holder.UserAdditional = &UserAdditional{
                    EnName: additional.EnName,
                    CnName: additional.Name,
                    Avatar: additional.Avatar,
                }
            }
            holders = append(holders, holder)
        }

        project.Owners = holders
        project.Sort()
    }
    return nil
}
