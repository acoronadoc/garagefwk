package main

import (
	"fmt"
	"garagefwk"
	"time"
)

func main() {
	app := make(map[string]interface{})
	app["screens"] = map[string]interface{}{
		"function": Dashboard,
	}

	go garagefwk.InitGarageFWK(&app)

	for {
		time.Sleep(time.Hour * 999999)
	}
}

func Dashboard(parts *[]map[string]interface{}, request *garagefwk.GarageRequest) {

	*parts = append(*parts, map[string]interface{}{
		"component": "clean",
	})

	*parts = append(*parts, map[string]interface{}{
		"component": "renderTag",
		"tag":       "h2",
		"innerHTML": "Dashboard",
	})

	regs := garagefwk.ReadDataObjectsByFilter(request.Db, "", "server", &[]garagefwk.DataObjectFilter{})
	for _, regp := range *regs {
		reg := (*regp.Reg)
		html := ""
		html += fmt.Sprint("<div><strong>", reg["name"], "</strong> (", reg["status"], ")</div>")
		html += fmt.Sprint("<div>", reg["_lastcheck"], " (", reg["_duration"], " milis.)</div>")
		html += fmt.Sprint("<div>", reg["_results"], "</div>")

		*parts = append(*parts, map[string]interface{}{
			"component": "renderTag",
			"tag":       "div",
			"style":     "margin: 15px; padding: 15px; border: 1px solid #666; border-radius: 5px;",
			"innerHTML": html,
		})
	}

	regs = garagefwk.ReadDataObjectsByFilter(request.Db, "", "code", &[]garagefwk.DataObjectFilter{})
	for _, regp := range *regs {
		reg := (*regp.Reg)
		html := ""
		html += fmt.Sprint("<div>", reg["name"], "</div>")

		*parts = append(*parts, map[string]interface{}{
			"component": "renderTag",
			"tag":       "div",
			"innerHTML": html,
		})
	}
}
