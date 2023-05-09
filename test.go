package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func formatVertical(code string) string {
	// Replace all spaces with newline characters, including spaces inside strings
	insideString := false
	compactCode := ""
	for _, char := range code {
		if char == '"' || char == '\'' || char == '`' {
			insideString = !insideString
		}
		if insideString {
			if char == '\n' {
				compactCode += " "
			} else {
				compactCode += string(char)
			}
		} else {
			if char == ' ' {
				compactCode += "\n"
			} else {
				compactCode += string(char)
			}
		}
	}

	// Split the code into an array of tokens
	tokens := regexp.MustCompile(`([a-zA-Z_$][a-zA-Z_$0-9]*)|([\{\}\(\)\[\];,])|(\".*?\")|('.*?')|(`+"`"+`[\s\S]*?`+"`"+`)|(\d*\.\d+|\d+)|(\S)`).FindAllStringSubmatch(compactCode, -1)

	// Process each token
	resultingCode := ""
	for i, match := range tokens {
		token := match[0]
		if token == "{" {
			resultingCode += "\n" + token + "\n"
		} else if token == "}" {
			resultingCode += "\n" + token + "\n"
		} else if strings.HasSuffix(token, ";") {
			resultingCode += token + "\n"
		} else if token == "function" || token == "if" || token == "else" || token == "for" || token == "while" || token == "return" || token == "true" || token == "false" {
			resultingCode += "\n" + token + " "
		} else if strings.HasPrefix(token, "\"") || strings.HasPrefix(token, "'") || strings.HasPrefix(token, "`") {
			if strings.Contains(token, "\n") {
				resultingCode += "\n" + token + "\n"
			} else {
				resultingCode += " " + token + " "
			}
		} else if i+1 < len(tokens) && (token == "===" || token == "&&" || token == "||") && tokens[i+1][0] != ";" && tokens[i+1][0] != "," && tokens[i+1][0] != ")" {
			resultingCode += " " + token + " "
		} else {
			resultingCode += "\n" + token
		}
	}

	return resultingCode
}

func main() {
	// get jquery code as text from url
	resp, err := http.Get("https://code.jquery.com/jquery-3.6.4.min.js")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	inputCode := string(body)

	// remove comments from code
	inputCode = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(inputCode, "")

	formattedCode := formatVertical(inputCode)

	// save output to file
	err = ioutil.WriteFile("output.js", []byte(formattedCode), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

}
