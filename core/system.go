package core

import "fmt"

// https://en.wikipedia.org/wiki/Entity_component_system

type Component interface {
	Type() string
}

type Entity struct {
	Id         string
	Components map[string]Component
}

type System struct {
	Entities []Entity
}

func NewEntity(id string) Entity {
	return Entity{
		Id:         id,
		Components: make(map[string]Component),
	}
}

func (e *Entity) HasComponent(withType string) bool {
	_, ok := e.Components[withType]
	return ok
}

func (e *Entity) GetComponent(withType string) Component {
	return e.Components[withType]
}

func (e *Entity) AttachComponent(component Component) {
	// TODO: Silently overwriting an existing component
	e.Components[component.Type()] = component
}

func (s *System) AddEntity(entity Entity) {
	s.Entities = append(s.Entities, entity)
}

func (s *System) FindComponentsOfType(withType string) []Component {
	components := make([]Component, 0, len(s.Entities))

	for _, entity := range s.Entities {
		if component := entity.GetComponent(withType); component != nil {
			components = append(components, component)
		}
	}

	return components
}

func (s *System) OneComponentOfType(withType string) Component {
	var found Component

	for _, entity := range s.Entities {
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
