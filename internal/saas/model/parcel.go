package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

/*
CREATE TABLE `project` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `name` text COMMENT 'name of project',
  `service_tree_node` varchar(512) DEFAULT NULL COMMENT 'ID of service tree node this project is bound to',
  `service_tree_path` varchar(512) DEFAULT NULL COMMENT 'name of service tree node this project is bound to',
  `accepted` tinyint(1) DEFAULT '0' COMMENT 'if project was accepted by SaaS admin',
  `modification` text COMMENT 'pending project modification',
  `state` text COMMENT 'project state',
  `state_comment` text COMMENT 'comment for the current project state',
  `appid` varchar(128) DEFAULT NULL COMMENT 'appid of project',
  `workflow_id` text COMMENT 'currently pending workflow id',
  `alert_user_emails` text COMMENT 'emails of users in the project to alert on events',
  `deactivated` tinyint(1) DEFAULT '0' COMMENT 'if project was deactivated',
  `kingdom` varchar(256) DEFAULT NULL COMMENT 'system from which the project originates',
  `description` text COMMENT 'project description, required by VolcEngine',
  `created_timestamp` bigint(20) NOT NULL DEFAULT '0' COMMENT 'timestamp of project creation, required by VolcEngine',
  `top_account_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'TOP AccountId that owns the project, required by VolcEngine',
  `is_test` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'if project is for test',
  `platform` varchar(512) NOT NULL DEFAULT '' COMMENT 'platform this project is bound to; empty means default platform for each kingdom',
  `volc_project_name` varchar(255) DEFAULT NULL COMMENT '火山项目名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_appid` (`appid`),
  UNIQUE KEY `uk_top_account_volc_project` (`top_account_id`,`volc_project_name`) COMMENT '确保同一账号下火山项目名称唯一',
  KEY `idx_top_account_id` (`top_account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=286752 DEFAULT CHARSET=utf8 COMMENT='projects'
*/

// Parcel groups all Species, and is the most general type of object that Garden user can interact with.
type Parcel struct {
    ID    int    `gorm:"column:id;primary_key;AUTO_INCREMENT" db:"id" json:"id"`
    AppID string `gorm:"column:appid" db:"appid" json:"appid"`

    // TopAccountID is AccountID of top account that owns this Parcel.
    //
    // As of 19.11.2021, AccountID has a numeric value. Because no numeric
    // operations make sense on IDs, we store it as a string.
    //
    // This field is not empty only if the Parcel was created by an outside system such as
    // VolcanoEngine or BytePlus. For Parcels created through  SaaS, this field should be equal to "".
    TopAccountID string `gorm:"column:top_account_id" db:"top_account_id" json:"top_account_id,omitempty"`

    Name             string      `gorm:"column:name" db:"name" json:"name,omitempty"`
    Description      string      `gorm:"column:description" db:"description" json:"description,omitempty"`
    ServiceTreeNode  string      `gorm:"column:service_tree_node" db:"service_tree_node" json:"service_tree_node,omitempty"`
    ServiceTreePath  string      `gorm:"column:service_tree_path" db:"service_tree_path" json:"service_tree_path,omitempty"`
    CreatedTimestamp int64       `gorm:"column:created_timestamp" db:"created_timestamp" json:"created_timestamp,omitempty"`
    IsTest           *bool       `gorm:"column:is_test" db:"is_test" json:"test_project"`
    Kingdom          string      `gorm:"column:kingdom" db:"kingdom" json:"system"`
    State            ParcelState `gorm:"column:state" db:"state" json:"state"`
    Deactivated      *bool       `gorm:"column:deactivated" db:"deactivated" json:"deactivated,omitempty"`
    AlertUserEmails  StringSlice `gorm:"column:alert_user_emails" db:"alert_user_emails" json:"alert_user_emails"`

    // backwards compatibility
    Accepted *bool `gorm:"column:accepted" db:"accepted" json:"accepted"`

    StateComment           *string             `gorm:"column:state_comment" db:"state_comment" json:"state_comment,omitempty"`
    PendingSignboardChange *ParcelModification `gorm:"column:modification" db:"modification" json:"modification,omitempty"`
    PendingInquiryID       *string             `gorm:"column:workflow_id" db:"workflow_id" json:"workflow_id,omitempty"`
    VolcanoProjectName     string              `gorm:"column:volc_project_name" db:"volc_project_name" json:"volc_project_name,omitempty"`

    // Platform is the platform that this parcel belongs to; default is "", which means
    // default platform for each kingdom. Nonempty value means this parcel is not limited
    // to total limits of appid on volcano's front-end page and is not displayed on the page either.
    Platform ParcelPlatform `gorm:"column:platform" db:"platform" json:"platform,omitempty"`

    // TODO add
    // UserPermissions []permissions.Type `gorm:"-" json:"user_permissions"`
    // Padlocks        []Padlock          `gorm:"-" json:"grants,omitempty"`
}

type ParcelModification ParcelSignboard

type ParcelSignboard struct {
    Name            string         `db:"name" json:"name,omitempty"`
    Description     string         `db:"description" json:"description,omitempty"`
    ServiceTreeNode string         `db:"service_tree_node" json:"service_tree_node,omitempty"`
    ServiceTreePath string         `db:"service_tree_path" json:"service_tree_path,omitempty"`
    Platform        ParcelPlatform `db:"platform" json:"platform,omitempty"`

    OptionalAppID   string   `db:"-" json:"appid,omitempty"`
    OptionalCreator string   `db:"-" json:"creator,omitempty"`
    OptionalOwners  []string `db:"-" json:"owners,omitempty"`
}

func (m ParcelModification) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *ParcelModification) Scan(src any) error {
    srcBytes, ok := src.([]byte)
    if !ok {
        panic(fmt.Sprintf("invalid ParcelModification.Scan source type: %T", src))
    }

    if len(srcBytes) == 0 {
        return nil
    }

    return json.Unmarshal(srcBytes, m)
}

const (
    ParcelStateCreatePending ParcelState = "create_pending" // 创建待处理
    ParcelStateUpdatePending ParcelState = "update_pending" // 更新待处理
    ParcelStateAccepted      ParcelState = "accepted"       // 接受
    ParcelStateCreateDenied  ParcelState = "create_denied"  // 创建被拒绝

    // NOTE: new parcel state added for admin-console

    ParcelStatePaused     ParcelState = "paused"      // 暂停
    ParcelStateDeleted    ParcelState = "deleted"     // 删除
    ParcelStateNotCreated ParcelState = "not_created" // 未创建
)

type ParcelState string

func (s *ParcelState) Scan(src any) error {
    if src == nil {
        *s = ""
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        panic(fmt.Sprintf("invalid ParcelState.Scan source type: %T", src))
    }

    *s = ParcelState(srcBytes)
    return nil
}

func (s ParcelState) Value() (driver.Value, error) { return string(s), nil }

const (
    // DefaultParcelPlatform is the default platform for each kingdom.
    DefaultParcelPlatform ParcelPlatform = ""

    // AvatarParcelPlatform is the platform all avatar projects belong to;
    // its value comes from the product name set in the volcano cost center.
    AvatarParcelPlatform ParcelPlatform = "AvatarAPI"

    // ArkParcelPlatform is the platform all ark tts projects belong to;
    // its value is negotiated with front-end.
    ArkParcelPlatform ParcelPlatform = "ArkAPI"

    OutboundCallParcelPlatform ParcelPlatform = "OutboundCallAPI"
    V3ExperienceAPI            ParcelPlatform = "V3ExperienceAPI"
)

type ParcelPlatform string

func (p *ParcelPlatform) Scan(src any) error {
    if src == nil {
        *p = ""
        return nil
    }

    switch value := src.(type) {
    case []byte:
        *p = ParcelPlatform(value)
    case string:
        *p = ParcelPlatform(value)
    default:
        return fmt.Errorf("invalid ParcelPlatform.Scan source type: %T", src)
    }

    return nil
}

func (p ParcelPlatform) Value() (driver.Value, error) { return string(p), nil }
