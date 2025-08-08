package top

import (
    "context"
    "fmt"
    "strconv"

    // "code.byted.org/lab-speech/saas_backend/gardening"
    // "code.byted.org/lab-speech/saas_backend/permissions"
    // "code.byted.org/lab-speech/saas_backend/volcengine/top"

    "study/internal/saas/permissions"
)

type parcelRegistry interface {
    SelectParcels(ctx context.Context, filter gardening.ParcelFilter) ([]gardening.Parcel, error)
    UpdateParcel(ctx context.Context, parcel gardening.Parcel, filter gardening.ParcelFilter) (gardening.Parcel, error)
}

type PermissionManager struct {
    parcelRegistry parcelRegistry
}

func NewPermissionManager(parcelRegistry parcelRegistry) *PermissionManager {
    return &PermissionManager{parcelRegistry: parcelRegistry}
}

func (pm PermissionManager) UserAccess(ctx context.Context) (permissions.UserAccess, error) {
    account, err := top.AccountFromCtx(ctx)
    if err != nil {
        return permissions.UserAccess{}, fmt.Errorf("get top account from context: %w", err)
    }
    filter := gardening.ParcelFilter{TOPAccountIDs: []string{strconv.FormatInt(account.AccountID, 10)}}
    parcels, err := pm.parcelRegistry.SelectParcels(ctx, filter)
    if err != nil {
        return permissions.UserAccess{}, fmt.Errorf("[top_acount_id=%v] get parcels for top account: %w",
            account.AccountID, err)
    }
    var userAccess permissions.UserAccess
    for _, parcel := range parcels {
        userAccess = append(userAccess, permissions.UserAccessControl{
            ProjectID:       parcel.ID,
            UserPermissions: []permissions.PermissionType{permissions.PermissionTypeOwner},
        })
    }
    return userAccess, nil
}

func (pm PermissionManager) InitProjectPermissions(ctx context.Context, projectID int) error {
    account, err := top.AccountFromCtx(ctx)
    if err != nil {
        return fmt.Errorf("get top account from context: %w", err)
    }
    update := gardening.Parcel{TOPAccountID: strconv.FormatInt(account.AccountID, 10)}
    filter := gardening.ParcelFilter{IDs: []int{projectID}}
    _, err = pm.parcelRegistry.UpdateParcel(ctx, update, filter)
    if err != nil {
        return fmt.Errorf("update parcel: %w", err)
    }
    return nil
}

// InitProjectPermissionsExplicit has the same meaning as InitProjectPermissions in the realm of top.
func (pm PermissionManager) InitProjectPermissionsExplicit(ctx context.Context, projectID int, creator string, owners []string) error {
    return pm.InitProjectPermissions(ctx, projectID)
}
