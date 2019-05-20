package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	//"regexp"
	"strconv"

	//"fmt"

	//"fmt"
	_ "fmt"
	"log"
	"os"
	"strings"
	//ui "../ui"
)

func GetDescriptionTextForUpdate(configFile string, target string) (int,[]string) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	helpText := "Description for \""+strings.Replace(target,"_"," ",-1)+"\":\n"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "$"+target) {
			helpText += strings.Replace(strings.Replace(scanner.Text()+"\n", "$"+target, "", -1), "#", "", -1)
		}

	}
	helpByLines := strings.Split(helpText,"\n")
	helpByLines = delete_empty(helpByLines)

	return len(helpByLines), helpByLines
}

func delete_empty (s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func ValidateConfigKey(configFile string, keyName string, valueToCheck string) bool {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	isValid := true
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "$v_"+keyName+"_not") {
			//check for "not" values
			//rules := strings.ToLower(strings.Split(scanner.Text(),"=")[1])
			//if valueToCheck ==rules{
			//	return false
			//
			//}
			testVal := strings.Split(scanner.Text(),"=")[1]
			if valueToCheck==testVal {
				isValid=false
			}

		}else if strings.Contains(scanner.Text(), "["){
			continue
		}else{
			//check validity of values
			if strings.Contains(scanner.Text(), "$v_"+keyName) {
				if strings.Contains(scanner.Text(), "="){
					rules := strings.ToLower(strings.Split(scanner.Text(),"=")[1])
					//var TFMatch = regexp.MustCompile(`true|false`)
					if strings.Contains(rules, ","){
						//or matching
						validValues := strings.Split(rules, ",")
						isValid = false
						for _,v := range validValues{
							if valueToCheck==v{
								isValid = true
							}
						}
					}else if strings.Contains(rules, "int"){
						//interger matching
						isValid = false
						_, err := strconv.Atoi(valueToCheck)
						if err != nil {

						}else{
							isValid = true
						}
					}else if strings.Contains(rules, "email"){
						isValid = false
						if strings.Contains(valueToCheck, "@"){
							isValid=true
						}
					}
					//////////
				}

				//fmt.Println("\n\n\n")
				//fmt.Println(scanner.Text())
				//os.Exit(1)
			}
		}


	}
	//helpByLines := strings.Split(helpText,"\n")
	//helpByLines = delete_empty(helpByLines)
	//
	//return len(helpByLines), helpByLines
	return isValid
}

func UpdateConfigWithNewValue(config string, oldval string, newval string) {
	input, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte(oldval), []byte(newval), -1)

	if err = ioutil.WriteFile(config, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

var SectionMap map[string]string
var configSections map[string]string
var mapInit bool

func BuildMenuFromConfig(configFile string) map[string]string {
	if !mapInit {
		SectionMap = make(map[string]string)
		mapInit = true
	}
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	foundSection := false
	var tmpSectionName string
	var tmpSectionOptions string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#") {
			continue
		}
		if scanner.Text() == "" {
			continue
		}
		if foundSection {
			if strings.Contains(scanner.Text(), "=") {
				tmpSectionOptions += scanner.Text() + "||"
				//os.Exit(1)
			} else if strings.Contains(scanner.Text(), "[") {
				if strings.Contains(scanner.Text(), "]") {
					foundSection = false
					SectionMap[tmpSectionName] = tmpSectionOptions
					tmpSectionOptions = ""
					tmpSectionName = strings.Replace(strings.Replace(scanner.Text(), "[", "", -1), "]", "", -1)
					foundSection = true
				}
			}
		} else {

			if strings.Contains(scanner.Text(), "[") {
				if strings.Contains(scanner.Text(), "]") {
					tmpSectionName = strings.Replace(strings.Replace(scanner.Text(), "[", "", -1), "]", "", -1)
					foundSection = true
				}
			}
		}
	}

	return SectionMap
}