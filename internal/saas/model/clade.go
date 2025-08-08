package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

/*
CREATE TABLE `product` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `listing_usage_names` varchar(1024) NOT NULL DEFAULT '' COMMENT 'usages that are listed on product list page',
  `listing_order` int(11) NOT NULL DEFAULT '0' COMMENT 'position on which product will be displayed in group',
  `details_icon` mediumtext COMMENT 'product icon',
  `kingdom` varchar(256) NOT NULL DEFAULT '' COMMENT 'system for which product is designed',
  `business_line` varchar(256) NOT NULL DEFAULT '' COMMENT 'business line the product belongs to',
  `product_category` varchar(256) NOT NULL DEFAULT '' COMMENT 'product category',
  `api_service` varchar(1024) NOT NULL DEFAULT '[]' COMMENT 'list of api endpoint urls corresponding to the product',
  `service_tree_node` varchar(512) NOT NULL DEFAULT '' COMMENT 'service tree node used to push the cost of commodity',
  `service_tree_path` varchar(512) NOT NULL DEFAULT '' COMMENT 'service tree path used to push the cost of commodity',
  `psm` varchar(256) NOT NULL DEFAULT '' COMMENT 'product psm',
  `volc_product` varchar(256) NOT NULL DEFAULT '' COMMENT 'english name of commodity in volcano-cost-center, used when sending an order to volc',
  `product_doc` varchar(1024) NOT NULL DEFAULT '' COMMENT 'product doc for volcano-console display',
  `pricing_doc` varchar(1024) NOT NULL DEFAULT '' COMMENT 'pricing doc for volcano-console display',
  `api_doc` varchar(1024) NOT NULL DEFAULT '' COMMENT 'api integration doc for volcano-console display',
  `sdk_doc` varchar(1024) NOT NULL DEFAULT '' COMMENT 'sdk integration doc for volcano-console display',
  `launch_status` varchar(32) NOT NULL DEFAULT '' COMMENT 'enum within: new,launching,launched',
  `delisted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'delisted item should not be shown in user side',
  `listing_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'if service blueprint should be listed in SaaS Platform',
  `is_big_model` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'is big model product',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10034 DEFAULT CHARSET=utf8 COMMENT='speech team product'
*/

const (
    CladeLaunchStatusNew       CladeLaunchStatus = "new"
    CladeLaunchStatusLaunching CladeLaunchStatus = "launching"
    CladeLaunchStatusLaunched  CladeLaunchStatus = "launched"
)

type Clade struct {
    ID              int               `gorm:"column:id;primary_key;AUTO_INCREMENT" db:"id" json:"id"`
    Icon            string            `gorm:"column:details_icon" db:"details_icon" json:"details_icon,omitempty"`
    ListingOrder    int               `gorm:"column:listing_order" db:"listing_order" json:"listing_order"`
    Kingdom         string            `gorm:"column:kingdom" db:"kingdom" json:"system"`
    BusinessLine    string            `gorm:"column:business_line" db:"business_line" json:"business_line"`
    ProductCategory string            `gorm:"column:product_category" db:"product_category" json:"product_category"`
    ApiServices     ApiServiceUrls    `gorm:"column:api_service" db:"api_service" json:"api_service"`
    ServiceTreeNode string            `gorm:"column:service_tree_node" db:"service_tree_node" json:"service_tree_node"`
    ServiceTreePath string            `gorm:"column:service_tree_path" db:"service_tree_path" json:"service_tree_path"`
    PSM             string            `gorm:"column:psm" db:"psm" json:"psm"`
    VolcanoProduct  string            `gorm:"column:volc_product" db:"volc_product" json:"volc_product"`
    ProductDoc      string            `gorm:"column:product_doc" db:"product_doc" json:"product_doc"`
    PricingDoc      string            `gorm:"column:pricing_doc" db:"pricing_doc" json:"pricing_doc"`
    APIDoc          string            `gorm:"column:api_doc" db:"api_doc" json:"api_doc"`
    SDKDoc          string            `gorm:"column:sdk_doc" db:"sdk_doc" json:"sdk_doc"`
    LaunchStatus    CladeLaunchStatus `gorm:"column:launch_status" db:"launch_status" json:"launch_status"`

    // make Delisted a pointer to make it compatible with KORM
    Delisted *bool `gorm:"column:delisted" db:"delisted" json:"delisted"`

    // ListingPublic for saas console, decide whether to show it on saas console
    ListingPublic *bool `gorm:"column:listing_public" db:"listing_public" json:"listing_public,omitempty"`

    IsBigModel       bool `gorm:"column:is_big_model" db:"is_big_model" json:"is_big_model"`
    CladeDescription `gorm:"-" db:"product_translation" json:",inline"`
}

/*
CREATE TABLE `product_translation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `product_id` int(11) NOT NULL COMMENT 'ID of product this translation applies to',
  `lang` varchar(16) NOT NULL COMMENT 'language shortcode is the language of description for given clade',
  `listing_group` text COMMENT 'name of group under which product is displayed',
  `details_name` text COMMENT 'name of product',
  `details_description_short` text COMMENT 'description for listing view',
  `details_description_full` text COMMENT 'description for details view',
  `details_name_short` text COMMENT 'short name of product',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=137 DEFAULT CHARSET=utf8 COMMENT='product translation for multilanguage support'
*/

type CladeDescription struct {
    CladeID          int     `gorm:"column:product_id;primary_key;AUTO_INCREMENT" db:"product_id" json:"product_id"`
    Language         string  `gorm:"column:lang" db:"lang" json:"-"`
    Name             string  `gorm:"column:details_name" db:"details_name" json:"details_name"`
    ShortName        *string `gorm:"column:details_name_short" db:"details_name_short" json:"details_name_short"`
    ListingGroup     string  `gorm:"column:listing_group" db:"listing_group" json:"listing_group"`
    DescriptionShort string  `gorm:"column:details_description_short" db:"details_description_short" json:"details_description_short,omitempty"`
    DescriptionFull  string  `gorm:"column:details_description_full" db:"details_description_full" json:"details_description_full,omitempty"`
}

/*
CREATE TABLE `project2product` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key column',
  `project_id` int(11) DEFAULT NULL COMMENT 'ID of project',
  `product_id` int(11) DEFAULT NULL COMMENT 'ID of activated product',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7369 DEFAULT CHARSET=utf8 COMMENT='marks enabled products for project'
*/

type CladeHabitat struct {
    ID       int `db:"id"`
    ParcelID int `db:"project_id"`
    CladeID  int `db:"product_id"`
}

type CladeLaunchStatus string

func (c *CladeLaunchStatus) Scan(src any) error {
    if src == nil {
        *c = ""
        return nil
    }

    switch value := src.(type) {
    case []byte:
        *c = CladeLaunchStatus(value)
    case string:
        *c = CladeLaunchStatus(value)
    default:
        return fmt.Errorf("invalid Kingdom.Scan source type: %T", src)
    }

    return nil
}

func (c CladeLaunchStatus) Value() (driver.Value, error) { return string(c), nil }

type ApiServiceUrls []string

func (c *ApiServiceUrls) Scan(src any) error {
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

func (c ApiServiceUrls) Value() (driver.Value, error) {
    if len(c) == 0 {
        return "", nil
    }

    b, err := json.Marshal(c)
    if err != nil {
        return nil, fmt.Errorf("marshal Clusters: %v", err)
    }

    return string(b), nil
}
