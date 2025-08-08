package router

import (
    "net/url"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/route"
    "github.com/samber/lo"
)

type Factory struct {
    engine          *server.Hertz
    url             *url.URL
    baseMiddleware  []app.HandlerFunc
    innerMiddleware []app.HandlerFunc
    openMiddleware  []app.HandlerFunc
}

type Routers struct {
    InnerRouter *route.RouterGroup
    OpenRouter  *route.RouterGroup
}

func NewFactory(e *server.Hertz) *Factory { return &Factory{engine: e, url: &url.URL{}} }

func (f *Factory) Clone() *Factory {
    cloned := lo.ToPtr(*f)

    cloned.url = lo.ToPtr(*cloned.url)
    cloned.baseMiddleware = f.cloneMiddleware(cloned.baseMiddleware)
    cloned.innerMiddleware = f.cloneMiddleware(cloned.innerMiddleware)
    cloned.openMiddleware = f.cloneMiddleware(cloned.openMiddleware)

    return cloned
}

func (f *Factory) cloneMiddleware(in []app.HandlerFunc) (out []app.HandlerFunc) {
    return append(out, in...)
}

func (f *Factory) AddPath(path string) *Factory {
    f.url = f.url.JoinPath(path)
    return f
}

func (f *Factory) AppendBaseMiddleware(fn ...app.HandlerFunc) *Factory {
    f.baseMiddleware = append(f.baseMiddleware, fn...)
    return f
}

func (f *Factory) AppendInnerMiddleware(fn ...app.HandlerFunc) *Factory {
    f.innerMiddleware = append(f.innerMiddleware, fn...)
    return f
}

func (f *Factory) AppendOpenMiddleware(fn ...app.HandlerFunc) *Factory {
    f.openMiddleware = append(f.openMiddleware, fn...)
    return f
}

func (f *Factory) Groups() *Routers {
    baseRouter := f.engine.Group(f.url.String(), f.baseMiddleware...)
    return &Routers{
        InnerRouter: baseRouter.Group("/inner", f.innerMiddleware...),
        OpenRouter:  baseRouter.Group("/open", f.openMiddleware...),
    }
}
