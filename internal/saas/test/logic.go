package test

/*

整体数据流向:
		先发布一个没有服务配置的商品版本(创建 commodity_version)
		然后添加服务类型(服务类型的相关配置 => commodity_version_service_type 中的 basic_info)
		最后发布商品(basic_info 配置 => service_type 相关字段)

核心表:
	commodity_version:
		id:
		product_id:
		status: new、launching、launched
		basic_info: 控制台 admin 创建时的 BasicInfo
		release_note: 最终发布时的 ReleaseNote

	commodity_version_service_type:
		service_type_id:
		commodity_version_id:
		basic_info: 控制台 admin 中 AddServiceType 时添加的相关配置
	service_type:
		product_id:
		resource_id:
		trial_limits:
		business_limit_types:
		volcengine_config:

多限流 && 多计费:
	增加服务类型:
		Method: PUT
		path: commodity_management/product/100518/version/1600/service_type/10283
		逻辑:
			1.通过 product_id 和 version_id 获取 commodity_version
			2.通过 commodity_version 中的 commodity_version_id 和 service_type_id 更新
				commodity_version_service_type 表中的 basic_info 字段(限流、资源包等相关配置)

	商品发布:
		Method: POST
		path: commodity_management/product/100520/version/1601/launch
		逻辑:
			1.通过 product_id 获取当前的商品信息; 校验商品状态(不能等于 launching)
			2.构建 commodityCreateReq(CladeCreateReq + SpeciesCreateReqs) 并 Validate
				2.1 通过 product_id 和 version_id 获取 CladeCreateReq
				2.2 通过 product_id 和 version_id 获取 serviceTypeVersions
					2.2.1 通过 product_id 和 version_id 获取商品版本信息(commodity_version)
					2.2.2 通过 commodity_version_id 获取 serviceTypeVersions
				2.3 遍历 serviceTypeVersions 进行数据组装形成 SpeciesCreateReqs
			3. 获取 launchStateMachine 并试运行
				3.1 加载插件
				3.2 加载 product 和 serviceType 的 launchState
			4. 更新商品状态( new => launching )
			5. launchStateMachine 运行
				5.1 遍历所有的 launchState 进行 Launch,如果失败进行回滚、成功则保存之前状态的结果
					5.1.1 SpeciesLaunchState 的 Launch
			6. 更新 commodity_version 中的 release_note
			7. 更新商品状态( launching => launched )
			8. 根据 commodity_version 最后一条记录制作新版本(commodity_version、commodity_version_service_type)

		SpeciesLaunchState 的 Launch:
			1.校验 speciesCreateReq
			2.对之前的商品状态结果进行类型断言
			3.获取老的服务配置表(service_blueprint)及配置拆分表(service_blueprint_translation)
			4.进行 LaunchSpeciesUpdate
				4.1 获取老的服务配置(Species)
				4.2 获取后付费配置项
				4.3 主计费项处理: 获取主计费项服务配置
					4.3.1 对 service_blueprint 表中的字段进行初始化
						收集 volcano_engine_config 字段及其他字段的上下文
					4.3.2 对 service_blueprint_translation 表中的字段进行初始化
						speciesDescCN、speciesDescEN
				4.4 次计费项处理:
					与主计费项相似，通过 bundle_service_blueprint_ids 关联起来


*/

type Item struct {
    Hour  string
    Month string
}

func GroupItemState(slice []string) []Item {
    var items []Item
    var pendingHourly bool

    for _, current := range slice {
        switch current {
        case "hourly":
            if pendingHourly {
                // 如果已有未配对的hourly，先保存它
                items = append(items, Item{Hour: "hourly"})
            }
            pendingHourly = true
        case "monthly":
            if pendingHourly {
                // 配对成功
                items = append(items, Item{Hour: "hourly", Month: "monthly"})
                pendingHourly = false
            } else {
                // 单独的monthly
                items = append(items, Item{Month: "monthly"})
            }
        }
    }

    // 处理最后剩余的hourly
    if pendingHourly {
        items = append(items, Item{Hour: "hourly"})
    }

    return items
}

func GroupItemDoubleIndex(slice []string) []Item {
    var items []Item
    left, right := 0, 0

    for left < len(slice) {
        if slice[left] == "hourly" {
            right = left + 1
            if right < len(slice) && slice[right] == "monthly" {
                items = append(items, Item{Hour: slice[left], Month: slice[right]})
                left = right + 1
            } else {
                items = append(items, Item{Hour: slice[left]})
                left++
            }
        } else if slice[left] == "monthly" {
            items = append(items, Item{Month: slice[left]})
            left++
        } else {
            left++
        }
    }

    return items
}

func GroupItem(slice []string) []Item {
    var items []Item
    i := 0

    for i < len(slice) {
        current := slice[i]

        if current == "hourly" &&
                i+1 < len(slice) && slice[i+1] == "monthly" {
            // hourly后面跟着monthly，配对
            items = append(items, Item{Hour: "hourly", Month: "monthly"})
            i += 2 // 跳过这两个元素
        } else if current == "monthly" {
            // 单独的monthly
            items = append(items, Item{Month: "monthly"})
            i++
        } else if current == "hourly" {
            // 单独的hourly
            items = append(items, Item{Hour: "hourly"})
            i++
        } else {
            i++
        }
    }

    return items
}

type ChargeItem struct {
    Method string // 有三种类型: hour、month、realtime
    Id     int    // 代表计费项编号,hour+month+realtime 组成一个计费项
}

func GroupChargeItemWithPriority(slice []ChargeItem) []ChargeItem {
    result := make([]ChargeItem, len(slice))
    copy(result, slice)
    id := 1
    n := len(slice)
    i := 0

    // automatically determine priority
    priority := "month"
    for _, item := range slice {
        if item.Method == "realtime" {
            priority = "realtime"
            break
        }
    }

    for i < n {
        if priority == "realtime" {
            // 优先匹配长度 3
            if i+2 < n && slice[i].Method == "hour" &&
                    slice[i+1].Method == "month" &&
                    slice[i+2].Method == "realtime" {
                result[i].Id = id
                result[i+1].Id = id
                result[i+2].Id = id

                i += 3
                id++
                continue
            }
        }

        // 匹配长度 2
        if i+1 < n {
            m1 := slice[i].Method
            m2 := slice[i+1].Method
            matched := false

            if priority == "month" {
                if m1 == "hour" && m2 == "month" {
                    result[i].Id = id
                    result[i+1].Id = id
                    i += 2

                    id++
                    matched = true
                }
            } else { // realtime
                if m1 == "hour" && m2 == "month" ||
                        m1 == "hour" && m2 == "realtime" ||
                        m1 == "month" && m2 == "realtime" {
                    result[i].Id = id
                    result[i+1].Id = id
                    i += 2

                    id++
                    matched = true
                }
            }

            if matched {
                continue
            }
        }

        // 单独元素
        result[i].Id = id
        i++
        id++
    }

    return result
}

type postpaidConfigItem struct {
    hourlyPostpaid  ChargeItem
    monthlyPostpaid ChargeItem
}

func GetConfigsTest(slice []ChargeItem) []postpaidConfigItem {
    configs := slice

    if len(configs) <= 0 {
        return []postpaidConfigItem{{}}
    }
    if len(configs) == 1 {
        item := postpaidConfigItem{}

        if configs[0].Method == "hourly" {
            item.hourlyPostpaid = configs[0]
        } else if configs[0].Method == "monthly" {
            item.monthlyPostpaid = configs[0]
        }

        return []postpaidConfigItem{item}
    }

    items := make([]postpaidConfigItem, 0)
    itemMap := make(map[int]postpaidConfigItem)

    for _, config := range configs {
        item, ok := itemMap[config.Id]
        if !ok {
            tmp := postpaidConfigItem{}
            if config.Method == "hourly" {
                tmp.hourlyPostpaid = config
            }
            if config.Method == "monthly" {
                tmp.monthlyPostpaid = config
            }

            itemMap[config.Id] = tmp
        }

        if config.Method == "hourly" {
            item.hourlyPostpaid = config
        }
        if config.Method == "monthly" {
            item.monthlyPostpaid = config
        }
    }

    if len(itemMap) == 1 {
        return []postpaidConfigItem{{
            hourlyPostpaid:  itemMap[1].hourlyPostpaid,
            monthlyPostpaid: itemMap[1].monthlyPostpaid,
        }}
    }

    for i, item := range itemMap {
        if i == 1 {
            items = append(items, item)
        }
    }
    for _, item := range itemMap {
        items = append(items, item)
    }

    return items
}

func GetConfigs(configs []ChargeItem) []postpaidConfigItem {
    itemMap := make(map[int]*postpaidConfigItem)
    for _, c := range configs {
        if _, ok := itemMap[c.Id]; !ok {
            itemMap[c.Id] = &postpaidConfigItem{}
        }
        switch c.Method {
        case "hourly":
            itemMap[c.Id].hourlyPostpaid = c
        case "monthly":
            itemMap[c.Id].monthlyPostpaid = c
        }
    }

    items := make([]postpaidConfigItem, 0, len(itemMap))
    if item, ok := itemMap[1]; ok {
        items = append(items, *item)
    }
    for id, item := range itemMap {
        if id != 1 {
            items = append(items, *item)
        }
    }
    return items
}
