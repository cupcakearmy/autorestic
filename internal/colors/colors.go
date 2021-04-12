package colors

import (
	"github.com/fatih/color"
)

var Body = color.New()
var Primary = color.New(color.Underline, color.Bold, color.BgBlue)
var Secondary = color.New(color.Bold, color.FgCyan)
var Success = color.New(color.FgGreen)
var Error = color.New(color.FgRed, color.Bold)
var Faint = color.New(color.Faint)
