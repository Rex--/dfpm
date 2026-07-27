package main

var Dfpm_Version = "1.0"
var Dfpm_Year = "2026"

func init() {
	InitCLI()
}

func main() {
	args := ParseCLI()

	switch args.Command {
	case "install":
		Dfpm_install(args)
	case "list":
		Dfpm_list(args)
	case "search":
		Dfpm_search(args)
	// case "update":
	// 	// Dfpm_update()
	// 	fmt.Println("Updating packs is currently not supported. Please re-download the pack for the latest version")
	case "xc8":
		Dfpm_xc8(args)
	default:
		usageErr()
	}

}
