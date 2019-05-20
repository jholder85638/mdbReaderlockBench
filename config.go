package main

import (
	zui "./ui"
	"bufio"
	"bytes"
	"fmt"
	_ "fmt"
	ui "github.com/VladimirMarkelov/clui"
	"github.com/gizak/termui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"os"
	"strings"
	utils "./utils"
)

func setVariables() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	c = SafeCounter{}
	c.testsRun = 0
	mdbPath = "/opt/zimbra/data/ldap/mdb/db/"
	haveLoginData = false
	debugEnabled = true
	disableClearScreen = false
	canContinue = false
	needBreak = false
	url = "https://192.168.1.17/"


}

func initConfig(version string) {
	//setVariables()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	ui.SetThemePath(dir + "/ui/themes/")
	//os.Exit(1)
	ui.SetCurrentTheme("turbovision")

	// These variables are needed for init.
	// Anything else should be in the setVariables function

	configFile = dir + "/settings.ini"
	if _, err := os.Stat(configFile); err == nil {

		// path/to/whatever exists
		//

	} else if os.IsNotExist(err) {
		if !disableClearScreen {
			print("\033[H\033[2J")
		}
		if err := termui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}

		//shouldExit := false
		//editConfig := false
		buttonsString := "Yes||No||Exit"
		buttons := strings.Split(buttonsString, "||")
		configAnswer := zui.CreateQuestionDialog("Zimbra LMDB Testing Utility", "A valid configuration could not be found.\nWould you like to create one now?", buttons, 2)
		ui.MainLoop()
		switch configAnswer.Result {
		case 1:
			ui.InitLibrary()
			ui.SetThemePath(dir + "/ui/themes/")
			ui.SetCurrentTheme("turbovision")
			buildConfiguration(configFile)
			configSections := BuildMenuFromConfig(configFile)
			zui.ConfigurationEditor(configSections, configFile)
			ui.MainLoop()
			os.Exit(1)
			break
		case 2:
			break
		case 3:
			ui.DeinitLibrary()
			//zui.Cleanup()
			exitReason = "You have chosen to exit without creating a configuration."
			return
			//break
		}


		//for {
		//	if dialogResult.View.Active(){
		//		ui.MainLoop()
		//		switch dialogResult.Result {
		//		case 1:
		//			//Yes
		//			dialogResult.View.SetActive(false)
		//			ui.WindowManager().DestroyWindow(dialogResult.View)
		//			ui.DeinitLibrary()
		//			zui.Cleanup()
		//			editConfig = true
		//			break
		//		case 2:
		//			shouldExit = true
		//			//No
		//			break
		//		case 3:
		//			shouldExit = true
		//			exitReason = "You have chosen to exit. Farewell!"
		//			//Exit
		//			break
		//		}
		//
		//	}else{
		//		if editConfig{
		//			ui.InitLibrary()
		//			ui.SetThemePath(dir + "/ui/themes/")
		//			fmt.Println(ui.ThemePath())
		//			//os.Exit(1)
		//			ui.SetCurrentTheme("turbovision")
		//			buildConfiguration(configFile)
		//			configSections = BuildMenuFromConfig(configFile)
		//			tmpMainMenuString :=""
		//			for k,_ := range configSections{
		//				tmpMainMenuString +=k+"||"
		//			}
		//			listBoxOptions := strings.Split(tmpMainMenuString, "||")
		//			result := zui.ConfigurationEditor(listBoxOptions)
		//			//ui.MainLoop()
		//			if result ==99{
		//				shouldExit=true
		//				break
		//			}
		//			//fmt.Println(result)
		//			//os.Exit(1)
		//			//switch result {
		//			//case 0:
		//			//
		//			//	break
		//			//case 1:
		//			//	shouldExit = true
		//			//	exitReason = "You have chosen to exit. Farewell!"
		//			//	//ui.WindowManager().DestroyWindow(dialogResult.View)
		//			//	ui.DeinitLibrary()
		//			//	zui.Cleanup()
		//			//	//os.Exit(1)
		//			//	//quit during config
		//			//	break
		//			//}
		//		}
		//
		//
		//	}
		//	if shouldExit{
		//		return
		//	}
		//}


		//go ui.Stop()
		//answerResult := result.Result

		//return
		//os.Exit(1)

		//options = []string{"Yes", "No"}
		//answer := AskQuestionMultiple("A valid configuration could not be found. Would you like to create one now?", options, "missingConfig")
		//location = "mainMenu"
		//lastLocation := ""
		//subsection := false
		//switch answer {
		//
		//case "Yes":
		//	buildConfiguration(configFile)
		//	if _, err := os.Stat(configFile); err == nil {
		//		for {
		//			configSections = BuildMenuFromConfig(configFile)
		//			tmpMainMenuString 	:= ""
		//			if location == "mainMenu"{
		//				for k,_ := range configSections{
		//					tmpMainMenuString +=k+"||"
		//				}
		//				tmpMainMenuString +="Start Tests Now ->||[Exit]"
		//				options = strings.Split(tmpMainMenuString, "||")
		//				location = AskQuestionMultiple("Configuration Options: ", options, location)
		//			}else if location == "<- Back to Main Menu"{
		//				if !disableClearScreen{
		//					print("\033[H\033[2J")
		//				}
		//				location = "mainMenu"
		//				lastLocation = ""
		//				subsection = false
		//			}else if location == "Start Tests Now ->"{
		//				if !disableClearScreen{
		//					print("\033[H\033[2J")
		//				}
		//				log.Println("Configuration saved, Starting Up...")
		//				break
		//			}else if location == "[Exit]"{
		//				if !disableClearScreen{
		//					print("\033[H\033[2J")
		//				}
		//				log.Println("Exiting. Farewell!")
		//				os.Exit(0)
		//			}else{
		//				if subsection==false {
		//					for k,v := range configSections{
		//						if k == location{
		//							v +="<- Back to Main Menu||Start Tests Now ->||[Exit]"
		//							options = strings.Split(v, "||")
		//							lastLocation = location
		//							location = AskQuestionMultiple("Editing "+k, options, "subSection")
		//							lastLocation += " => "+location
		//							subsection = true
		//						}else{
		//							break
		//						}
		//					}
		//				}else{
		//					if strings.Contains(location, ": "){
		//						location = AskQuestionMultiple(lastLocation, options, location)
		//						//reload the config
		//						configSections := BuildMenuFromConfig(configFile)
		//						tmpMainMenuString := ""
		//						for k,_ := range configSections{
		//							tmpMainMenuString +=k+"||"
		//						}
		//						tmpMainMenuString +="<- Back to Main Menu||Start Tests Now ->||[Exit]"
		//						options = strings.Split(tmpMainMenuString, "||")
		//						location = strings.Split(lastLocation, " => ")[0]
		//						subsection = false
		//					}else{
		//						break
		//					}
		//				}
		//			}
		//
		//		}
		//	} else if os.IsNotExist(err) {
		//		log.Fatal("Tried to write the configuration, but it failed for some reason: " + err.Error())
		//	}
		//	break
		//case "No":
		//	break
		//}
	}
}

func DisplayMenu(){

}
func AskQuestionMultiple(questionText string, options []string, location string) string {
	if !disableClearScreen {
		print("\033[H\033[2J")
	}
	answer := ""
	switch location {
	case "missingConfig":
		prompt := &survey.Select{
			Message:  questionText,
			Options:  options,
			PageSize: 10,
		}
		survey.AskOne(prompt, &answer, nil)
		break
	case "mainMenu":
		prompt := &survey.Select{
			Message:  questionText,
			Options:  options,
			PageSize: 10,
		}
		survey.AskOne(prompt, &answer, nil)
		break
	case "subSection":
		//split these options into editable components

		for k, v := range options {
			options[k] = strings.Replace(v, "=", ": ", -1)
		}
		prompt := &survey.Select{
			Message:  questionText,
			Options:  options,
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
			fmt.Println(utils.GetDescriptionTextForUpdate(configFile, tmpKey))

			promptQuestion := "Please enter a new value for " + tmpKey + " (was: " + prevVal + ")"
			prompt := &survey.Input{
				Message: promptQuestion + ": ",
			}
			survey.AskOne(prompt, &answer, nil)
			oldval := tmpKey + "=" + prevVal
			newval := tmpKey + "=" + answer
			fmt.Println("Updating Configuration: " + configFile)
			updateConfigWithNewValue(configFile, oldval, newval)
		}
	}

	return answer
}

func BuildMenuFromConfig(configFile string) map[string]string {

	SectionMap := make(map[string]string)
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



func isError(err error) bool {
	if err != nil {
		fmt.Println("FoundError:")
		fmt.Println(err.Error())
	}

	return (err != nil)
}
