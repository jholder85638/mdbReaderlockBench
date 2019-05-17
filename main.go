package main

import (
	"fmt"
	"github.com/bmatsuo/lmdb-go/lmdb"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	testsRun int
	mux sync.Mutex
}
var version string
var c SafeCounter
var dbi lmdb.DBI
var haveLoginData bool
var authMethod int
var csrfToken string
var adminAuthToken string
var preAuthToken string
var mdbPath string

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc() {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.testsRun += 1
	c.mux.Unlock()
}

func main(){
	version := "1.2-Alpha"
	initConfig(version)
	if !canContinue{
		setVariables()
		fmt.Println("The configuration safety check has failed. This can happen if you hit ctrl-c.")
		if debugEnabled{
			fmt.Println("Debug information will be printed. This may help find the issue. Cannot continue, sorry.")
			debug.PrintStack()
		}
		os.Exit(1)
	}
	setVariables()
	fmt.Println("======")
	fmt.Println("Welcome to the Zimbra LDAP Database Utility " + version + "!")
	fmt.Println("You should never run this on a production system.")
	fmt.Println("License: GPL - John Holder (john@johnholder.net)")
	fmt.Println("Disclaimer: Synacor & Zimbra are in no way responsible for this app.")
	fmt.Println("Your use of this program constitutes an agreement that you accept the above disclaimer.")
	fmt.Println("If you do not, you must exit now.")
	fmt.Println("======")
	defer os.Exit(1)

	env, err := lmdb.NewEnv()
	err = env.Open(mdbPath, 0, 0644)

	if err != nil {
		// ...
	}

	var readers []string
	var freePageCount string
	env.ReaderList(func(msg string) error {
		//log.Printf("reader: %q", msg)
		readers = append(readers, msg)
		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}


	//defer env.Close()
	numStale, err := env.ReaderCheck()
	if err != nil {
		log.Fatal(err.Error())
	}
	if numStale > 0 {
		log.Printf("Released locks for %d dead readers", numStale)
	}
	t1 := time.Now()
	maxLockTime :=65
	threads := 5
	threadCounter :=0
	finished := false
	log.Println("There are currently "+strconv.Itoa(len(readers))+" readers connected to the database.")
	log.Println("This program will place a lock on the database. While this lock in place by this app, it will also generate activity which writes to ldap.")
	log.Println("This should artificially increase the free page count for testing.")
	fmt.Println("Waiting 5 seconds before starting. If you wish to cancel, now is the time to hit [crtl]+[c]..")
	cancelBuffer := 0
	stat, _ := env.Stat()
	fmt.Println(stat)
	for{
		if cancelBuffer <5{
			fmt.Print(" .")
			cancelBuffer = cancelBuffer+1
			time.Sleep(1*time.Second)
		}else{
			fmt.Print(" .")
			fmt.Println("")
			log.Println("Starting...")
			break
		}
	}
	go func() {

		err = env.View(func(txn *lmdb.Txn) (err error) {
			cur, err := txn.OpenCursor(dbi)
			if err != nil {
				return err
			}

			for {
				_,_,err := cur.Get(nil, nil, lmdb.Next)
				if lmdb.IsNotFound(err) {
					return nil
				}
				if err != nil {
					return err
				}

				for{
					t2 := time.Now()
					delta := t2.Sub(t1).Seconds()

					if  int(delta)>maxLockTime {
						finished = true
						break
					}else{
						time.Sleep(1*time.Second)
					}

				}
			}

		})
	}()

	log.Print("A reader lock has been placed on the database. The lock will remain for "+strconv.Itoa(maxLockTime)+" seconds.")
	printedTestNotice := false
	testsStarted := false
	for{
		if finished{
			log.Println("\nClearing Reader Lock...")
			env.Close()
			os.Exit(0)
		}else{
			if !printedTestNotice{
				log.Println("Lock Detected. Starting Tests...")
				printedTestNotice = true
			}
			if !testsStarted{
				//threads := 10
				//threadCounter :=0
				fmt.Print("Starting thread: ")
				for{
					if threadCounter <= threads{
						threadCounter = threadCounter+1
						if threads ==threadCounter{
							fmt.Print(" "+ strconv.Itoa(threadCounter)+". Done.\n")
						}else{
							fmt.Print(" "+ strconv.Itoa(threadCounter)+", ")
						}
						go func() {
							LoginOperation()
						}()
					}else{
						break
					}
					time.Sleep(50*time.Millisecond)
				}
				testsStarted = true
			}
			//awk '{print $3 * 4096/1024/1024 " MB"}'
			FreePagesInt, err := strconv.Atoi(freePageCount)
			FreePagesInt = FreePagesInt * 4096/1024/1024
			//FreePagesMB := math.Round(FreePagesInt*100)/100)
			fmt.Print("\rEvents Run: "+strconv.Itoa(c.testsRun)+" / Current Free Page Count: "+freePageCount +" ("+strconv.Itoa(FreePagesInt)+" MB)")
			//./openldap-2.4.39.2z/bin/mdb_stat
			///opt/zimbra/common/bin/mdb_stat
			out, err := exec.Command("/opt/zimbra/openldap-2.4.39.2z/bin/mdb_stat","-a","-e","-f", "/opt/zimbra/data/ldap/mdb/db").Output()
			if err != nil {
				log.Fatal(err)
			}
			result := strings.Split(string(out),"\n")
			for _,v := range result{
				if strings.Contains(v, "Free pages"){
					freePageCount = strings.Split(v, ": ")[1]
				}
			}

			time.Sleep(1*time.Second)
		}

	}


}
