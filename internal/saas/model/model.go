package model

type Order struct {
    Type         string
    InstanceNO   string
    SubOrder     string
    IsResource   int
    InstanceType int
}

type ResourceTag string
