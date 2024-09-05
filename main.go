package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.dendromeda.gpt-mud/gpt"
	"go.dendromeda.gpt-mud/world"
)

//go:file secrets.json

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func main() {
	//Use your API KEY here
	b, err := os.ReadFile("secrets.json")
	if err != nil {
		panic(err)
	}
	secrets := make(map[string]string)
	err = json.Unmarshal(b, &secrets)
	if err != nil {
		panic(err)
	}

	apiKey := secrets["api-key"]

	//Create a new client
	client := gpt.NewClient(apiKey, apiEndpoint)

	// Create a new character
	character, err := world.NewCharacter(client, []string{"Your goal is to get out of the room"}, "Guy")

	if err != nil {
		panic(err)
	}

	scene := Room2()

	var result string

	for {
		fmt.Println(result + "\n" + scene.Describe(character))
		resp, err := character.Prompt(result + "\n" + scene.Describe(character))
		if err != nil {
			panic(err)
		}
		fmt.Println(">" + resp)

		result = scene.PerformAction(character, resp) + "\n"

		time.Sleep(time.Second)
	}

}

func Room1() *world.Scene {
	table := world.NewObject("table", "a wooden table")
	key := world.NewObject("key", "a shiny key")
	door := world.NewObject("door", "a locked door")
	scene := world.NewScene("You are in a big room", []*world.Object{table, key, door})

	key.SetPosition("on", table)
	key.SetAction("take", func(c *world.Character) string {
		key.SetPosition("", nil)
		c.Inventory = append(c.Inventory, key)
		scene.RemoveObject("key")
		return "You take the key"
	})

	door.SetAction("open", func(c *world.Character) string {
		if door.Description == "an unlocked door" {
			door.Description = "an open door"
			door.SetAction("exit", func(c *world.Character) string {
				fmt.Print("Congratulations! You exit the room")
				os.Exit(0)
				return ""
			})
			return "You open the door"
		}
		return "The door is locked"
	})

	door.SetInteraction("key", func(c *world.Character) string {
		door.Description = "an unlocked door"
		for i, o := range c.Inventory {
			if o.Name == "key" {
				var trail []*world.Object
				if i < len(c.Inventory)-1 {
					trail = c.Inventory[i+1:]
				}
				c.Inventory = append(c.Inventory[:i], trail...)
			}
		}

		return "You unlock the door"
	})

	table.SetAction("flip", func(c *world.Character) string {
		for _, o := range scene.Objects {
			if rel, obj := o.GetPosition(); rel == "on" && obj.Name == "table" {
				o.SetPosition("on", world.NewObject("floor", ""))
			}
		}
		return "You flip the table over. Anything on it falls to the floor."
	})
	return scene
}

func Room2() *world.Scene {
	table := world.NewObject("table", "a wooden table")
	key := world.NewObject("key", "a shiny key")
	door := world.NewObject("door", "a locked door")
	plate := world.NewObject("plate", "a metal plate that seem to be electrified")
	scene := world.NewScene("You are in a big room", []*world.Object{table, key, door, plate})

	key.SetPosition("on", plate)
	plate.SetPosition("on", table)
	key.SetAction("take", func(c *world.Character) string {
		if _, o := key.GetPosition(); o.Name == "plate" {
			return "The metal plate the key sits on zaps you"
		}
		key.SetPosition("", nil)
		c.Inventory = append(c.Inventory, key)
		scene.RemoveObject("key")
		return "You take the key"
	})

	door.SetAction("open", func(c *world.Character) string {
		if door.Description == "an unlocked door" {
			door.Description = "an open door"
			door.SetAction("exit", func(c *world.Character) string {
				fmt.Print("Congratulations! You exit the room")
				os.Exit(0)
				return ""
			})
			return "You open the door"
		}
		return "The door is locked"
	})

	door.SetInteraction("key", func(c *world.Character) string {
		door.Description = "an unlocked door"
		for i, o := range c.Inventory {
			if o.Name == "key" {
				var trail []*world.Object
				if i < len(c.Inventory)-1 {
					trail = c.Inventory[i+1:]
				}
				c.Inventory = append(c.Inventory[:i], trail...)
			}
		}

		return "You unlock the door"
	})

	table.SetAction("flip", func(c *world.Character) string {
		for _, o := range scene.Objects {
			if rel, obj := o.GetPosition(); rel == "on" && (obj.Name == "plate" || obj.Name == "table") {
				o.SetPosition("on", world.NewObject("floor", ""))
			}
		}
		return "You flip the table over. Anything on it falls to the floor."
	})
	return scene
}

func Room3() *world.Scene {
	return nil
}
