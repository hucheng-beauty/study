package model

import (
    "database/sql/driver"
    "fmt"
)

/*
CREATE TABLE `service` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `level` varchar(128) DEFAULT NULL COMMENT 'level of service - business or trial',
  `project_id` int(11) DEFAULT NULL COMMENT 'id of project under which service operates',
  `service_blueprint_id` int(11) DEFAULT NULL COMMENT 'blueprint for service',
  `pending_quotas` text COMMENT 'pending quota update (waiting for accept in BPM)',
  `name` text COMMENT 'name of user service',
  `active_quotas` text COMMENT 'Active quota limits for service',
  `uprooted` tinyint(1) DEFAULT '0' COMMENT 'set to true instead of deleting service on service deactivation, false means service is active',
  `state` varchar(128) DEFAULT NULL COMMENT 'state of service, enumerated in code',
  `volc_instance_number` varchar(256) DEFAULT NULL COMMENT 'Hourly bill InstanceNO from VolcEngine',
  `volc_instance_number_monthly` varchar(256) DEFAULT NULL COMMENT 'Monthly bill InstanceNO from VolcEngine',
  `volc_instance_type_preferred` varchar(256) NOT NULL DEFAULT '' COMMENT 'preferred instance type to bill set from admin-console, empty string means monthly if possible',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time of entry',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  `volc_instance_exiled` json DEFAULT NULL COMMENT 'list of exiled volcano instances, item should at least include the field volc_instance_number',
  `is_standalone` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'marks if this service is a pseudo service created just for its standalone resource pack, if true then their state change should sync with each other',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_id_service_blueprint_id` (`project_id`,`service_blueprint_id`),
  KEY `idx_blueprint_id` (`service_blueprint_id`),
  KEY `idx_volc_instance_number` (`volc_instance_number`),
  KEY `idx_volc_instance_number_monthly` (`volc_instance_number_monthly`)
) ENGINE=InnoDB AUTO_INCREMENT=1736944 DEFAULT CHARSET=utf8 COMMENT='blueprints enabled for projects'
*/

type Plant struct {
    ID       int     `gorm:"column:id;primary_key;AUTO_INCREMENT" db:"id" json:"id"`
    Name     *string `gorm:"column:name" db:"name" json:"name"`
    ParcelID int     `gorm:"column:project_id" db:"project_id" json:"project_id"`

    // Uprooted should be set to `true` if service was already created, but was later uprooted.
    // If the service is active, it should be set to `false`.
    Uprooted *bool `gorm:"column:uprooted" db:"uprooted" json:"inactive"`

    State                PlantState                 `gorm:"column:state" db:"state" json:"state"`
    Rank                 Rank                       `gorm:"column:level" db:"level" json:"level"`
    HarvestLimits        HarvestLimitations         `gorm:"column:active_quotas" db:"active_quotas" json:"active_quotas"`
    PendingHarvestLimits *PendingHarvestLimitations `gorm:"column:pending_quotas" db:"pending_quotas" json:"pending_quotas,omitempty"`

    // NOTE: VolcanoEngine_xxx fields below should be managed by `VolcanoPlantInstanceRegistrar`
    // to sync with db automatically

    // VolcanoEngineInstanceNo is used to bill users hourly. If VolcanoEngineMonthlyInstanceNo is empty,
    // billing is pushed to this instance.
    VolcanoEngineInstanceNo *string `gorm:"column:volc_instance_number" db:"volc_instance_number" json:"volc_instance_number,omitempty"`

    // VolcanoEngineMonthlyInstanceNo is used to bill users monthly. It's empty by default.
    // If a user is migrated to monthly billing, this field is filled in and billing is pushed
    // to this instance instead of VolcanoEngineInstanceNo.
    // I repeat: if this field is not empty or nil, BILLING WILL NOT BE PUSHED to VolcanoEngineInstanceNo,
    // but it will be pushed to VolcanoEngineMonthlyInstanceNo instead.
    VolcanoEngineMonthlyInstanceNo *string `gorm:"column:volc_instance_number_monthly" db:"volc_instance_number_monthly" json:"volc_instance_number_monthly,omitempty"`

    // VolcanoInstanceTypePreferred is preferred instance type to
    // bill set from admin-console, empty string means monthly if possible.
    VolcanoInstanceTypePreferred string `gorm:"column:volc_instance_type_preferred" db:"volc_instance_type_preferred" json:"volc_instance_type_preferred,omitempty"`

    // IsStandalone notice that it cannot be nil in db;
    // there we use pointer only to allow form filter to pick it up.
    IsStandalone *bool `gorm:"column:is_standalone" db:"is_standalone" json:"is_standalone,omitempty"`

    SpeciesID int `gorm:"column:service_blueprint_id" db:"service_blueprint_id" json:"-" log:"service_blueprint_id,omitempty"`
    Species   `gorm:"-" db:"service_blueprint" json:"blueprint" log:"-"`

    // 记录关联的服务包信息,避免多传参数
    BundleSpeciesList []Species `gorm:"-" json:"-" log:"-"`
}

const (
    // service state

    PlantStateCreating   PlantState = "creating"   // 创建中
    PlantStateActive     PlantState = "active"     // 活跃的
    PlantStateUpgrading  PlantState = "upgrading"  // 升级中
    PlantStatePaused     PlantState = "paused"     // 暂停
    PlantStateInactive   PlantState = "inactive"   // 不活跃
    PlantStateExpired    PlantState = "expired"    // 过期
    PlantStateOverdue    PlantState = "overdue"    // 逾期
    PlantStateTerminated PlantState = "terminated" // 终止
    PlantStateReclaimed  PlantState = "reclaimed"  // 回收

    // service level

    RankTrial    Rank = "trial"    // 试用
    RankBusiness Rank = "business" // 商业
)

type PlantState string

func (ps *PlantState) Scan(src any) error {
    if src == nil {
        *ps = ""
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        panic(fmt.Sprintf("invalid PlantState.Scan source type: %T", src))
    }

    *ps = PlantState(srcBytes)

    return nil
}

func (ps PlantState) Value() (driver.Value, error) { return string(ps), nil }

type Rank string

func (r *Rank) Scan(src any) error {
    srcBytes, ok := src.([]byte)
    if !ok {
        panic(fmt.Sprintf("invalid Rank.Scan source type: %T", src))
    }

    *r = Rank(srcBytes)

    return nil
}

func (r Rank) Value() (driver.Value, error) { return string(r), nil }
