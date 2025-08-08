package registry

import (
    "reflect"
    "strings"
)

// tagInfo 包含从 db 标签解析出的信息
type tagInfo struct {
    FieldName string
    // 可以根据需要添加更多字段，如 omitempty 等选项
}

// dbTagInfo 解析结构体字段的 db 标签信息
func dbTagInfo(field reflect.StructField) (tagInfo, bool, error) {
    // 获取字段的 db 标签
    dbTag := field.Tag.Get("db")
    if dbTag == "" {
        // 如果没有 db 标签，返回 false 表示不处理该字段
        return tagInfo{}, false, nil
    }

    // db 标签通常格式为 `db:"field_name"` 或 `db:"field_name,opt1,opt2"`
    // 我们只需要字段名部分
    parts := strings.Split(dbTag, ",")
    fieldName := strings.TrimSpace(parts[0])

    // 如果字段名为空，使用结构体字段名作为默认值
    if fieldName == "" {
        fieldName = field.Name
    }

    return tagInfo{FieldName: fieldName}, true, nil
}
