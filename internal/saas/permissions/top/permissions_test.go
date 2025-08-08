package top

import (
    "context"
    "testing"

    "code.byted.org/lab-speech/saas_backend/commons/more_testing"
    "code.byted.org/lab-speech/saas_backend/gardening"
    "code.byted.org/lab-speech/saas_backend/permissions"
    "code.byted.org/lab-speech/saas_backend/volcengine/top"
)

type stubParcelRegistry struct {
    selectParcelsFunc func(ctx context.Context, filter gardening.ParcelFilter) ([]gardening.Parcel, error)
    updateParcelFunc  func(ctx context.Context, parcel gardening.Parcel, filter gardening.ParcelFilter) (gardening.Parcel, error)
}

func (spr stubParcelRegistry) SelectParcels(ctx context.Context, filter gardening.ParcelFilter) ([]gardening.Parcel, error) {
    return spr.selectParcelsFunc(ctx, filter)
}

func (spr stubParcelRegistry) UpdateParcel(ctx context.Context, parcel gardening.Parcel, filter gardening.ParcelFilter) (gardening.Parcel, error) {
    return spr.updateParcelFunc(ctx, parcel, filter)
}

func Test_PermissionManager_UserAccess(t *testing.T) {
    // given
    ctx := top.CtxWithAccount(context.Background(), top.Account{AccountID: 2137})
    expectedFilter := gardening.ParcelFilter{TOPAccountIDs: []string{"2137"}}
    allowedParcels := []gardening.Parcel{{ID: 420}, {ID: 1337}, {ID: 69}}
    parcelRegistry := stubParcelRegistry{
        selectParcelsFunc: func(ctx context.Context, filter gardening.ParcelFilter) ([]gardening.Parcel, error) {
            more_testing.ErrorIfNotEqual(t, "invalid parcel filter", expectedFilter, filter)
            return allowedParcels, nil
        },
    }
    permissionManager := NewPermissionManager(parcelRegistry)
    // when
    userAccess, err := permissionManager.UserAccess(ctx)
    // then
    more_testing.FailOnError(t, "", err)
    ownerPermissions := []permissions.PermissionType{permissions.PermissionTypeOwner}
    expectedUserAccess := permissions.UserAccess{
        {ProjectID: 420, UserPermissions: ownerPermissions},
        {ProjectID: 1337, UserPermissions: ownerPermissions},
        {ProjectID: 69, UserPermissions: ownerPermissions},
    }
    more_testing.ErrorIfNotEqual(t, "invalid user access returned", expectedUserAccess, userAccess)
}

func Test_PermissionManager_InitProjectPermissions(t *testing.T) {
    // given
    ctx := top.CtxWithAccount(context.Background(), top.Account{AccountID: 2137})
    updateCalled := false
    parcelRegistry := stubParcelRegistry{
        updateParcelFunc: func(ctx context.Context, parcel gardening.Parcel, filter gardening.ParcelFilter) (gardening.Parcel, error) {
            updateCalled = true
            expectedUpdate := gardening.Parcel{TOPAccountID: "2137"}
            more_testing.ErrorIfNotEqual(t, "invalid update", expectedUpdate, parcel)
            expectedFilter := gardening.ParcelFilter{IDs: []int{420}}
            more_testing.ErrorIfNotEqual(t, "invalid filter", expectedFilter, filter)
            // ignore returned parcel, because it's ignored in the caller
            return gardening.Parcel{}, nil
        },
    }
    permissionManager := NewPermissionManager(parcelRegistry)
    // when
    err := permissionManager.InitProjectPermissions(ctx, 420)
    // then
    more_testing.FailOnError(t, "", err)
    more_testing.ErrorIfNotEqual(t, "registry update not called", true, updateCalled)
}
