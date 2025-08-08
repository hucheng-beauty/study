package model

/*
CREATE TABLE `resource_packs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'table primary key',
  `configuration_code` varchar(256) DEFAULT NULL COMMENT 'resource pack configuration code',
  `instance_number` varchar(256) DEFAULT NULL COMMENT 'number of volcengine instance',
  `service_id` int(11) DEFAULT '0' COMMENT 'fk to service',
  `is_quota` tinyint(1) DEFAULT NULL COMMENT 'if it needs quota recalculation',
  `is_access` tinyint(1) DEFAULT NULL COMMENT 'if it needs resource_tree rebuild',
  `created` bigint(20) DEFAULT NULL COMMENT 'milliseconds timestamp of order placed',
  `purchased` bigint(20) DEFAULT NULL COMMENT 'milliseconds timestamp of purchase done, confirmation from MQ',
  `quota_type` varchar(256) DEFAULT NULL COMMENT 'type of resource purchased',
  `months` int(11) DEFAULT NULL COMMENT 'for how many months is resource pack available',
  `value` bigint(20) NOT NULL DEFAULT '0' COMMENT 'how much was purchased',
  `special_system_package` tinyint(1) DEFAULT NULL COMMENT 'true if business initial resource pack',
  `expires` bigint(20) DEFAULT NULL COMMENT 'milliseconds timestamp of when the instance should expire, taken from MQ callback',
  `is_given` tinyint(1) DEFAULT '0' COMMENT 'if it is given by operator',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time of entry',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  `pack_type` varchar(64) NOT NULL COMMENT 'resource pack type',
  `begins` bigint(20) NOT NULL DEFAULT '0' COMMENT 'ms timestamp of when the instance should begin, 0 if it should take effect immediately',
  `is_from_prepaid_service` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'marks if such resource pack is generated as part of prepaid service',
  `is_standalone` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'marks if this resource pack owns a pseudo serivce created by template, if true then their state change should sync with each other',
  `quota_type_meta` json DEFAULT NULL COMMENT 'meta data for quota type',
  `attributes` json DEFAULT NULL COMMENT 'key-value pairs only meaningful for certain pack type',
  `alias` varchar(512) DEFAULT '' COMMENT 'client preferred name of a resource pack',
  `state` varchar(20) DEFAULT '' COMMENT 'resource pack state',
  `top_account_id` varchar(20) DEFAULT '0' COMMENT 'volc account id',
  `volc_product` varchar(128) DEFAULT '' COMMENT 'volc product name',
  `train_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT 'train id, like speaker',
  `volc_project_name` varchar(255) NOT NULL DEFAULT '' COMMENT '火山项目名',
  PRIMARY KEY (`id`),
  KEY `idx_instance_number` (`instance_number`),
  KEY `idx_service_id` (`service_id`),
  KEY `idx_top_account_id` (`top_account_id`),
  KEY `idx_state` (`state`),
  KEY `idx_quota_type` (`quota_type`),
  KEY `idx_expires` (`expires`),
  KEY `idx_special_system_package` (`special_system_package`),
  KEY `idx_pack_type` (`pack_type`),
  KEY `idx_configuration_code` (`configuration_code`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_begins` (`begins`),
  KEY `idx_train_id` (`train_id`),
  KEY `idx_volc_project_name` (`volc_project_name`)
) ENGINE=InnoDB AUTO_INCREMENT=358967251257606 DEFAULT CHARSET=utf8 COMMENT='information on VolcEngine resource packs'
*/

// ResourcePackage is a template to generate a real volcano engine.
// it will at least be used in two place
//  1. when we send the volcano the final order
//     The final.ResourcePack will only use data here as a reference
//     It's their choice to generate any VolcanoEngine.NewResources
//  2. when we calculate the price we show in our console
type ResourcePackage struct {
    // Code is an optional field aims to identify each resource pack;
    // it is set when `ConfigurationCode` is alone not enough to differentiate
    // packs inside species; NEVER call it directly, use `PreferCode()` instead;
    // it is defined by SaaS, should never have changed after created.
    Code string `json:"code"`

    // ConfigurationCode is used to place volcano order; it must be
    // from the volcano product config directly; however, it may not
    // be the same for each pack inside species; it is defined by SaaS,
    // should never have changed after created.
    ConfigurationCode string `json:"configuration_code"`

    ChargeItemCode  string         `json:"charge_item_code"` // should never have changed after created
    Period          string         `json:"period"`           // we expect it should not be changed after created, because user may see different price before order & after order
    NeverExpire     bool           `json:"never_expire"`     // never expire, be carefully to set, only set, it will ignore ExpireQuantity use long time as expiration time
    ExpireQuantity  int            `json:"expireQuantity"`   // $expireQuantity $Period is the time this package will expire, e.g. 12 monthly
    ResourceDisplay string         `json:"resource_display"` // used in ConfigDescription and InstanceName, change carefully
    PriceDisplay    string         `json:"price_display"`
    DurationDisplay string         `json:"duration_display"`
    DisplayAsBought bool           `json:"display_as_bought"`
    Details         map[string]any `json:"details"`

    // DO NOT USE DIRECTLY, USE `GetPackType` instead.
    IsQuota  bool                `json:"is_quota"`  // IsQuota marks if ResourcePack is quota pack
    IsAccess bool                `json:"is_access"` // IsAccess marks if ResourcePack is access pack
    Type     ResourcePackageType `json:"pack_type"`

    CropType     CropType      `json:"quota_type"`
    CropTypeMeta *CropTypeMeta `json:"quota_type_meta,omitempty"`

    // IsNotDeductByVolcano marks if this is a resource package deducted by volcano engine;
    // default to false; if true, it means that the usage of this resource pack should
    // be deducted monthly by volcano instead of daily; and its real-time usage is NOT
    // managed by volcano.
    //
    // NOTE: Maybe I named it wrongly in the beginning...
    // `IsDeductMonthlyByResourcePackByVolcano` could be a name more precise.
    IsNotDeductByVolcano bool `json:"not_deduct_by_volc,omitempty"`

    // IsStandalone allows resource pack to be ordered when
    // there's no project or service if set to true; default to false.
    IsStandalone bool `json:"is_standalone,omitempty"`

    // Limits is used when to replace `quota_type` when it is not enough to store only
    // quota type without limit; `IsStandalone` and `IsNotDeductByVolc` packs MUST set this
    // field by nature.
    Limits []HarvestLimit `json:"limits,omitempty"`

    // BundlePacks are bundled packages this package has.
    BundlePacks []BundleResourcePackage `json:"bundle_packs,omitempty"`

    // CodeAttributeTree is a multi-fork tree to carry additional info related to "PreferCode()"
    // of the current resource pack; it can act as an additional filter when selecting resource
    // pack template; it can also carry additional info to downstream tenants.
    CodeAttributeTree *ResourcePackCodeAttributeTree `json:"code_attribute_tree,omitempty"`

    AliasGenerator AliasGenerator `json:"alias_generator,omitempty"`
    SpecifiedName  string         `json:"specified_name"`

    // PersistentCodeAttributes will be copied into `resource_packs.attributes`
    // if no val with the same key exists.
    PersistentCodeAttributes map[string]any `json:"persistent_code_attributes,omitempty"`

    // ForcedExpireQuantity overwrites times number when creating preorder,
    // no matter what the input is and expireQuantity is.
    ForcedExpireQuantity *int `json:"forced_expire_quantity,omitempty"`

    // FormalizedChargeItemCode give_resource_pack can be formalized with this charge item code
    FormalizedChargeItemCode string `json:"formalized_charge_item_code,omitempty"`

    // GivenAfterFirstFormalize means if this resource pack should be given after
    // first formalized in the account operate in the new console.
    GivenAfterFirstFormalize bool `json:"given_after_first_plant,omitempty"`

    // GivenValue means the value of the given prepaid resource pack
    GivenValue int `json:"given_value,omitempty"`

    // GivenAfterDataAuthorization means if this resource pack should be given after data authorization
    GivenAfterDataAuthorization bool `json:"given_after_data_authorization,omitempty"`
}

const (
    QuotaType   ResourcePackageType = "quota"
    PrepaidType ResourcePackageType = "prepaid"
    AccessType  ResourcePackageType = "access"
    DummyType   ResourcePackageType = "dummy"

    // YearsOfNeverExpire overwrite time of order when resource set NeverExpire
    YearsOfNeverExpire = 99
)

type ResourcePackageType string

// ResourcePackageSchemaServicePrepaid marks that this struct is `PrepaidServiceResourcePackage`
// if the `ResourcePackageSchemaType` field is set as the value of this const.
const ResourcePackageSchemaServicePrepaid ResourcePackageSchemaType = "service_resource_pack"

type ResourcePackageSchemaType string

type PrepaidServiceResourcePackage struct {
    ResourcePackage `json:",inline"`

    // ResourcePackageSchemaType must have the value of
    // ResourcePackageSchemaServicePrepaid for PrepaidServiceResourcePackage.
    ResourcePackageSchemaType ResourcePackageSchemaType `json:"resource_package_schema_type"`

    // IsTrial means that if this prepaid service pack will be initialized as a service in trial version.
    IsTrial bool `json:"is_trial"`
}
