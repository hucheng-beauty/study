package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type StringSlice []string

func (ss *StringSlice) Scan(src any) error {
    if src == nil {
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        return fmt.Errorf("invalid StringSlice.Scan source type: %T", src)
    }

    if len(srcBytes) == 0 {
        return nil
    }

    return json.Unmarshal(srcBytes, ss)
}

func (ss StringSlice) Value() (driver.Value, error) { return json.Marshal(ss) }
