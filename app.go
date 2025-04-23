package main

import (
	"context"
	"encoding/json"
	server "myapp/src"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(funcName, param string) interface{} {
	res := make(map[string]interface{})
	fn, ok := server.FuncMap[funcName]
	if !ok {
		res["data"] = nil
		res["error"] = "func not found"
		return res
	}
	jsonData := make(map[string]interface{})
	err := json.Unmarshal([]byte(param), &jsonData)
	if err != nil {
		return err.Error()
	}
	result := fn(jsonData)
	res["data"] = result
	return res
}
