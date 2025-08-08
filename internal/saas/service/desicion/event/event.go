package event

import (
    "study/internal/saas/model"
    "study/internal/saas/toolkit/pubsub"
)

// NOTE: reactor events contains events that third-party (except direct db client)
// should listen to after rds entries changed by the reactor itself.
// Currently, only prepaid service & prepaid resource support these events.

type (
    BaseArg struct {
        Order *model.Order
    }

    ParcelArg struct {
        BaseArg `json:",inline"`
        Parcel  *model.Parcel

        ResourceTag *model.ResourceTag `json:"ResourceTag"` // optional
    }
    PlantArg struct {
        BaseArg `json:",inline"`
        Plant   *model.Plant

        Parcel      *model.Parcel      // optional
        ResourceTag *model.ResourceTag `json:"ResourceTag"` // optional
    }
    UpdatePlantArg PlantArg // 更新

    ExpirePlantArg TerminatePlantArg // 过期

    // TerminatePlantArg could be from trial plant or service plant.
    TerminatePlantArg struct {
        BaseArg `json:",inline"`

        Parcel *model.Parcel

        // PlantBeforeTermination may have data more useful
        // such as it previous state, and for this reason
        // event arg here requires a plant before change.
        PlantBeforeTermination *model.Plant

        // NOTE: ResourcesBeforeTermination is deprecated since
        // resource pack is terminated or expired using other events
        // ResourcesBeforeTermination includes all resource packs
        // expired along with it and owned by this plant;
        // NOT INCLUDING standalone packet if there's any.
        ResourcesBeforeTermination []*model.ResourcePackage
    }                        // 终止
    ResumePlantArg  PlantArg // 恢复
    ReclaimPlantArg PlantArg // 回收

    NewResourcesArg struct {
        BaseArg               `json:",inline"`
        PlantWithUpdatedQuota *model.Plant
        Resources             []*model.Resource

        ResourceTag *model.ResourceTag `json:"ResourceTag"` // optional
    }
    ExpireResourceArg TerminateResourceArg // 过期
    RenewResourceArg  struct {
        BaseArg `json:",inline"`

        PlantWithUpdatedQuota *model.Plant
        ResourcePackUpdated   *model.Resource
    } // 续约
    TerminateResourceArg struct {
        BaseArg `json:",inline"`

        PlantWithUpdatedQuota *model.Plant
        ResourcesExpired      []*model.Resource
    }                                       // 终止
    ReclaimResourceArg TerminateResourceArg // 回收
)

type Events struct {
    Parcel         pubsub.Event[*ParcelArg]
    Plant          pubsub.Event[*PlantArg]
    UpdatePlant    pubsub.Event[*UpdatePlantArg]    // 更新
    ExpirePlant    pubsub.Event[*ExpirePlantArg]    // 过期
    ResumePlant    pubsub.Event[*ResumePlantArg]    // 恢复
    TerminatePlant pubsub.Event[*TerminatePlantArg] // 终止
    ReclaimPlant   pubsub.Event[*ReclaimPlantArg]   // 回收

    NewResources      pubsub.Event[*NewResourcesArg]
    ExpireResource    pubsub.Event[*ExpireResourceArg]    // 过期
    RenewResource     pubsub.Event[*RenewResourceArg]     // 续约
    TerminateResource pubsub.Event[*TerminateResourceArg] // 终止
    ReclaimResource   pubsub.Event[*ReclaimResourceArg]   // 回收
}

func NewEvents() *Events {
    e := &Events{}

    // service
    e.Parcel = pubsub.New[*ParcelArg]()
    e.Plant = pubsub.New[*PlantArg]()
    e.UpdatePlant = pubsub.New[*UpdatePlantArg]()
    e.ExpirePlant = pubsub.New[*ExpirePlantArg]()
    e.ResumePlant = pubsub.New[*ResumePlantArg]()
    e.TerminatePlant = pubsub.New[*TerminatePlantArg]()
    e.ReclaimPlant = pubsub.New[*ReclaimPlantArg]()

    // resource package
    e.NewResources = pubsub.New[*NewResourcesArg]()
    e.ExpireResource = pubsub.New[*ExpireResourceArg]()
    e.RenewResource = pubsub.New[*RenewResourceArg]()
    e.TerminateResource = pubsub.New[*TerminateResourceArg]()
    e.ReclaimResource = pubsub.New[*ReclaimResourceArg]()

    return e
}

// Clear clears all subscribers.
// Usually used in tests.
func (e *Events) Clear() {
    if e == nil {
        return
    }

    e.clear()
}

func (e *Events) clear() {
    e.Parcel.Clear()
    e.Plant.Clear()
    e.UpdatePlant.Clear()
    e.ExpirePlant.Clear()
    e.ResumePlant.Clear()
    e.TerminatePlant.Clear()
    e.ReclaimPlant.Clear()

    e.NewResources.Clear()
    e.ExpireResource.Clear()
    e.RenewResource.Clear()
    e.TerminateResource.Clear()
    e.ReclaimResource.Clear()
}
