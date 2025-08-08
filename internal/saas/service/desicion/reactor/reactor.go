package reactor

import (
    "context"
    "fmt"

    "study/internal/saas/model"
    "study/internal/saas/service/desicion/event"
    "study/internal/saas/service/desicion/trade"

    "github.com/benbjohnson/clock"
)

var _ trade.OrderHandler = (*Reactor)(nil)

type Reactor struct {
    merchant                order
    orderQuotaAlarm         orderQuotaAlarm
    purchaseConfigLoader    purchaseConfigLoader
    noticeSettingsGetter    noticeSettingsGetter
    locksmith               locksmith
    accountant              Accountant
    orderCollectorFactory   OrderCollectorFactory
    instanceRegistryFactory VolcanoPlantInstanceRegistryFactory

    cladeRegistry          cladeRegistry
    parcelRegistry         parcelRegistry
    speciesLoader          speciesLoader
    plantRegistry          plantRegistry
    enforcer               harvestLimitsEnforcer
    resourceTemplateGetter resourceTemplateGetter
    resourceGormRegistry   resourceGromRegistry
    resourceReplacer       resourceReplacer
    resourceRegistry       ResourceRegistry
    resourceGroupLoader    ResourceGroupLoader
    resourceForester       VolcanoResourceForester

    register *register
    events   *event.Events
    clock    clock.Clock
}

func (r *Reactor) initRegister() {
    r.register = &register{
        root: newStrategy(
            "root", passAll,
            r.registerOrderNew(),
            r.registerOrderFormalize(),
            r.registerOrderModify(),
            r.registerOrderExpiring(),
            r.registerOrderExpire(),
            r.registerOrderOverdue(),
            r.registerOrderRenew(),
            r.registerOrderResume(),
            r.registerOrderTerminate(),
            r.registerOrderReclaim(),
            r.registerOrderCancelNew(),
            r.registerOrderCancelModify(),
            r.registerOrderCancelFormalize(),
            r.registerOrderPreNew(),
            r.registerOrderPreModify(),
            r.registerOrderPreTerminate(),
            newAction("unknown_order_type", passAll,
                func(ctx context.Context, orderContext *orderContext) error {
                    // return fmt.Errorf("handle order, unknown order type: %v", orderContext.order.Type)
                    return fmt.Errorf("handle order, unknown order type: %v", nil)
                },
            ),
        ),
    }
}

func (r *Reactor) HandleOrder(ctx context.Context, order model.Order) (err error) {
    return nil
}
