package demo

import (
    "fmt"
    "reflect"
    "testing"
)

// 8 个测试用例，覆盖优先级为 hour+month 优先组合
func TestGroupItemId_PriorityMonth_All8(t *testing.T) {
    tests := []struct {
        name  string
        input []ChargeItem
        want  []ChargeItem
    }{
        {"hour+hour+hour",
            []ChargeItem{{"hour", 0}, {"hour", 0}, {"hour", 0}},
            []ChargeItem{{"hour", 1}, {"hour", 2}, {"hour", 3}}},
        {"hour+hour+month",
            []ChargeItem{{"hour", 0}, {"hour", 0}, {"month", 0}},
            []ChargeItem{{"hour", 1}, {"hour", 2}, {"month", 2}}},
        {"hour+month+hour",
            []ChargeItem{{"hour", 0}, {"month", 0}, {"hour", 0}},
            []ChargeItem{{"hour", 1}, {"month", 1}, {"hour", 2}}},
        {"hour+month+month",
            []ChargeItem{{"hour", 0}, {"month", 0}, {"month", 0}},
            []ChargeItem{{"hour", 1}, {"month", 1}, {"month", 2}}},

        {"month+hour+hour",
            []ChargeItem{{"month", 0}, {"hour", 0}, {"hour", 0}},
            []ChargeItem{{"month", 1}, {"hour", 2}, {"hour", 3}}},
        {"month+hour+month",
            []ChargeItem{{"month", 0}, {"hour", 0}, {"month", 0}},
            []ChargeItem{{"month", 1}, {"hour", 2}, {"month", 2}}},
        {"month+month+hour",
            []ChargeItem{{"month", 0}, {"month", 0}, {"hour", 0}},
            []ChargeItem{{"month", 1}, {"month", 2}, {"hour", 3}}},
        {"month+month+month",
            []ChargeItem{{"month", 0}, {"month", 0}, {"month", 0}},
            []ChargeItem{{"month", 1}, {"month", 2}, {"month", 3}}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GroupItemIdWithPriority(tt.input); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GroupItemId() = %v, want %v", got, tt.want)
            }
        })
    }
}

// 39 个测试用例，覆盖优先级为 hour+month+realtime 优先组合
func TestGroupItemId_PriorityRealtime_All39_T(t *testing.T) {
    type args struct {
        slice []ChargeItem
    }
    tests := []struct {
        name string
        args args
        want []ChargeItem
    }{}

    methods := []string{"hour", "month", "realtime"}
    seqIdx := 1

    // 生成长度 1
    for _, m1 := range methods {
        slice := []ChargeItem{{Method: m1}}
        want := []ChargeItem{{Method: m1, Id: 1}}
        tests = append(tests, struct {
            name string
            args args
            want []ChargeItem
        }{
            name: fmt.Sprintf("seq%d_len1_%s", seqIdx, m1),
            args: args{slice},
            want: want,
        })
        seqIdx++
    }

    // 生成长度 2
    for _, m1 := range methods {
        for _, m2 := range methods {
            slice := []ChargeItem{{Method: m1}, {Method: m2}}
            want := make([]ChargeItem, 2)
            copy(want, slice)
            currId := 1
            i := 0
            for i < 2 {
                if i+1 < 2 {
                    m1c := slice[i].Method
                    m2c := slice[i+1].Method
                    switch {
                    case m1c == "hour" && m2c == "month":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    case m1c == "month" && m2c == "realtime":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    case m1c == "hour" && m2c == "realtime":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    default:
                        want[i].Id = currId
                        i++
                    }
                } else {
                    want[i].Id = currId
                    i++
                }
                currId++
            }
            tests = append(tests, struct {
                name string
                args args
                want []ChargeItem
            }{
                name: fmt.Sprintf("seq%d_len2_%s_%s", seqIdx, m1, m2),
                args: args{slice},
                want: want,
            })
            seqIdx++
        }
    }

    // 生成长度 3
    for _, m1 := range methods {
        for _, m2 := range methods {
            for _, m3 := range methods {
                slice := []ChargeItem{{Method: m1}, {Method: m2}, {Method: m3}}
                want := make([]ChargeItem, 3)
                copy(want, slice)
                currId := 1
                i := 0
                for i < 3 {
                    if i+2 < 3 &&
                            slice[i].Method == "hour" &&
                            slice[i+1].Method == "month" &&
                            slice[i+2].Method == "realtime" {
                        want[i].Id = currId
                        want[i+1].Id = currId
                        want[i+2].Id = currId
                        i += 3
                    } else if i+1 < 3 {
                        m1c := slice[i].Method
                        m2c := slice[i+1].Method
                        switch {
                        case m1c == "hour" && m2c == "month":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        case m1c == "month" && m2c == "realtime":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        case m1c == "hour" && m2c == "realtime":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        default:
                            want[i].Id = currId
                            i++
                        }
                    } else {
                        want[i].Id = currId
                        i++
                    }
                    currId++
                }
                tests = append(tests, struct {
                    name string
                    args args
                    want []ChargeItem
                }{
                    name: fmt.Sprintf("seq%d_len3_%s_%s_%s", seqIdx, m1, m2, m3),
                    args: args{slice},
                    want: want,
                })
                seqIdx++
            }
        }
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GroupItemIdWithPriority(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GroupItemId() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGroupItemId_PriorityRealtime_All39(t *testing.T) {
    tests := []struct {
        name  string
        input []ChargeItem
        want  []ChargeItem
    }{
        // ---------- length 1 (1条) ----------
        {"len1_realtime",
            []ChargeItem{{"realtime", 0}},
            []ChargeItem{{"realtime", 1}}},

        // ---------- length 2 (6条) ----------
        // 13 23 31 32 12 21
        {"len2_hour_realtime",
            []ChargeItem{{"hour", 0}, {"realtime", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}}},
        {"len2_month_realtime",
            []ChargeItem{{"month", 0}, {"realtime", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}}},
        {"len2_realtime_hour",
            []ChargeItem{{"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}}},
        {"len2_realtime_month",
            []ChargeItem{{"realtime", 0}, {"month", 0}},
            []ChargeItem{{"realtime", 1}, {"month", 2}}},
        {"len2_hour_month",
            []ChargeItem{{"hour", 0}, {"month", 0}},
            []ChargeItem{{"hour", 1}, {"month", 1}}},
        {"len2_month_hour",
            []ChargeItem{{"month", 0}, {"hour", 0}},
            []ChargeItem{{"month", 1}, {"hour", 2}}},

        // ---------- length 3 (32条) ----------
        // 113 123 131 132 133 213
        {"len3_hour_hour_realtime",
            []ChargeItem{{"hour", 0}, {"hour", 0}, {"realtime", 0}},
            []ChargeItem{{"hour", 1}, {"hour", 2}, {"realtime", 2}}},
        {"len3_hour_month_realtime",
            []ChargeItem{{"hour", 0}, {"month", 0}, {"realtime", 0}},
            []ChargeItem{{"hour", 1}, {"month", 1}, {"realtime", 1}}},
        {"len3_hour_realtime_hour",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"hour", 2}}},
        {"len3_hour_realtime_month",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"month", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"month", 2}}},
        {"len3_hour_realtime_realtime",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"realtime", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"realtime", 2}}},
        {"len3_month_hour_realtime",
            []ChargeItem{{"month", 0}, {"hour", 0}, {"realtime", 0}},
            []ChargeItem{{"month", 1}, {"hour", 2}, {"realtime", 2}}},

        // 223 231 232 233 311 312
        {"len3_month_month_realtime",
            []ChargeItem{{"month", 0}, {"month", 0}, {"realtime", 0}},
            []ChargeItem{{"month", 1}, {"month", 2}, {"realtime", 2}}},
        {"len3_month_realtime_hour",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"hour", 2}}},
        {"len3_month_realtime_month",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"month", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"month", 2}}},
        {"len3_month_realtime_realtime",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"realtime", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"realtime", 2}}},
        {"len3_realtime_hour_hour",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"hour", 3}}},
        {"len3_realtime_hour_month",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"month", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"month", 2}}},

        // 313 321 322 323 331 332 333
        {"len3_realtime_hour_realtime",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"realtime", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"realtime", 2}}},
        {"len3_realtime_month_hour",
            []ChargeItem{{"realtime", 0}, {"month", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"month", 2}, {"hour", 3}}},
        {"len3_realtime_month_month",
            []ChargeItem{{"realtime", 0}, {"month", 0}, {"month", 0}},
            []ChargeItem{{"realtime", 1}, {"month", 2}, {"month", 3}}},
        {"len3_realtime_month_realtime",
            []ChargeItem{{"realtime", 0}, {"month", 0}, {"realtime", 0}},
            []ChargeItem{{"realtime", 1}, {"month", 2}, {"realtime", 2}}},
        {"len3_realtime_realtime_hour",
            []ChargeItem{{"realtime", 0}, {"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"realtime", 2}, {"hour", 3}}},
        {"len3_realtime_realtime_month",
            []ChargeItem{{"realtime", 0}, {"realtime", 0}, {"month", 0}},
            []ChargeItem{{"realtime", 1}, {"realtime", 2}, {"month", 3}}},
        {"len3_realtime_realtime_realtime",
            []ChargeItem{{"realtime", 0}, {"realtime", 0}, {"realtime", 0}},
            []ChargeItem{{"realtime", 1}, {"realtime", 2}, {"realtime", 3}}},

        // 补充长度 3 其余组合至总数 39 条
        // 131 132 133 231 232 233
        {"len3_hour_realtime_hour2",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"hour", 2}}},
        {"len3_hour_realtime_month2",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"month", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"month", 2}}},
        {"len3_hour_realtime_realtime2",
            []ChargeItem{{"hour", 0}, {"realtime", 0}, {"realtime", 0}},
            []ChargeItem{{"hour", 1}, {"realtime", 1}, {"realtime", 2}}},
        {"len3_month_realtime_hour2",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"hour", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"hour", 2}}},
        {"len3_month_realtime_month2",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"month", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"month", 2}}},
        {"len3_month_realtime_realtime2",
            []ChargeItem{{"month", 0}, {"realtime", 0}, {"realtime", 0}},
            []ChargeItem{{"month", 1}, {"realtime", 1}, {"realtime", 2}}},

        // 311 312 313 321
        {"len3_realtime_hour_hour2",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"hour", 3}}},
        {"len3_realtime_hour_month2",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"month", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"month", 2}}},
        {"len3_realtime_hour_realtime2",
            []ChargeItem{{"realtime", 0}, {"hour", 0}, {"realtime", 0}},
            []ChargeItem{{"realtime", 1}, {"hour", 2}, {"realtime", 2}}},
        {"len3_realtime_month_hour2",
            []ChargeItem{{"realtime", 0}, {"month", 0}, {"hour", 0}},
            []ChargeItem{{"realtime", 1}, {"month", 2}, {"hour", 3}}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GroupItemIdWithPriority(tt.input); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GroupItemIdWithPriority() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGroupItemId_PriorityRealtime_All39_T1(t *testing.T) {
    tests := []struct {
        name  string
        input []ChargeItem
        want  []ChargeItem
    }{}

    methods := []string{"hour", "month", "realtime"}
    seqIdx := 1

    // ---------- length 1 (3种) ----------
    for _, m1 := range methods {
        slice := []ChargeItem{{Method: m1}}
        want := []ChargeItem{{Method: m1, Id: 1}}
        tests = append(tests, struct {
            name  string
            input []ChargeItem
            want  []ChargeItem
        }{
            name:  fmt.Sprintf("len1_seq%d_%s", seqIdx, m1),
            input: slice,
            want:  want,
        })
        seqIdx++
    }

    // ---------- length 2 (9种) ----------
    for _, m1 := range methods {
        for _, m2 := range methods {
            slice := []ChargeItem{{Method: m1}, {Method: m2}}
            want := make([]ChargeItem, 2)
            copy(want, slice)
            currId := 1
            i := 0
            for i < 2 {
                if i+1 < 2 {
                    m1c := slice[i].Method
                    m2c := slice[i+1].Method
                    switch {
                    case m1c == "hour" && m2c == "month":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    case m1c == "hour" && m2c == "realtime":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    case m1c == "month" && m2c == "realtime":
                        want[i].Id = currId
                        want[i+1].Id = currId
                        i += 2
                    default:
                        want[i].Id = currId
                        i++
                    }
                } else {
                    want[i].Id = currId
                    i++
                }
                currId++
            }
            tests = append(tests, struct {
                name  string
                input []ChargeItem
                want  []ChargeItem
            }{
                name:  fmt.Sprintf("len2_seq%d_%s_%s", seqIdx, m1, m2),
                input: slice,
                want:  want,
            })
            seqIdx++
        }
    }

    // ---------- length 3 (27种) ----------
    for _, m1 := range methods {
        for _, m2 := range methods {
            for _, m3 := range methods {
                slice := []ChargeItem{{Method: m1}, {Method: m2}, {Method: m3}}
                want := make([]ChargeItem, 3)
                copy(want, slice)
                currId := 1
                i := 0
                for i < 3 {
                    if i+2 < 3 &&
                            slice[i].Method == "hour" &&
                            slice[i+1].Method == "month" &&
                            slice[i+2].Method == "realtime" {
                        want[i].Id = currId
                        want[i+1].Id = currId
                        want[i+2].Id = currId
                        i += 3
                    } else if i+1 < 3 {
                        m1c := slice[i].Method
                        m2c := slice[i+1].Method
                        switch {
                        case m1c == "hour" && m2c == "month":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        case m1c == "hour" && m2c == "realtime":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        case m1c == "month" && m2c == "realtime":
                            want[i].Id = currId
                            want[i+1].Id = currId
                            i += 2
                        default:
                            want[i].Id = currId
                            i++
                        }
                    } else {
                        want[i].Id = currId
                        i++
                    }
                    currId++
                }
                tests = append(tests, struct {
                    name  string
                    input []ChargeItem
                    want  []ChargeItem
                }{
                    name:  fmt.Sprintf("len3_seq%d_%s_%s_%s", seqIdx, m1, m2, m3),
                    input: slice,
                    want:  want,
                })
                seqIdx++
            }
        }
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GroupItemIdWithPriority(tt.input); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GroupItemIdWithPriority() = %v, want %v", got, tt.want)
            }
        })
    }
}

// 53 个测试用例，覆盖所有组合
func TestGroupItemId_AllPriorities(t *testing.T) {
    type args struct {
        slice []ChargeItem
    }

    tests := []struct {
        name string
        args args
        want []ChargeItem
    }{}

    seqIdx := 1

    // --------------------
    // Priority = month: 只包含 hour+month 的长度 3 全组合（8种）
    // --------------------
    monthSeqs := [][]string{
        {"hour", "hour", "hour"},
        {"hour", "hour", "month"},
        {"hour", "month", "hour"},
        {"hour", "month", "month"},
        {"month", "hour", "hour"},
        {"month", "hour", "month"},
        {"month", "month", "hour"},
        {"month", "month", "month"},
    }

    for _, seq := range monthSeqs {
        slice := make([]ChargeItem, len(seq))
        copyWant := make([]ChargeItem, len(seq))
        for i, m := range seq {
            slice[i] = ChargeItem{Method: m}
            copyWant[i] = ChargeItem{Method: m}
        }

        // 生成 want
        id := 1
        i := 0
        n := len(slice)
        for i < n {
            if i+1 < n && slice[i].Method == "hour" && slice[i+1].Method == "month" {
                copyWant[i].Id = id
                copyWant[i+1].Id = id
                i += 2
                id++
            } else {
                copyWant[i].Id = id
                i++
                id++
            }
        }

        tests = append(tests, struct {
            name string
            args args
            want []ChargeItem
        }{
            name: fmt.Sprintf("month_seq%d", seqIdx),
            args: args{slice},
            want: copyWant,
        })
        seqIdx++
    }

    // --------------------
    // Priority = realtime: 长度 1~3 所有组合（39种）
    // --------------------
    methods := []string{"hour", "month", "realtime"}
    for l := 1; l <= 3; l++ {
        var generate func(pos int, cur []string)
        generate = func(pos int, cur []string) {
            if pos == l {
                slice := make([]ChargeItem, l)
                copyWant := make([]ChargeItem, l)
                priority := "month"
                for i, m := range cur {
                    slice[i] = ChargeItem{Method: m}
                    copyWant[i] = ChargeItem{Method: m}
                    if m == "realtime" {
                        priority = "realtime"
                    }
                }

                // 生成 want
                id := 1
                i := 0
                n := l
                for i < n {
                    if priority == "realtime" && i+2 < n &&
                            slice[i].Method == "hour" &&
                            slice[i+1].Method == "month" &&
                            slice[i+2].Method == "realtime" {
                        copyWant[i].Id = id
                        copyWant[i+1].Id = id
                        copyWant[i+2].Id = id
                        i += 3
                        id++
                        continue
                    }
                    if i+1 < n {
                        m1 := slice[i].Method
                        m2 := slice[i+1].Method
                        matched := false
                        if priority == "month" {
                            if m1 == "hour" && m2 == "month" {
                                copyWant[i].Id = id
                                copyWant[i+1].Id = id
                                i += 2
                                id++
                                matched = true
                            }
                        } else {
                            switch {
                            case m1 == "hour" && m2 == "month":
                                copyWant[i].Id = id
                                copyWant[i+1].Id = id
                                i += 2
                                id++
                                matched = true
                            case m1 == "hour" && m2 == "realtime":
                                copyWant[i].Id = id
                                copyWant[i+1].Id = id
                                i += 2
                                id++
                                matched = true
                            case m1 == "month" && m2 == "realtime":
                                copyWant[i].Id = id
                                copyWant[i+1].Id = id
                                i += 2
                                id++
                                matched = true
                            }
                        }
                        if matched {
                            continue
                        }
                    }
                    copyWant[i].Id = id
                    i++
                    id++
                }

                tests = append(tests, struct {
                    name string
                    args args
                    want []ChargeItem
                }{
                    name: fmt.Sprintf("realtime_seq%d_len%d", seqIdx, l),
                    args: args{slice},
                    want: copyWant,
                })
                seqIdx++
                return
            }
            for _, m := range methods {
                generate(pos+1, append(cur, m))
            }
        }
        generate(0, []string{})
    }

    // --------------------
    // 执行测试
    // --------------------
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := GroupItemIdWithPriority(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GroupItemIdWithPriority() = %v, want %v", got, tt.want)
            }
        })
    }
}
