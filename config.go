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
var options []string
var configFile string
var location string
func initConfig(version string) {
	mapInit = false
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configFile = dir + "/settings.ini"
	if _, err := os.Stat(configFile); err == nil {

		// path/to/whatever exists
		//

	} else if os.IsNotExist(err) {
		print("\033[H\033[2J")
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
					//print("\033[H\033[2J")
					//Build Main menu
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
						print("\033[H\033[2J")
						location = "mainMenu"
						lastLocation = ""
						subsection = false
					}else if location == "Start Tests Now ->"{
						print("\033[H\033[2J")
						log.Println("Configuration saved, Starting Up...")
						break
					}else if location == "[Exit]"{
						print("\033[H\033[2J")
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
	//os.Exit(1)
}

func handleConfigAnswer(s string) {
	if strings.Contains(s, "[") {
		if strings.Contains(s, "]") {
			//this is a key.
		}
	}
}

func setVariables() {
	c = SafeCounter{}
	c.testsRun = 0
	mdbPath = "/opt/zimbra/data/ldap/mdb/db/"
	haveLoginData = false
}

func AskQuestionMultiple(questionText string, options []string, location string) string {
	print("\033[H\033[2J")
	print("Got: "+location+"\n\n")

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
			//we're editing a subsection
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

func buildConfiguration(fileLocation string) {
	configText := `#testingSettings
# Please do note remove the $ from the lines.
# $ is the help text in the menu.

[Mailbox Server Config]
# The target mailbox server where the http requests are sent. $server
# Must be reachable from the host, and must have the mailbox service installed and running. $server
server=192.168.1.17

# The target mailbox server protocol where the http requests are sent. $protocol
protocol=https

# The target mailbox server port where the http requests are sent. $port
# You should only change this if you're running http or http on ports different than $port
# port 80 and port 443 respectively. $port
port=443

# The username the tests will use to authenticate. Doesn't need to be an administrator. $username
username=john@johnholder.net

# The password the tests will use to authenticate. Doesn't need to be an administrator. $password
password=1233456


[Threads Config]
# How many threads to create for testing. $threads
# These are not Zimbra threads, rather testing threads. $threads
# For instance, if you set 10 threads, then it will perform $threads
#     10 tests at one time, continually until the goal is met. $threads
threads=10

# Tests are performed in a loop.  $delayBetweenThreadRestart
# This is the delay between restarts for each loop/thread. $delayBetweenThreadRestart
# 0 means no delay, this is in milliseconds $delayBetweenThreadRestart
delayBetweenThreadRestart=0

[Goals Config]
# The goal type is what decides when the test should end.
# First, you set the goal type. The for that type, you set the value.

# The type of goal to use. Goal types and values with examples can be viewed by typing ? and hitting enter. goalType
goalType=mdbfreepagesize

# The goal value to use.  Goal values and types with examples can be viewed by typing ? and hitting enter. goalValue
goalValue=36000

[Types]
# Types:
#   name: mdbfreepagesize
#   description: Converts the free page count, to MB
#   goal unit: number/MB
#   example:
#       type=mdbfreepagesize
#       goal=3000
#       ^ this would end the test when the free page size converts to 3000MB or higher.

#   name: mdbfreepagecount
#   description: Monitors the free page count
#   goal unit: count/number
#   example:
#       type=mdbfreepagescount
#       goal=64000
#       ^ this would end the test when the free page count reaches 64000 or higher

#   name: timer
#   description: Ends the test after a certain number of seconds
#   goal unit: seconds
#   example:
#       type=timer
#       goal=300
#       ^ this would end the test after 300 seconds have passed

#   name: events
#   description: Ends the test after a certain number events were performed
#   goal unit: count
#   example:
#       type=events
#       goal=300
#       ^ this would end the test after 300 events were performed.`

	var _, err = os.Create(fileLocation)
	if isError(err) {
		return
	}
	// open file using READ & WRITE permission
	var file, err2 = os.OpenFile(fileLocation, os.O_RDWR, 0644)
	if isError(err2) {
		return
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(configText)
	if isError(err) {
		return
	}

	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}

	//fmt.Println("==> done writing to file")
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
	//var options string
	//foundSectionOptions := false
	//var heading string
	//defaultsAdded := false
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
					//fmt.Println("Adding " + tmpSectionOptions + " to " + tmpSectionName)
					//fmt.Println(tmpSectionName)
					//fmt.Println(tmpSectionOptions)
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

		//	switch menuType {
		//
		//		print("\033[H\033[2J")
		//
		//		newval := ""
		//		fmt.Println(sectionName)
		//		thisKey := strings.Split(sectionName, "=")[0]
		//		thisVal := strings.Split(sectionName, "=")[1]
		//		fmt.Println("You are updating the value for " + thisKey+" which is currently set to: " + thisVal)
		//		descriptionText := getDescriptionTextForUpdate(configFile, thisKey)
		//		fmt.Println(descriptionText)
		//		prompt := &survey.Password{
		//			Message: thisKey + ": ",
		//		}
		//		survey.AskOne(prompt, &newval, nil)
		//		newval = thisKey+"="+newval
		//		sectionName = strings.Replace(sectionName, ": ","=",-1)
		//		fmt.Println(newval)
		//		fmt.Println(sectionName)
		//		fmt.Println(configFile)
		//		//os.Exit(1)
		//		updateConfigWithNewValue(configFile, sectionName, newval)
		//	}
		//}
		//if err := scanner.Err(); err != nil {
		//	log.Fatal(err)
		//}
		//if !defaultsAdded {
		//	options += "Edit config using your browser" + "||"
		//	options += "Exit and Don't Save" + "||"
		//	defaultsAdded = true
		//}
		//
		//parsedOptions := strings.Split(options, "||")
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
