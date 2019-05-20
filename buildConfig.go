package main

import "os"

func buildConfiguration(fileLocation string) {
	configText := `#testingSettings
# Please do note remove the $ from the lines.
# $ is the help text in the menu.

[Mailbox Server Config]
# The target mailbox server where the http $server 
# requests are sent. $server $server
# Must be reachable from the host, and must $server
# have the mailbox service installed and running. $server
# $v_server=[ip],[hostname]
# $v_server_not=0.0.0.0
server=0.0.0.0

# The target mailbox server protocol where the $protocol
# http requests are sent. $protocol
# $v_protocol=http,https
protocol=https

# The target mailbox server port where the http $port
# requests are sent. $port
# $port
# You should only change this if you're running $port
# http or http on ports different than $port
# port 80 and port 443 respectively. $port
# $v_port=range[1-6553]
port=443

# The username the tests will use to authenticate. $username
# Doesn't need to be an administrator. $username
# $v_username=email
# $v_username_not=domain@example.com
username=domain@example.com

# The password the tests will use to authenticate. $password
# Doesn't need to be an administrator. $password
# $v_password_not=empty
password=1233456


[Threads Config]
# How many threads to create for testing. $threads
# These are not Zimbra threads, rather testing threads. $threads
# $threads
# For instance, if you set 10 threads, then it will $threads
# perform 10 tests at one time, continually until the $threads
# goal is met. $threads
# $v_threads_int
# $v_threads_not=empty
# $v_threads_warn=>5
threads=10

# Tests are performed in a loop.  $delayBetweenThreadRestart
# This is the delay between restarts for each loop/thread. $delayBetweenThreadRestart
# $delayBetweenThreadRestart
# 0 means no delay. Time is in milliseconds $delayBetweenThreadRestart
# $v_delayBetweenThreadRestart=int
# $v_delayBetweenThreadRestart_not=>100
delayBetweenThreadRestart=0

[Goals Config]
# The goal type is what decides when the test should end.
# First, you set the goal type. The for that type, you set the value.

# The type of goal to use. Goal types are: $goalType
# mdbfreepagesize,mdbfreepagecount,timer,events $goalType
# $v_goalType=mdbfreepagesize,mdbfreepagecount,timer,events
goalType=mdbfreepagesize

# The goal value to use. 
# $v_goalValue=int
goalValue=36000

[Environment]
# Free Page stats are read using the mdb_stat binary which $use_builtin_mdbstat
# is part of the Zimbra installation. $use_builtin_mdbstat
# $use_builtin_mdbstat
# Set this to true to use the tool which comes with this testing tool. $use_builtin_mdbstat
# $v_use_builtin_mdbstat=true,false
use_builtin_mdbstat=false

# mdb_stat is part of the Zimbra installation. $lmdb_stat_location
# $lmdb_stat_location
# In Zimbra 8.6, the location is: $lmdb_stat_location
#     /opt/zimbra/openldap-2.4.39.2z/bin/mdb_stat $lmdb_stat_location
# In Zimbra 8.7+ the location is: $lmdb_stat_location
#     /opt/zimbra/common/bin/mdb_stat $lmdb_stat_location
# $v_lmdb_stat_location=file
lmdb_stat_location=/opt/zimbra/openldap-2.4.39.2z/bin/mdb_stat

# The MDB database location is often referred to as the "environment". $lmdb_location
# This setting defines the location of the lmdb $lmdb_location
# $v_lmdb_location=filetype_mdb
lmdb_location=/opt/zimbra/data/ldap/mdb/db/

# You may wish to enable debug mode. This can print useful information. $debug_mode
# $v_debug_mode=true,false
debug_mode=false

# When the menus appear, the screen is cleared to make it more readable. $disable_clear_screen
# Set this to true to not clear the screen and preserve it. $disable_clear_screen
# $v_disable_clear_screen=true,false
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