package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s can %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s cannot %s %s\n", sub, act, obj)
	}
}

func main() {
	e, err := casbin.NewEnforcer("./internal/rbac/model.conf", "./internal/rbac/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}
	check(e, "zxp", "data1", "read")
	check(e, "zhang", "data2", "write")
	check(e, "zxp", "data1", "write")
	check(e, "zxp", "data2", "read")
}
