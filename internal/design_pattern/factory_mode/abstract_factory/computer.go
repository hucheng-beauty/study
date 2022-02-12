package abstract_factory

type Product interface {
	SetName(name string)
	GetName() string
}

type HuaWei struct {
	name string
}

func (h *HuaWei) SetName(name string) {
	h.name = name
}

func (h *HuaWei) GetName() string {
	return h.name
}

type Apple struct {
	name string
}

func (a *Apple) SetName(name string) {
	a.name = name
}

func (a *Apple) GetName() string {
	return a.name
}

type ComputerFactory struct{}

func (f *ComputerFactory) Create(parameterType string) Product {
	switch parameterType {
	case "HuaWei":
		return &HuaWei{}
	case "Apple":
		return &Apple{}
	default:
		return nil
	}
}
