package main

import (
	"bufio"
	"bytes"
	"fmt"
	_ "fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var sectionMap map[string]string
var mapInit bool
var canContinue bool
var needBreak bool
var options []string
var configFile string
var location string
var debugEnabled bool
var disableClearScreen bool


func setVariables() {
	c = SafeCounter{}
	c.testsRun = 0
	mdbPath = "/opt/zimbra/data/ldap/mdb/db/"
	haveLoginData = false
	debugEnabled = false
	disableClearScreen = false
	mapInit = false
	canContinue = false
	needBreak = false
}

func initConfig(version string) {
	setVariables()
	// These variables are needed for init.
	// Anything else should be in the setVariables function

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configFile = dir + "/settings.ini"
	if _, err := os.Stat(configFile); err == nil {

		// path/to/whatever exists
		//

	} else if os.IsNotExist(err) {
		if !disableClearScreen{
			print("\033[H\033[2J")
		}

		options = []string{"Yes", "No"}
		answer := AskQuestionMultiple("A valid configuration could not be found. Would you like to create one now?", options, "missingConfig")
		location = "mainMenu"
		lastLocation := ""
		subsection := false
		switch answer {

		case "Yes":
			buildConfiguration(configFile)
			if _, err := os.Stat(configFile); err == nil {
				for {
					configSections := buildMenuFromConfig(configFile)
					tmpMainMenuString 	:= ""
					if location == "mainMenu"{
						for k,_ := range configSections{
							tmpMainMenuString +=k+"||"
						}
						tmpMainMenuString +="Start Tests Now ->||[Exit]"
						options = strings.Split(tmpMainMenuString, "||")
						location = AskQuestionMultiple("Configuration Options: ", options, location)
					}else if location == "<- Back to Main Menu"{
						if !disableClearScreen{
							print("\033[H\033[2J")
						}
						location = "mainMenu"
						lastLocation = ""
						subsection = false
					}else if location == "Start Tests Now ->"{
						if !disableClearScreen{
							print("\033[H\033[2J")
						}
						log.Println("Configuration saved, Starting Up...")
						break
					}else if location == "[Exit]"{
						if !disableClearScreen{
							print("\033[H\033[2J")
						}
						log.Println("Exiting. Farewell!")
						os.Exit(0)
					}else{
						if subsection==false {
							for k,v := range configSections{
								if k == location{
									v +="<- Back to Main Menu||Start Tests Now ->||[Exit]"
									options = strings.Split(v, "||")
									lastLocation = location
									location = AskQuestionMultiple("Editing "+k, options, "subSection")
									lastLocation += " => "+location
									subsection = true
								}else{
									break
								}
							}
						}else{
							if strings.Contains(location, ": "){
								location = AskQuestionMultiple(lastLocation, options, location)
								//reload the config
								configSections := buildMenuFromConfig(configFile)
								tmpMainMenuString := ""
								for k,_ := range configSections{
									tmpMainMenuString +=k+"||"
								}
								tmpMainMenuString +="<- Back to Main Menu||Start Tests Now ->||[Exit]"
								options = strings.Split(tmpMainMenuString, "||")
								location = strings.Split(lastLocation, " => ")[0]
								subsection = false
							}else{
								break
							}
						}
					}

				}
			} else if os.IsNotExist(err) {
				log.Fatal("Tried to write the configuration, but it failed for some reason: " + err.Error())
			}
			break
		case "No":
			break
		}
	}
}



func AskQuestionMultiple(questionText string, options []string, location string) string {
	if !disableClearScreen{
		print("\033[H\033[2J")
	}
	answer := ""
	switch location {
	case "missingConfig":
		prompt := &survey.Select{
			Message: questionText,
			Options: options,
			PageSize: 10,
		}
		survey.AskOne(prompt, &answer, nil)
		break
	case "mainMenu":
		prompt := &survey.Select{
			Message: questionText,
			Options: options,
			PageSize: 10,
		}
		survey.AskOne(prompt, &answer, nil)
		break
	case "subSection":
		//split these options into editable components

		for k,v := range options{
			options[k] = strings.Replace(v, "=",": ",-1)
		}
		prompt := &survey.Select{
			Message: questionText,
			Options: options,
			PageSize: 10,
		}
		survey.AskOne(prompt, &answer, nil)
		break
	default:
		if strings.Contains(location, ": ") {
			tmpKey := strings.Split(location, ": ")[0]
			prevVal := strings.Split(location, ": ")[1]
			ActionDescriptionText := strings.Split(questionText, ": ")[0]
			fmt.Println("You are editing:")
			fmt.Println(ActionDescriptionText)
			fmt.Println(getDescriptionTextForUpdate(configFile, tmpKey))

			promptQuestion := "Please enter a new value for "+tmpKey+" (was: "+prevVal+")"
			prompt := &survey.Input{
				Message: promptQuestion + ": ",
			}
			survey.AskOne(prompt, &answer, nil)
			oldval := tmpKey+"="+prevVal
			newval := tmpKey+"="+answer
			fmt.Println("Updating Configuration: "+configFile)
			updateConfigWithNewValue(configFile, oldval, newval)

		}
	}
	return answer
}


func buildMenuFromConfig(configFile string) map[string]string {
	if !mapInit {
		sectionMap = make(map[string]string)
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
					sectionMap[tmpSectionName] = tmpSectionOptions
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
	return sectionMap
}

func updateConfigWithNewValue(config string, oldval string, newval string) {
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

func getDescriptionTextForUpdate(configFile string, target string) interface{} {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	helpText := "\nDescription:\n"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "$"+target) {
			helpText += strings.Replace(strings.Replace(scanner.Text()+"\n", "$"+target, "", -1), "#", "", -1)
		}

	}
	return helpText
}

func isError(err error) bool {
	if err != nil {
		fmt.Println("FoundError:")
		fmt.Println(err.Error())
	}

	return (err != nil)
}
