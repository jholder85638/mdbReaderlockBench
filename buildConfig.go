package main

import "os"

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

# The type of goal to use. Goal types and values with examples can be viewed by typing ? and hitting enter. $goalType
goalType=mdbfreepagesize

# The goal value to use.  Goal values and types with examples can be viewed by typing ? and hitting enter. $goalValue
goalValue=36000

[Environment]
# Free Page stats are read using the mdb_stat binary which is part of the Zimbra installation. $use_builtin_mdbstat
# Set this to true to use the tool which comes with this testing tool. $use_builtin_mdbstat
use_builtin_mdbstat=false

# mdb_stat is part of the Zimbra installation. This sets the location. $lmdb_stat_location
# In Zimbra 8.6, the location is: $lmdb_stat_location
#     /opt/zimbra/openldap-2.4.39.2z/bin/mdb_stat $lmdb_stat_location
# In Zimbra 8.7+ the location is: $lmdb_stat_location
#     /opt/zimbra/common/bin/mdb_stat $lmdb_stat_location
lmdb_stat_location=/opt/zimbra/openldap-2.4.39.2z/bin/mdb_stat

# The MDB database location is often referred to as the "environment". $lmdb_location
# This setting defines the location of the lmdb $lmdb_location
lmdb_location=/opt/zimbra/data/ldap/mdb/db/

# You may wish to enable debug mode. This can print useful information. $debug_mode
debug_mode=false

# When the menus appear, the screen is cleared to make it more readable. $disable_clear_screen
# Set this to true to not clear the screen and preserve it. $disable_clear_screen
disable_clear_screen=false

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