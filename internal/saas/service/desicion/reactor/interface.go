package reactor

import (
    "context"

    "study/internal/saas/model"
)

type Accountant interface{}

type VolcanoPlantInstanceRegistryFactory interface{}

type OrderCollectorFactory interface{}

type order interface {
    OrderCallback(ctx context.Context, order model.Order, status model.InstanceChangeStatus) error
    TerminateOrder(ctx context.Context, instanceNO string) error
    IsPrimaryInstance(instanceNo string) bool
    InstanceStatus(ctx context.Context, instanceNo string) (*model.Instance, error)
    SubOrderDetail(ctx context.Context, subOrderNo string) (*model.SubOrder, error)
}

type locksmith interface{}

type purchaseConfigLoader interface{}

type noticeSettingsGetter interface{}

type orderQuotaAlarm interface{}

type cladeRegistry interface{}

type parcelRegistry interface{}

type speciesLoader interface{}

type plantRegistry interface{}

type harvestLimitsEnforcer interface {
    Enforce(ctx context.Context, parcelID int, resourceID string, limits model.HarvestLimitations) error
}

type resourceReplacer interface {
    Replace(ctx context.Context, old, new model.ServiceInstanceNumber) error
}

type resourceTemplateGetter interface {
    Get(ctx context.Context, rp *model.ResourcePackage) (*model.ResourcePackTemplateWithConfig, error)
}

type resourceGromRegistry interface{}

type ResourceRegistry interface{}

type ResourceGroupLoader interface{}

type VolcanoResourceForester interface{}
