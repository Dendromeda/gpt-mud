package world

import (
	"strings"
)

type Scene struct {
	Description string
	Objects     []*Object
}

func NewScene(desc string, objects []*Object) *Scene {
	return &Scene{
		Description: desc,
		Objects:     objects,
	}
}

func (s *Scene) RemoveObject(name string) {
	for i, o := range s.Objects {
		if o.Name == name {
			var trail []*Object
			if i < len(s.Objects)-1 {
				trail = s.Objects[i+1:]
			}
			s.Objects = append(s.Objects[:i], trail...)
			return
		}
	}
}

func (s *Scene) Describe(c *Character) string {
	resp := s.Description + "\n"
	for _, obj := range s.Objects {
		resp += "There is" + obj.Describe() + "\n"
	}
	if len(c.Inventory) > 0 {
		resp += "You have the following items in your inventory:\n"
		for _, e := range c.Inventory {
			resp += e.Name + "\n"
		}
	}

	resp += "You can examine objects by typing \"examine [object] to get available actions\n"
	resp += "You can use inventory objects with other objects by typing \"use [inventory object] [object]\n"

	return resp
}

func (s *Scene) PerformAction(c *Character, action string) string {

	actionParts := strings.Split(action, " ")

	if actionParts[0] == "examine" && len(actionParts) == 2 {
		for _, obj := range s.Objects {
			if obj.Name == actionParts[1] {
				return obj.Examine()
			}
		}
	}

	if actionParts[0] == "use" && len(actionParts) == 3 {
		for _, obj := range c.Inventory {
			if obj.Name == actionParts[1] {
				for _, obj2 := range s.Objects {
					if obj2.Name == actionParts[2] {
						return obj2.PerformInteraction(c, obj.Name)
					}
				}
			}
		}

	}

	for _, obj := range s.Objects {
		if len(actionParts) == 2 && obj.Name == actionParts[1] {
			return obj.PerformAction(c, actionParts[0])
		}
	}
	return "You can't do that."
}
