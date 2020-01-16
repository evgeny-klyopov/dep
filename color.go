package main

import "fmt"

type Code struct {
	Default string
	Black string
	Red string
	Green string
	Yellow string
	Purple string
	Magenta string
	Teal string
	White string
}

type Color struct {
	Code Code
}

func NewColor() Color {
	return Color{
		Code: Code{
			Default: "\033[0m",
			Black: "\033[1;30m",
			Red: "\033[1;31m",
			Green: "\033[1;32m",
			Yellow: "\033[1;33m",
			Purple: "\033[1;34m",
			Magenta: "\033[1;35m",
			Teal: "\033[1;36m",
			White: "\033[1;37m",
		},
	}
}

func(c *Color) Format(color string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(color,
			fmt.Sprint(args...))
	}
	return sprint
}
func(c *Color) Print(f func(...interface{}) string, s string) {
	fmt.Println(f(s))
}


func(c *Color) Info(message ...interface{}) string {
	return c.Teal(message...)
}
func(c *Color) Success(message ...interface{}) string {
	return c.Green(message...)
}
func(c *Color) Warning(message ...interface{}) string {
	return c.Yellow(message...)
}
func(c *Color) Fatal(message ...interface{}) string {
	return c.Red(message...)
}

func(c *Color) Black(message ...interface{}) string {
	return c.Format(c.Code.Black + "%s" + c.Code.Default)(message...)
}
func(c *Color) Red(message ...interface{}) string {
	return c.Format(c.Code.Red + "%s" + c.Code.Default)(message...)
}
func(c *Color) Green(message ...interface{}) string {
	return c.Format(c.Code.Green + "%s" + c.Code.Default)(message...)
}
func(c *Color) Yellow(message ...interface{}) string {
	return c.Format(c.Code.Yellow + "%s" + c.Code.Default)(message...)
}
func(c *Color) Purple(message ...interface{}) string {
	return c.Format(c.Code.Purple + "%s" + c.Code.Default)(message...)
}
func(c *Color) Magenta(message ...interface{}) string {
	return c.Format(c.Code.Magenta + "%s" + c.Code.Default)(message...)
}
func(c *Color) Teal(message ...interface{}) string {
	return c.Format(c.Code.Teal + "%s" + c.Code.Default)(message...)
}
func(c *Color) White(message ...interface{}) string {
	return c.Format(c.Code.White + "%s" + c.Code.Default)(message...)
}

