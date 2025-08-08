package action

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app"
)

type Action interface {
    Do(ctx context.Context, c *app.RequestContext)
}

type Version struct {
    Version string `json:"version"`
    Action  string `json:"action"`
}

func NewVersion(version string, action string) Version {
    return Version{Version: version, Action: action}
}

func NewVersionOne(action string) Version {
    return Version{Version: "v1", Action: action}
}

func NewVersionTwo(action string) Version {
    return Version{Version: "v2", Action: action}
}

var Factory = map[Version]app.HandlerFunc{}

func SelectAction(version Version) app.HandlerFunc {
    action, ok := Factory[version]
    if ok {
        return action
    }
    return nil
}
