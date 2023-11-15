package main

import "ddns/routes"

func main() {
	r := routes.DnspodRouter()
	r.Run(":88")
}
