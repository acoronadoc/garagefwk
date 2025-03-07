package garagefwk

import (
	"fmt"
	"reflect"
)

func initScreensTypes(screens *map[string]interface{}) *map[string]interface{} {
	app := make(map[string]interface{})

	app["list"] = listScreenType
	app["form"] = formScreenType

	if screens != nil {
		for key, value := range *screens {
			app[key] = value
		}
	}

	return &app
}

func execScreenType(parts *[]map[string]interface{}, request *GarageRequest) {
	f, exist := (*request.ScreenTypes)[request.Screen.Scrtype]

	if exist {
		myFunc := reflect.ValueOf(f)
		myFunc.Call([]reflect.Value{
			reflect.ValueOf(parts),
			reflect.ValueOf(request),
		})
		//ret := myResult.Interface().(map[string]interface{})
		//return &ret
	}
}

func listScreenType(parts *[]map[string]interface{}, request *GarageRequest) {

	table := ""
	dos := readDataObjects(request.Db, request.User, request.Screen.Options["dataobject"].(string))

	for _, field := range request.Screen.Options["columns"].([]interface{}) {
		table += fmt.Sprint("<th>", field.(map[interface{}]interface{})["label"], "</th>")
	}
	table = fmt.Sprint("<tr>", table, "</tr>")

	for _, row := range *dos {
		l := ""
		r := *row.Reg
		urlVars := map[string]interface{}{
			"id": row.Id,
		}
		url := setVarsURL(request.Screen.Options["regurl"].(string), &urlVars)

		for _, field := range request.Screen.Options["columns"].([]interface{}) {
			f := field.(map[interface{}]interface{})["name"].(string)

			l += fmt.Sprint("<td><a href='", url, "'>", r[f], "</a></td>")
		}

		table += fmt.Sprint("<tr>", l, "</tr>")
	}

	renderHeader(parts, request)

	*parts = append(*parts, map[string]interface{}{
		"component":   "renderTag",
		"tag":         "table",
		"cellspacing": "0",
		"class":       "app-table",
		"innerHTML":   table,
	})
}

func formScreenType(parts *[]map[string]interface{}, request *GarageRequest) {
	idReg, idRegOk := (*(*request).VarsUrl)["id"]

	if (*request).Body != nil {
		saveReg, saveRegOk := (*(*request).Body)["saveForm"]
		checkReg, checkRegOk := (*(*request).Body)["checkForm"]

		if saveRegOk && idRegOk {
			formScreenSave(parts, request, idReg.(string), &saveReg)
			return
		}

		if checkRegOk && idRegOk {
			formScreenCheck(parts, request, &checkReg)
			return
		}

	}

	renderHeader(parts, request)

	*parts = append(*parts, map[string]interface{}{
		"component": "renderTag",
		"tag":       "form",
		"id":        "mainform",
		"innerHTML": "",
	})

	var reg *DataObject
	if idReg != "-1" {
		reg = readDataObject(request.Db, idReg.(string))
	}

	for _, field := range request.Screen.Options["fields"].([]interface{}) {
		name := interfaceValueStr(&field, "name", "")
		ftype := interfaceValueStr(&field, "type", "textbox")
		options := ""

		if ftype == "select" {
			for i, op := range field.(map[interface{}]interface{})["options"].([]interface{}) {
				id := interfaceValueStr(&op, "value", "")
				label := interfaceValueStr(&op, "label", id)

				if i != 0 {
					options += ","
				}

				options += "{ \"id\": \"" + id + "\", \"label\":\"" + label + "\" }"
			}
		}

		*parts = append(*parts, map[string]interface{}{
			"component": "renderTag",
			"tag":       "app-input",
			"selector":  "#mainform",
			"name":      name,
			"options":   "[" + options + "]",
			"onchange":  "javascript: checkForm(event);",
			"label":     interfaceValueStr(&field, "label", name),
			"type":      ftype,
			"value":     DataObjectValueStr(reg, name, ""),
		})
	}

	*parts = append(*parts, map[string]interface{}{
		"component": "renderTag",
		"tag":       "div",
		"selector":  "#mainform",
		"class":     "app-button-bar",
		"innerHTML": renderButtonBar(),
	})

	*parts = append(*parts, map[string]interface{}{
		"component": "eval",
		"script":    "checkForm();",
	})
}

func formScreenSave(parts *[]map[string]interface{}, request *GarageRequest, idReg string, saveReg *interface{}) {
	var err *string

	if idReg == "-1" {
		idReg = createUuid()
		err = insertDataObject(request.Db, request.User, map[string]interface{}{
			"id":         idReg,
			"objecttype": request.Screen.Options["dataobject"],
			"perms":      "{}",
			"reg":        saveReg,
		})
	} else {
		err = UpdateDataObject(request.Db, request.User, idReg, map[string]interface{}{
			"reg": saveReg,
		})
	}

	if err == nil {
		(*(*request).VarsUrl)["id"] = idReg
		*parts = append(*parts, map[string]interface{}{
			"component": "eval",
			"script":    "navigate('" + setVarsURL(request.Screen.Url, request.VarsUrl) + "?ok=Reg Saved.');",
		})
	} else {
		*parts = append(*parts, map[string]interface{}{
			"component":    "renderTag",
			"tag":          "app-msg",
			"msg":          err,
			"insertBefore": "form > app-input",
		})
	}

}

func formScreenCheck(parts *[]map[string]interface{}, request *GarageRequest, checkReg *interface{}) {

	for _, field := range request.Screen.Options["fields"].([]interface{}) {
		name := interfaceValueStr(&field, "name", "")
		visibleif, visibleifOk := field.(map[interface{}]interface{})["visibleif"]

		if visibleifOk {
			v := checkVisibleDisplay(&visibleif, checkReg)

			*parts = append(*parts, map[string]interface{}{
				"component": "eval",
				"script":    "document.querySelector('app-input[name=" + name + "]').style.display = '" + v + "';",
			})
		}

	}

}

func checkVisibleDisplay(visibleif *interface{}, checkReg *interface{}) string {
	field, fieldOk := (*visibleif).(map[interface{}]interface{})["field"]
	eq, eqOk := (*visibleif).(map[interface{}]interface{})["eq"]

	if fieldOk && eqOk {
		value, valueOk := (*checkReg).(map[string]interface{})[field.(string)]

		if valueOk && value == eq.(string) {
			return "block"
		}
	}

	return "none"
}

func renderHeader(parts *[]map[string]interface{}, request *GarageRequest) {
	actions := ""
	returnlink := ""
	toolbar, toolbarok := request.Screen.Options["toolbar"].([]interface{})
	if toolbarok {
		for _, action := range toolbar {
			actions += fmt.Sprint("<a class='action-top' href='", action.(map[interface{}]interface{})["url"], "'>", action.(map[interface{}]interface{})["label"], "</a>")
		}
	}

	returnurl, returnurlok := request.Screen.Options["returnurl"].(string)
	if returnurlok {
		returnlink = "<a class='returnlink' title='Volver' href='" + setVarsURL(returnurl, request.VarsUrl) + "'></a>"
	}

	*parts = append(*parts, map[string]interface{}{
		"component": "clean",
	})

	*parts = append(*parts, map[string]interface{}{
		"component": "renderTag",
		"tag":       "h2",
		"innerHTML": request.Screen.Options["title"].(string) + returnlink + actions,
	})

	okparam, okparamok := (*request.Params)["ok"]
	if okparamok {
		*parts = append(*parts, map[string]interface{}{
			"component": "renderTag",
			"tag":       "app-msg",
			"msg":       okparam,
			"class":     "app-msg-ok",
		})
	}
}

func interfaceValueStr(option *interface{}, field string, defaultValue string) string {
	of, ok := (*option).(map[interface{}]interface{})[field]

	if !ok {
		return defaultValue
	}

	return of.(string)
}

func DataObjectValueStr(do *DataObject, field string, defaultValue string) string {
	if do == nil {
		return defaultValue
	}

	of, ok := (*do.Reg)[field]
	if !ok {
		return defaultValue
	}

	return fmt.Sprint(of)
}

func renderButtonBar() string {
	r := ""

	r += "<button onclick='javascript: saveForm(event);'>Guardar</button>"

	return r
}
