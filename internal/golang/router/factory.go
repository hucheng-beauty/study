package router

import (
    "net/url"

    "github.com/gin-gonic/gin"
    "github.com/samber/lo"
)

type Factory struct {
    engine          *gin.Engine
    url             *url.URL
    baseMiddleware  []gin.HandlerFunc
    innerMiddleware []gin.HandlerFunc
    openMiddleware  []gin.HandlerFunc
}

type Routers struct {
    InnerRouter *gin.RouterGroup
    OpenRouter  *gin.RouterGroup
}

func NewFactory(e *gin.Engine) *Factory { return &Factory{engine: e, url: &url.URL{}} }

func (f *Factory) Clone() *Factory {
    cloned := lo.ToPtr(*f)

    cloned.url = lo.ToPtr(*cloned.url)
    cloned.baseMiddleware = f.cloneMiddleware(cloned.baseMiddleware)
    cloned.innerMiddleware = f.cloneMiddleware(cloned.innerMiddleware)
    cloned.openMiddleware = f.cloneMiddleware(cloned.openMiddleware)

    return cloned
}

func (f *Factory) AddPath(path string) *Factory {
    f.url = f.url.JoinPath(path)
    return f
}

func (f *Factory) AppendBaseMiddleware(fn ...gin.HandlerFunc) *Factory {
    f.baseMiddleware = append(f.baseMiddleware, fn...)
    return f
}

func (f *Factory) AppendInnerMiddleware(fn ...gin.HandlerFunc) *Factory {
    f.innerMiddleware = append(f.innerMiddleware, fn...)
    return f
}

func (f *Factory) AppendOpenMiddleware(fn ...gin.HandlerFunc) *Factory {
    f.openMiddleware = append(f.openMiddleware, fn...)
    return f
}

func (f *Factory) Routers() *Routers {
    baseRouter := f.engine.Group(f.url.String(), f.baseMiddleware...)
    return &Routers{
        InnerRouter: baseRouter.Group("/inner", f.innerMiddleware...),
        OpenRouter:  baseRouter.Group("/open", f.openMiddleware...),
    }
}

func (f *Factory) cloneMiddleware(in []gin.HandlerFunc) (out []gin.HandlerFunc) {
    return append(out, in...)
}
