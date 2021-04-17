package colors

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var Body = color.New()
var Primary = color.New(color.Bold, color.BgBlue, color.FgHiWhite)
var Secondary = color.New(color.Bold, color.FgCyan)
var Success = color.New(color.FgGreen)
var Error = color.New(color.FgRed, color.Bold)
var Faint = color.New(color.Faint)

func PrimaryPrint(msg string, args ...interface{}) {
	fmt.Printf("\n\n%s\n\n", Primary.Sprintf("  "+msg+"  ", args...))
}

func DisableColors(state bool) {
	color.NoColor = state
}

func PrintDescription(left string, right string) {
	right = strings.Trim(right, "\n")
	right = strings.Trim(right, "\t")
	Body.Printf("%s\t%s\n", Secondary.Sprint(left), right)
}
