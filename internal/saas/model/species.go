package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "strings"
)

/*
CREATE TABLE `service_blueprint` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `psm` varchar(256) DEFAULT NULL COMMENT 'psm this service_blueprint uses',
  `product_id` int(11) NOT NULL COMMENT 'product under which service operates',
  `resource_id` varchar(256) NOT NULL DEFAULT '' COMMENT 'id of resource this blueprint gives access to',
  `billing_product` varchar(256) NOT NULL DEFAULT '' COMMENT 'product in BABI system to issue bill under',
  `billing_service` varchar(256) NOT NULL DEFAULT '' COMMENT 'service in BABI system to issue bill under',
  `billing_service_node` varchar(512) DEFAULT NULL COMMENT 'If this field is not null, babi usage will be pushed to this service tree node ID. This setting overwrites service_tree_node in project table when pushing billing for this service_blueprint.',
  `billing_service_path` varchar(512) DEFAULT NULL COMMENT 'If this field is not null, babi usage will be pushed to this service tree node name. This setting overwrites service_tree_path in project table when pushing billing for this service_blueprint.',
  `billing_location` varchar(256) NOT NULL DEFAULT '' COMMENT 'location in BABI system to issue bill under',
  `billing_price_source` varchar(256) NOT NULL DEFAULT '' COMMENT 'usage that should be sent to BABI system to calculate bill from',
  `billing_price_expr` varchar(1024) NOT NULL DEFAULT '' COMMENT 'expression to calculate price for billing display',
  `billing_username` varchar(128) NOT NULL DEFAULT '' COMMENT 'username to authorize in BABI service',
  `billing_password` varchar(128) NOT NULL DEFAULT '' COMMENT 'username to authorize in BABI service',
  `listing_usage_names` varchar(1024) NOT NULL DEFAULT '' COMMENT 'usages to list in service list view',
  `listing_order` int(11) NOT NULL DEFAULT '0' COMMENT 'position under which to display on product details view',
  `update_limit_workflow` varchar(128) NOT NULL DEFAULT '' COMMENT 'ID of BPM workflow to update limits',
  `trial_limits` varchar(1024) NOT NULL DEFAULT '' COMMENT 'trial limits for new created service',
  `business_limit_types` varchar(1024) NOT NULL DEFAULT '' COMMENT 'limit types for business service',
  `required_appid` varchar(1024) NOT NULL DEFAULT '' COMMENT 'required app ID for services operated by third party apps',
  `listing_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'if service blueprint should be listed in "available_services" in SaaS Platform',
  `kingdom` varchar(256) NOT NULL DEFAULT '' COMMENT 'system that is allowed to use given blueprint',
  `volcengine_config` json DEFAULT NULL COMMENT 'configuration of volcengine product associated with given blueprint',
  `clusters` varchar(256) NOT NULL DEFAULT '' COMMENT 'JSON list of clusters for given resource',
  `saas_config` json DEFAULT NULL COMMENT 'configuration of saas product associated with given blueprint',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10059 DEFAULT CHARSET=utf8 COMMENT='blueprint used to create and access services'
*/

type Species struct {
    ID int `gorm:"column:id;primary_key;AUTO_INCREMENT" db:"id" json:"id"`

    // RequiredAppID should only be filled with third party appid if requests
    // using given resource require a third party appid.
    RequiredAppID string `gorm:"column:required_appid" db:"required_appid" json:"-"`

    PSM *string `gorm:"column:psm" db:"psm" json:"psm"`

    // CladeID points to clade that this species is part of.
    CladeID int `gorm:"column:product_id" db:"product_id" json:"product_id"`

    // ResourceID is used to identify species of requests.
    ResourceID string `gorm:"column:resource_id" db:"resource_id" json:"resource_id"`

    // UpdateLimitWorkflow is the ID of workflow to use when requesting quota limits update.
    UpdateLimitWorkflow string `gorm:"column:update_limit_workflow" db:"update_limit_workflow" json:"-"`

    ListingOrder  int      `gorm:"column:listing_order" db:"listing_order" json:"listing_order"`
    ListingPublic *bool    `gorm:"column:listing_public" db:"listing_public" json:"listing_public,omitempty"`
    Kingdom       string   `gorm:"column:kingdom" db:"kingdom" json:"system"`
    Clusters      Clusters `gorm:"column:clusters" db:"clusters" json:"clusters,omitempty"`

    SpeciesDescription `gorm:"-" db:"service_blueprint_translation" json:",inline"`

    // SpeciesBilling     `db:",inline" json:"billing"`
    BillingProduct     string  `gorm:"column:billing_product" db:"billing_product" json:"billing_product"`
    BillingService     string  `gorm:"column:billing_service" db:"billing_service" json:"billing_service"`
    BillingServiceNode *string `gorm:"column:billing_service_node" db:"billing_service_node" json:"billing_service_node"`
    BillingServicePath *string `gorm:"column:billing_service_path" db:"billing_service_path" json:"billing_service_path"`
    BillingLocation    string  `gorm:"column:billing_location" db:"billing_location" json:"billing_location"`
    BillingPriceSource string  `gorm:"column:billing_price_source" db:"billing_price_source" json:"billing_price_source"`
    BillingPriceExpr   string  `gorm:"column:billing_price_expr" db:"billing_price_expr" json:"billing_price_expr"`
    BillingUsername    string  `gorm:"column:billing_username" db:"billing_username" json:"billing_username"`
    BillingPassword    string  `gorm:"column:billing_password" db:"billing_password" json:"billing_password"`

    // SpeciesHarvesting `db:",inline" json:",inline"`

    TrialHarvestLimits HarvestLimitations `gorm:"column:trial_limits" db:"trial_limits" json:"quota_trial"`
    BusinessCropTypes  BusinessCropTypes  `gorm:"column:business_limit_types" db:"business_limit_types" json:"quota_business"`

    ListingUsageCropTypes CropTypes                   `gorm:"column:listing_usage_names" db:"listing_usage_names" json:"-"`
    VolcanoEngineConfig   *SpeciesVolcanoEngineConfig `gorm:"column:volcengine_config" db:"volcengine_config" json:"volcengine_config,omitempty"`
    SaaSConfig            *SpeciesSaaSConfig          `gorm:"column:saas_config" db:"saas_config" json:"saas_config,omitempty"`
}

/*
CREATE TABLE `service_blueprint_translation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `service_blueprint_id` int(11) DEFAULT NULL COMMENT 'primary key column',
  `lang` varchar(16) DEFAULT NULL COMMENT 'primary key column',
  `listing_group` text COMMENT 'name of group under which service blueprint is displayed',
  `details_name` text COMMENT 'name of service blueprint',
  `details_description` text COMMENT 'description of service blueprint',
  `details_price_trial` text COMMENT 'information about trial price',
  `details_price_full` text COMMENT 'information about business price',
  PRIMARY KEY (`id`),
  KEY `idx_blueprint_id` (`service_blueprint_id`)
) ENGINE=InnoDB AUTO_INCREMENT=345 DEFAULT CHARSET=utf8 COMMENT='service translation for multilanguage support'
*/

type SpeciesDescription struct {
    SpeciesID   int    `db:"service_blueprint_id" gorm:"column:service_blueprint_id" json:"-"`
    Name        string `db:"details_name" gorm:"column:details_name" json:"details_name"`
    Language    string `db:"lang" gorm:"column:lang" json:"-"`
    Genus       string `db:"listing_group" gorm:"column:listing_group" json:"listing_group"`
    Description string `db:"details_description" gorm:"column:details_description" json:"details_description"`
    PriceTrial  string `db:"details_price_trial" gorm:"column:details_price_trial" json:"details_price_trial"`
    PriceFull   string `db:"details_price_full" gorm:"column:details_price_full"  json:"details_price_full"`
}

type Clusters []string

func (c *Clusters) Scan(src any) error {
    if src == nil {
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        panic(fmt.Sprintf("invalid Clusters.Scan source type: %T", src))
    }

    if len(srcBytes) == 0 {
        return nil
    }

    return json.Unmarshal(srcBytes, c)
}

func (c Clusters) Value() (driver.Value, error) {
    if len(c) == 0 {
        return "", nil
    }

    b, err := json.Marshal(c)
    if err != nil {
        return nil, fmt.Errorf("marshal Clusters: %v", err)
    }

    return string(b), nil
}

type BusinessCropTypes []BusinessCropType

type BusinessCropType struct {
    CropType     CropType      `json:"name"`
    CropTypeMeta *CropTypeMeta `json:"crop_type_with_meta,omitempty"`
}

type CropTypes []CropType

func (ct *CropTypes) Scan(src any) error {
    srcBytes, ok := src.([]byte)
    if !ok {
        return fmt.Errorf("invalid CropTypes.Scan type: %T", src)
    }
    srcString := string(srcBytes)
    if srcString == "" {
        return nil
    }

    stringSlice := strings.Split(srcString, ",")
    *ct = make(CropTypes, len(stringSlice))
    for i, name := range stringSlice {
        (*ct)[i] = CropType(name)
    }

    return nil
}

func (cd CropTypes) Value() (driver.Value, error) {
    cropTypes := make([]string, len(cd))
    for i, cropType := range cd {
        cropTypes[i] = string(cropType)
    }
    return []byte(strings.Join(cropTypes, ",")), nil
}

func (cd CropTypes) Contain(cropType CropType) bool {
    for _, singleCropType := range cd {
        if singleCropType == cropType {
            return true
        }
    }

    return false
}

func (ct CropTypes) Copy() CropTypes {
    var copied CropTypes
    copied = append(copied, ct...)
    return copied
}

const (
    // VolcanoEnginePricingDisplayDefault means to use info
    // inside `VolcanoEnginePricing` as the source of pricing display.
    VolcanoEnginePricingDisplayDefault VolcanoEnginePricingDisplayMethod = ""

    // VolcanoEnginePricingDisplayPrepaid means to use info
    // inside `PrepaidResourcePacks` as the source of pricing display.
    VolcanoEnginePricingDisplayPrepaid   VolcanoEnginePricingDisplayMethod = "prepaid_service_resource_pack"
    ResourcePackagePricingDisplayDefault VolcanoEnginePricingDisplayMethod = ""
    ResourcePackagePricingDisplayVolcano VolcanoEnginePricingDisplayMethod = "volc"
)

type (
    ResourcePackagePricingDisplayMethod     string
    PrepaidServiceResourcePackDisplayFilter string
    VolcanoEnginePricingDisplayMethod       string
)

type VolcanoEngineAccessCategory struct {
    Category      string                      `json:"category"`
    AccessMethods []VolcanoEngineAccessMethod `json:"access_methods"`
}

type VolcanoEngineAccessMethod struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

type VolcanoEnginePricing struct {
    APIHeader                  string `json:"api_header"`
    APIName                    string `json:"api_name"`
    GradientHeader             string `json:"gradient_header"`     // Sometimes we just regard [fix price] as a special case of gradient case
    GradientValueUnit          string `json:"gradient_value_unit"` // if the case is [fix price] it just means it's unit
    GradientBoundaryExpression string `json:"gradient_boundary_expression"`
    PriceHeader                string `json:"price_header"`
    ForceChargeItemCode        string `json:"force_charge_item_code,omitempty"` // 计费项
}

type SpeciesVolcanoEngineConfig struct {
    Product                     string `json:"product"`
    ConfigurationCode           string `json:"configuration_code"`
    FormalChargeItemCode        string `json:"formal_charge_item_code"`
    MonthlyConfigurationCode    string `json:"monthly_configuration_code"`
    MonthlyFormalChargeItemCode string `json:"monthly_formal_charge_item_code"`
    BillingRulesURL             string `json:"billing_rules"`
    DefaultBusinessQuota        int    `json:"default_business_quota"`

    // LaunchAgainForbidden means that if this species(serviceType) should forbid launch again.
    LaunchAgainForbidden bool `json:"launch_again_forbidden,omitempty"`

    VolcanoEnginePricing VolcanoEnginePricing `json:"pricing_display_config"`
    PrepaidResourcePacks []ResourcePackage    `json:"prepaid_resource_packs"`
    QuotaResourcePacks   []ResourcePackage    `json:"quota_resource_packs"`
    AccessResourcePacks  []ResourcePackage    `json:"access_resource_packs"`

    // DummyResourcePackages are packs hidden to client and plays no role in affecting
    // service quota or status; it could be part of a resource pack bundle that client
    // is forced to buy; its resource packs all have their type as "Dummy".
    DummyResourcePackages []ResourcePackage `json:"dummy_resource_packs"`

    // OfflineResourcePackages should record all resource packs taken offline so that
    // users who bought them can still access their config & display info.
    // Actually stored is the Resource Pack that has been offline, not only `access` type
    OfflineResourcePackages []ResourcePackage `json:"offline_access_resource_packs"`

    AccessDocumentation []VolcanoEngineAccessCategory `json:"access_documentation"`

    // MonitoringQuotaType
    // NOTE: MonitoringQuotaType is recommended not to be used isolated,
    // Use GetMonitoringQuotaTypeWithMeta() instead
    MonitoringQuotaType     CropType      `json:"monitoring_quota_type"`
    MonitoringQuotaTypeMeta *CropTypeMeta `json:"monitoring_quota_type_meta,omitempty"`

    // v3模型引入多个限流条件
    BundleMonitoringQuotaTypes CropTypes `json:"bundle_monitoring_quota_types,omitempty"`
    BundleDefaultBusinessQuota []int     `json:"bundle_default_business_quota,omitempty"`

    // MonitoringUsageType
    // NOTE: MonitoringUsageType is recommended not to be used isolated,
    // Use GetMonitoringUsageTypeWithMeta() instead
    MonitoringUsageType     CropType      `json:"monitoring_usage_type"`
    MonitoringUsageTypeMeta *CropTypeMeta `json:"monitoring_usage_type_with_meta,omitempty"`

    // Hidden means that if this species(serviceType) should display or hidden in volcano-console.
    Hidden bool `json:"hidden"`

    // TrialDurationDay means how long should a service's trial expire after it's first creation
    // it is something maintain in volcano engine side but not concept that saas consume.
    TrialDurationDay int `json:"trial_duration_day"`

    // VolcanoEnginePricingDisplayMethod decides how and from where
    // to get pricing display for plant / service.
    PricingDisplayMethod VolcanoEnginePricingDisplayMethod `json:"pricing_display_method"`

    // ResourcePackPricingDisplayMethod decides how and from where
    // to get pricing display for resource packs.
    ResourcePackPricingDisplayMethod ResourcePackagePricingDisplayMethod `json:"resource_pack_pricing_display_method"`

    // PrepaidServiceResourcePackDisplayFilters is a group of filters pre-defined to
    // apply to prepaid service resource packs before displaying their price.
    PrepaidServiceResourcePackDisplayFilters []PrepaidServiceResourcePackDisplayFilter `json:"prepaid_service_resource_pack_display_filters"`

    // PrepaidServiceResourcePacks marks resource packs that are born with prepaid services;
    // since it contains extra info.
    // if is for trial version, that won't be useful for normal resource_pack,
    // it is designed as a superset of it.
    PrepaidServiceResourcePacks []PrepaidServiceResourcePackage `json:"prepaid_service_resource_packs"`

    // BindingPlatform
    // Sometimes the project is not necessary when creating a service,
    // use the BindingPlatform to associate the name of the platform when project be created
    // such as: Avatar
    BindingPlatform ParcelPlatform `json:"binding_platform"`

    // GivenResourcePackages are packs given freely to services
    // whenever client requires with `is_given` marked as true.
    GivenResourcePackages []ResourcePackage `json:"given_resource_packs,omitempty"`

    GivenStandaloneSpeciesIDs []int `json:"given_standalone_species_ids,omitempty"`

    // IsTrialSkipped indicates if service should skip trial rank
    // and formalized into business directly.
    IsTrialSkipped bool `json:"is_trial_skipped,omitempty"`

    // IsBillingSkipped indicates if we should send any usage bill to volcano.
    IsBillingSkipped bool `json:"is_billing_skipped,omitempty"`

    // IsGrayStage 商品是否为灰度发布状态,灰度商品需要判断账号权限,
    // 通过 GenerateChildrenFromProductIDWithMap 控制前端渲染结果标签,如果没有权限是无法在前端展示使用入口和后端下单的
    // 灰度的商品依赖 Account 上打的标签,根据 Account 上的标签确认是否可以购买,
    // https://bytedance.larkoffice.com/docx/OcaRdjeIyowMecxfesycmgbPnwe#part-VjjTd281hoZ8kXxGf5ycC6m1nif
    IsGrayStage bool `json:"is_gray_stage,omitempty"`

    // GrayScope 通过计算对应账户的身份值 score,判断身份值在不在 GrayScope 中,如果在里面,则允许访问,反之亦然
    GrayScope []int `json:"gray_scope,omitempty"`

    // 引入一个火山实例多计费项新增的字段 =====> start

    // 关联计费项，这样设计的好处是:每一个 service_blueprint_id 的配置结构都不需要大改,只需要使用时候合并起来即可
    BundleServiceBlueprintIDs    []int    `json:"bundle_service_blueprint_ids,omitempty"`
    FormalChargeItemCodes        []string `json:"formal_charge_item_codes,omitempty"` // 火山官网计费项没有顺序，这里存放所有的计费项
    MonthlyFormalChargeItemCodes []string `json:"monthly_formal_charge_item_codes,omitempty"`
    // 引入一个火山实例多计费项新增的字段 =====> end
}

type SpeciesSaaSConfig struct {
    MonitoringQuotaType CropType `json:"monitoring_quota_type"`
    MonitoringUsageType CropType `json:"monitoring_usage_type"`

    // BusinessQuotaLimit BusinessUsageLimit
    // just a recommended usage, be filled in when approval business.
    BusinessQuotaLimit int64 `json:"business_quota_limit,omitempty"`
    BusinessUsageLimit int64 `json:"business_usage_limit,omitempty"`
}
