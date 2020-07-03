package main

type Registry struct {
	objects map[string]SimpleData
}

func NewRegistry() *Registry{
	return &Registry{
		objects: make(map[string]SimpleData),
	}
}

// 需要传入object指针
func (registry *Registry)Register(objectName string, objectPtr SimpleData){
	_, exist := registry.objects[objectName]
	if exist {
		return
	}
	registry.objects[objectName] = objectPtr
}

func (registry *Registry)Remove(objectName string){
	delete(registry.objects, objectName)
}

func (registry *Registry)GetObject(key string)SimpleData{
	return registry.objects[key]
}

func (registry *Registry)GetAllObjects() []SimpleData{
	objects := make([]SimpleData, 0, len(registry.objects))
	for _, object := range registry.objects{
		objects = append(objects, object)
	}
	return objects
}

func (registry *Registry)GetAllNames()[]string{
	names := make([]string, 0, len(registry.objects))
	for name := range registry.objects{
		names = append(names, name)
	}
	return names
}