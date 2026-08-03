package utils

import (
	"agent/internal/model"
	"flag"
	"fmt"
)

func ReadArgs(args []model.ArgInfo) (res map[string]*string, err error) {
	res = make(map[string]*string)
	for _, arg := range args {
		res[arg.Name] = flag.String(arg.Name, arg.DefaultValue, arg.Description)
	}

	flag.Parse()
	for _, arg := range args {
		if arg.Required && res[arg.Name] == &arg.DefaultValue {
			return nil, fmt.Errorf("flag '%s' is required", arg.Name)
		}
	}
	return res, nil
}
