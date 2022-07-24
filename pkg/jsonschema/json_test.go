package jsonschema

import (
	"encoding/json"
	"testing"
)

var tdData = `
{
"hostname" : "lazar",
"messaging": {
		"cluster":"google.com",
		"network":"mainnet"
	},
"network" : [
		{
			"nic":"enp36s0",
			"ipv4":"192.168.1.2/24",
			"gateway":"192.168.1.1"
		},
		{
			"nic":"lo",
			"ipv4":"192.168.2.2/24",
			"gateway":"192.168.2.1"
		}
	],
"disks":[
		"/dev/sda1",
		"/dev/sda2",
		"/dev/sda3"
	],
"number":30
}`

var tdSchema = `
{
"hostname%required" : "hostname%required",
"messaging%required": {
		"cluster":"lookup4%required",
		"network":"string%required"
	},
"network%required" : [
		{
			"nic":"nic4%required",
			"ipv4":"ipv4net%required",
			"gateway":"ipv4%required"
		}
	],
"disks%required":[
		"string"
	],
"filesystems":[
		"string"
	],
"numbers":[
		"float"
	],
"number": "int(min=1)"
}`

func TestJson(t *testing.T) {
	var schema, data interface{}
	if err := json.Unmarshal([]byte(tdSchema), &schema); err != nil {
		t.Fatalf("Unmarshal Schema: %s", err)
	}
	if err := json.Unmarshal([]byte(tdData), &data); err != nil {
		t.Fatalf("Unmarshal Data: %s", err)
	}
	if errPath, modifiedData, err := Validate(
		schema,
		data,
		//schema.(map[string]interface{})["messaging"],
		//data.(map[string]interface{})["messaging"],
	); err != nil {
		_ = modifiedData
		t.Errorf("Validate: %s %s", errPath, err)
		//} else {
		//	spew.Dump(modifiedData)
	}
}
