package world

import (
	"fmt"
	"strings"

	"go.dendromeda.gpt-mud/gpt"
)

type Character struct {
	name      string
	initials  []string
	Chat      *gpt.Chat
	Inventory []*Object
	pos       position
}

const initialisation = `You are a character in a text-based adventure game. You can type commands to interact with the world around you.
The first word in your respond is the command you want to execute. There is one action for each response.
To talk, you respond with "say" followed by your message.
To investigate an object, you respond with "examine" followed by the object name. This will give you additional commands specific for the object.
You can use your items on objects by responding with "use" followed by the item name and then the object name.
Remember to examine objects after interacting with them, the available actions might have changed.
Do not write anything other than the command in your response.
Actions will alwatys be in the form of a command followed by the object you want to interact with.
You can only interact with the objects that are listed as being in the scene.
`

func NewCharacter(client *gpt.Client, initials []string, name string) (*Character, error) {

	var chat *gpt.Chat
	fullIntitals := strings.Join(append(initials, initialisation, fmt.Sprintf("Your name is \"%s\"\n", name)), "\n")
	fmt.Println(fullIntitals)
	if client != nil {
		chat = gpt.NewChat(client, "gpt-4-turbo")
		_, err := chat.Chat("system", fullIntitals)
		if err != nil {
			return nil, err
		}
	}

	return &Character{
		name:     name,
		initials: initials,
		Chat:     chat,
	}, nil
}

func (c *Character) Prompt(entry string) (string, error) {
	return c.Chat.Chat("user", entry)
}

func (c *Character) Describe() string {
	return c.name
}

func (c *Character) SetRelation(object *Object) {
	c.Inventory = append(c.Inventory, object)
}
