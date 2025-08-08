package demo

type ChargeItem struct {
    Method string
    Id     int
}

func GroupItemIdWithPriority(slice []ChargeItem) []ChargeItem {
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
