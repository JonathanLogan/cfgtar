package schemareg

type dirRegistry struct {
	data  interface{}
	byDir map[string]*dirRegistry
}

func (reg *dirRegistry) Add(path []string, data interface{}) {
	if len(path) == 0 {
		reg.data = data
		return
	}
	if reg.byDir == nil {
		reg.byDir = make(map[string]*dirRegistry)
	}
	if _, ok := reg.byDir[path[0]]; !ok {
		reg.byDir[path[0]] = new(dirRegistry)
	}
	reg.byDir[path[0]].Add(path[1:], data)
}

func (reg *dirRegistry) Get(path []string) interface{} {
	if len(path) == 0 || reg.byDir == nil {
		return reg.data
	}
	if sub, ok := reg.byDir[path[0]]; ok {
		if r := sub.Get(path[1:]); r != nil {
			return r
		}
		return reg.data
	}
	return nil
}

type Registry struct {
	reg  *dirRegistry
	data interface{}
}

func New(data interface{}) *Registry {
	return &Registry{
		data: data,
	}
}

func (reg *Registry) Add(path []string, data interface{}) {
	if reg.reg == nil {
		reg.reg = new(dirRegistry)
	}
	reg.reg.Add(path, data)
}

func (reg *Registry) Get(path []string) interface{} {
	if reg.reg == nil {
		reg.reg = new(dirRegistry)
	}
	r := reg.reg.Get(path)
	if r != nil {
		return r
	}
	return reg.data
}
