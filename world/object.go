package world

type Object struct {
	Name         string
	Description  string
	actions      map[string]action
	interactions map[string]interaction
	pos          *position
}

type position struct {
	relation string
	object   *Object
}

type action struct {
	name   string
	action actionFunc
}

type actionFunc func(c *Character) string

type interaction struct {
	name   string
	action actionFunc
}

func NewObject(name, desc string) *Object {
	return &Object{
		Name:         name,
		Description:  desc,
		actions:      make(map[string]action),
		interactions: make(map[string]interaction),
	}
}

func (o *Object) SetPosition(relation string, object *Object) {
	if object == nil {
		o.pos = nil
		return
	}
	o.pos = &position{relation: relation, object: object}
}

func (o *Object) GetPosition() (string, *Object) {
	if o.pos != nil {
		return o.pos.relation, o.pos.object
	}
	return "", nil
}

func (o *Object) SetAction(name string, actionFunc actionFunc) {
	o.actions[name] = action{name: name, action: actionFunc}
}

func (o *Object) SetInteraction(name string, actionFunc actionFunc) {
	o.interactions[name] = interaction{name: name, action: actionFunc}
}

func (o *Object) Examine() string {
	resp := "You see a " + o.Description + ".\n"
	resp += "You can interact with it in the following ways:\n"
	for _, action := range o.actions {
		resp += action.name + "\n"
	}
	return resp
}

func (o *Object) Describe() string {
	resp := " a " + o.Name
	if o.pos != nil {
		resp += " " + o.pos.relation + " the " + o.pos.object.Name + ".\n"
	}

	return resp
}

func (o *Object) PerformAction(c *Character, action string) string {
	var resp string
	if act, ok := o.actions[action]; ok {
		resp = act.action(c)
	} else {
		resp = "You can't do that."
	}
	return resp + "\n" + o.Examine()
}

func (o *Object) PerformInteraction(c *Character, interaction string) string {
	if act, ok := o.interactions[interaction]; ok {
		return act.action(c)
	}
	return "You can't do that."
}
