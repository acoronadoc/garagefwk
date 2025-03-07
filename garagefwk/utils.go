package garagefwk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getVarsURL(url1 string, url2 string) *map[string]interface{} {
	r := map[string]interface{}{}

	purl1 := strings.Split(url1, "?")
	purl2 := strings.Split(url2, "?")

	purl1 = strings.Split(purl1[0], "/")
	purl2 = strings.Split(purl2[0], "/")

	if len(purl1) != len(purl2) {
		return nil
	}

	for i, _ := range purl1 {
		if strings.HasPrefix(purl1[i], "${") && strings.HasSuffix(purl1[i], "}") {
			r[purl1[i][2:len(purl1[i])-1]] = purl2[i]
		} else {
			if purl1[i] != purl2[i] {
				return nil
			}
		}
	}

	return &r
}

func setVarsURL(url string, varsURL *map[string]interface{}) string {
	r := ""
	purl := strings.Split(url, "?")

	parts := strings.Split(purl[0], "/")
	for i, _ := range parts {
		if i != 0 {
			r += "/"
		}

		if strings.HasPrefix(parts[i], "${") && strings.HasSuffix(parts[i], "}") {
			r += (*varsURL)[parts[i][2:len(parts[i])-1]].(string)
		} else {
			r += parts[i]
		}
	}

	return r
}

func ReadForm(req *http.Request) *map[string]interface{} {
	var data map[string]interface{}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	if len(body) == 0 {
		return nil
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		panic(err)
	}

	return &data
}

func ReadParams(req *http.Request) *map[string]string {
	data := map[string]string{}

	params, _ := url.ParseQuery(req.URL.RawQuery)
	url := params.Get("url")
	pos := strings.Index(url, "?")
	if pos != -1 {
		substr := string(url[pos+1 : len(url)])
		substrparts := strings.Split(substr, "&")

		for _, v := range substrparts {
			vp := strings.Split(v, "=")
			data[vp[0]] = vp[1]
		}
	}

	return &data
}
