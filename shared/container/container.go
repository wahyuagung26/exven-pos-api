package container

import (
	"fmt"
	"reflect"
	"sync"
)

type Container interface {
	Register(name string, factory interface{})
	RegisterSingleton(name string, factory interface{})
	Get(name string) (interface{}, error)
	MustGet(name string) interface{}
	Resolve(target interface{}) error
	Has(name string) bool
}

type DIContainer struct {
	services   map[string]service
	singletons map[string]interface{}
	mu         sync.RWMutex
}

type service struct {
	factory     interface{}
	isSingleton bool
}

func New() *DIContainer {
	return &DIContainer{
		services:   make(map[string]service),
		singletons: make(map[string]interface{}),
	}
}

func (c *DIContainer) Register(name string, factory interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if factory == nil {
		panic(fmt.Sprintf("cannot register nil factory for service: %s", name))
	}

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		panic(fmt.Sprintf("factory must be a function for service: %s", name))
	}

	c.services[name] = service{
		factory:     factory,
		isSingleton: false,
	}
}

func (c *DIContainer) RegisterSingleton(name string, factory interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if factory == nil {
		panic(fmt.Sprintf("cannot register nil factory for service: %s", name))
	}

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		panic(fmt.Sprintf("factory must be a function for service: %s", name))
	}

	c.services[name] = service{
		factory:     factory,
		isSingleton: true,
	}
}

func (c *DIContainer) Get(name string) (interface{}, error) {
	c.mu.RLock()
	svc, exists := c.services[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service not found: %s", name)
	}

	if svc.isSingleton {
		c.mu.RLock()
		if instance, exists := c.singletons[name]; exists {
			c.mu.RUnlock()
			return instance, nil
		}
		c.mu.RUnlock()

		c.mu.Lock()
		defer c.mu.Unlock()

		if instance, exists := c.singletons[name]; exists {
			return instance, nil
		}

		instance := c.createInstance(svc.factory)
		c.singletons[name] = instance
		return instance, nil
	}

	return c.createInstance(svc.factory), nil
}

func (c *DIContainer) MustGet(name string) interface{} {
	instance, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return instance
}

func (c *DIContainer) Resolve(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	targetElem := targetValue.Elem()
	if !targetElem.CanSet() {
		return fmt.Errorf("target cannot be set")
	}

	targetType := targetElem.Type()

	for name, svc := range c.services {
		factoryType := reflect.TypeOf(svc.factory)
		if factoryType.NumOut() == 0 {
			continue
		}

		outputType := factoryType.Out(0)
		if outputType.AssignableTo(targetType) {
			instance, err := c.Get(name)
			if err != nil {
				return err
			}
			targetElem.Set(reflect.ValueOf(instance))
			return nil
		}
	}

	return fmt.Errorf("no service found for type: %s", targetType.String())
}

func (c *DIContainer) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	return exists
}

func (c *DIContainer) createInstance(factory interface{}) interface{} {
	factoryValue := reflect.ValueOf(factory)
	factoryType := factoryValue.Type()

	args := make([]reflect.Value, factoryType.NumIn())
	for i := 0; i < factoryType.NumIn(); i++ {
		argType := factoryType.In(i)

		for _, svc := range c.services {
			svcType := reflect.TypeOf(svc.factory)
			if svcType.NumOut() > 0 && svcType.Out(0) == argType {
				instance, _ := c.Get(c.getServiceName(svc.factory))
				args[i] = reflect.ValueOf(instance)
				break
			}
		}

		if !args[i].IsValid() {
			args[i] = reflect.Zero(argType)
		}
	}

	results := factoryValue.Call(args)
	if len(results) > 0 {
		return results[0].Interface()
	}
	return nil
}

func (c *DIContainer) getServiceName(factory interface{}) string {
	for name, svc := range c.services {
		if reflect.ValueOf(svc.factory).Pointer() == reflect.ValueOf(factory).Pointer() {
			return name
		}
	}
	return ""
}

type ServiceRegistrar interface {
	Register(container Container)
}

func RegisterModules(container Container, modules ...ServiceRegistrar) {
	for _, module := range modules {
		module.Register(container)
	}
}
