package color

import "fmt"

//import "bytes"

func PrintPrompt(format string, args ...interface{}) {
	//	var buf bytes.Buffer
	//	buf.WriteString("\033[7m")
	//	buf.WriteString(format)
	//	buf.WriteString("\033[0m")
	//	format = buf.String()
	format = "\033[7m" + format + "\033[0m"
	fmt.Printf(format, args...)
}

func PrintError(format string, args ...interface{}) {
	format = "\033[41m" + format + "\033[0m"
	fmt.Printf(format, args...)
}

func PrintErrorln(format string, args ...interface{}) {
	format = "\033[41m" + format + "\033[0m\n"
	fmt.Printf(format, args...)
}

func PrintBlue(format string, args ...interface{}) {
	format = "\033[44m" + format + "\033[0m"
	fmt.Printf(format, args...)
}

func PrintBlueln(format string, args ...interface{}) {
	format = "\033[44m" + format + "\033[0m\n"
	fmt.Printf(format, args...)
}
