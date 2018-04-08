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
	//fmt.Printf("\033[7m")
	if len(args) == 0 {
		fmt.Print(format)
		//		fmt.Printf("\033[0m")
		return
	}
	fmt.Printf(format, args...)
	//	fmt.Printf("\033[0m")
}

func PrintError(format string, args ...interface{}) {
	format = "\033[41m" + format + "\033[0m"
	if len(args) == 0 {
		fmt.Print(format)
		return
	}
	fmt.Printf(format, args...)
}

func PrintErrorln(format string, args ...interface{}) {
	format = "\033[41m" + format + "\033[0m\n"
	if len(args) == 0 {
		fmt.Print(format)
		return
	}
	fmt.Printf(format, args...)
}
