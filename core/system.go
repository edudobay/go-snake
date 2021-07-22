package core

import "fmt"

// https://en.wikipedia.org/wiki/Entity_component_system

type Component interface {
	Type() string
}

type Entity interface {
	Id() string
	AttachComponent(component Component)
	HasComponent(withType string) bool
	GetComponent(withType string) Component
}

type System interface {
	AddEntity(entity Entity)
	FindEntityOrNilById(id string) Entity
	FindComponentsOfType(withType string) []Component
	OneComponentOfType(withType string) Component
}

type systemImpl struct {
	entities map[string]Entity
}

type entityImpl struct {
	id         string
	Components map[string]Component
}

func NewSystem() System {
	return &systemImpl{
		entities: map[string]Entity{},
	}
}

func NewEntity(id string) Entity {
	return &entityImpl{
		id:         id,
		Components: make(map[string]Component),
	}
}

func (e *entityImpl) Id() string {
	return e.id
}

func (e *entityImpl) HasComponent(withType string) bool {
	_, ok := e.Components[withType]
	return ok
}

func (e *entityImpl) GetComponent(withType string) Component {
	return e.Components[withType]
}

func (e *entityImpl) AttachComponent(component Component) {
	// TODO: Silently overwriting an existing component
	e.Components[component.Type()] = component
}

func (s *systemImpl) AddEntity(entity Entity) {
	s.entities[entity.Id()] = entity
}

func (s *systemImpl) FindEntityOrNilById(id string) Entity {
	return s.entities[id]
}

func (s *systemImpl) FindComponentsOfType(withType string) []Component {
	components := make([]Component, 0, len(s.entities))

	for _, entity := range s.entities {
		if component := entity.GetComponent(withType); component != nil {
			components = append(components, component)
		}
	}

	return components
}

func (s *systemImpl) OneComponentOfType(withType string) Component {
	var found Component

	for _, entity := range s.entities {
		if found != nil {
			panic(fmt.Sprintf("found more than one component of type '%s'", withType))
		}

		if component := entity.GetComponent(withType); component != nil {
			found = component
		}
	}

	if found == nil {
		panic(fmt.Sprintf("no '%s' components found", withType))
	}

	return found
}
