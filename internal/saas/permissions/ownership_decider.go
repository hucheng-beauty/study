package permissions

import (
    "context"
)

// ParcelOwnershipDecider decides who should be the owner of given parcel.
type ParcelOwnershipDecider interface {
    ParcelOwnership(ctx context.Context) (ParcelOwnership, error)
}

type ParcelOwnership struct {
    Creator string
    Owners  []string
}
