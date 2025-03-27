package list

import (
    "testing"
)

func Test_findKey(t *testing.T) {
    list := &List{
        Data: 1,
        Next: &List{
            Data: 2,
            Next: &List{
                Data: 3,
                Next: nil,
            },
        },
    }

    t.Logf("%#v", findKey(list, 5))
    t.Logf("%#v", findKey(list, 3))
    t.Logf("%#v", findKey(list, 1))
}
