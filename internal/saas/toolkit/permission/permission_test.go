package permission

import (
    "fmt"
    "testing"
)

func TestHolder_level(t *testing.T) {
    h := Holder{
        Permissions: []Permission{
            // {Description: Description{Action: "admin"}},
            // {Description: Description{Action: "read"}},
            {Description: Description{Action: "write"}},
        },
    }
    fmt.Println("holder.level(): ", h.level())
}
