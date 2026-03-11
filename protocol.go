package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tidwall/resp"
)

const (
	CMD_SET = "set"
	CMD_GET = "get"
)

type Command interface {
	Name() string
}

type SetCommand struct {
	Key   string
	Value string
}

func (c SetCommand) Name() string { return CMD_SET }

type GetCommand struct {
	Key string
}

func (c GetCommand) Name() string { return CMD_GET }

func parseCommand(msg string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(msg))
	v, _, err := rd.ReadValue()
	if err != nil {
		return nil, fmt.Errorf("failed to read resp value: %w", err)
	}

	if v.Type() != resp.Array {
		return nil, fmt.Errorf("expected array, got %s", v.Type())
	}

	arr := v.Array()
	if len(arr) == 0 {
		return nil, fmt.Errorf("empty command array")
	}

	cmdName := strings.ToLower(arr[0].String())
	switch cmdName {
	case CMD_SET:
		if len(arr) != 3 {
			return nil, fmt.Errorf("SET requires 2 args, got %d", len(arr)-1)
		}
		return SetCommand{
			Key:   arr[1].String(),
			Value: arr[2].String(),
		}, nil

	case CMD_GET:
		if len(arr) != 2 {
			return nil, fmt.Errorf("GET requires 1 arg, got %d", len(arr)-1)
		}
		return GetCommand{Key: arr[1].String()}, nil

	default:
		return nil, fmt.Errorf("unsupported command: %q", cmdName)
	}
}
