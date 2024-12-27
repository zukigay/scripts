package main

import (
	"bufio"
    "os"
	"fmt"
    "io"
	"os/exec"
	"errors"
	"strings"
    "time"
    "strconv"
)

func changeTemp(temp string,tempCmd string, m chan string, verbose bool){
    //tempCmd := "busctl --user set-property rs.wl-gammarelay / rs.wl.gammarelay Temperature q "
    tempString := strings.Join([]string{tempCmd, temp}, "")
    //tempString := strings.Join([]string{"busctl --user set-property rs.wl-gammarelay / rs.wl.gammarelay Temperature q ", temp}, "")
    
    if verbose == true{
        fmt.Println(tempString)
    }
    execute(tempString, m)
}
func parseStrToInt(String string) (int) {
    Int, err := strconv.Atoi(String)
    if err != nil {
        panic("error: while parsing " + String + " into int.") 
    }
    return Int
}
//               nt      dt      mh   nh   mm   nm   verbose tempCmd
func argParse() (string, string, int, int, int, int, bool, string){
    //layer := "null" 

    // create argVars and set defualt values
    nightTemp := "3300"
    dayTemp := "6500"
    morningHour := 10
    nightHour := 21
    morningMin := 00
    nightMin := 00
    printArgs := false
    tempCmd := "busctl --user set-property rs.wl-gammarelay / rs.wl.gammarelay Temperature q "
    //getTempCmd := // not yet needed
    verbose := false 

    for i := 1; i < len(os.Args); i++ {
        arg := os.Args[i]
        //if layer == "null" {
            switch arg {
            case "-v":
                verbose = true
            case "-p":
                // toggle printing out args
                printArgs = true
            case "-nt":
                // set night temp(string)
                i++
                nightTemp = os.Args[i]

            case "-dt":
                // set day temp(string)
                i++
                dayTemp = os.Args[i]

            case "-mh":
                i++
                morningHour = parseStrToInt(os.Args[i])

            case "-nh":
                i++
                nightHour = parseStrToInt(os.Args[i])

            case "-mm":
                i++
                morningMin = parseStrToInt(os.Args[i])
            case "-nm":
                i++
                nightMin = parseStrToInt(os.Args[i])

            case "-h", "--help":
                fmt.Print("options:\n",
                            "   -h/--help print this text\n",
                            "   -v enable verbose logging\n",
                            "   -nt set night temp\n",
                            "   -dt set day temp\n",
                            "   -mh set moring hour\n",
                            "   -nh set night hour\n",
                            "   -nm set night min\n",
                            "   -mm set morning min\n",
                            "example\n",
                            "   autogamma -nt 3300 -dt 6500 -mh 10 -nh 21\n")
                os.Exit(0)

            /*
            case "-n":
                // set hours(two strings) example input "10,21" 
                i++
                parseHours := strings.Split(os.Args[i], ",")
                morningHour, err := strconv.Atoi(parseHours[0])
                if err != nil {
                    panic("error: ", err, " while tring to convert -n into morningHour") 
                }
                nightHour, err := strconv.Atoi(parseHours[1])
                if err != nil {
                    panic("error: ", err, " while tring to convert -n into nightHour") 
                } */

            default:
                panic("unknown option '" + arg + "'. exiting...")


            }
            if printArgs == true{
                fmt.Println(arg)
            }
        /* } else{
            if printArgs == true {
                fmt.Print(arg)
            }

        } */

    }
    return nightTemp, dayTemp, morningHour, nightHour, morningMin, nightMin, verbose, tempCmd

}

func main() {
    nightTemp, dayTemp, morningHour, nightHour, morningMin, nightMin, verbose, tempCmd := argParse()
    //messages := make(chan string)
    messagesChangeTemp := make(chan string)
    //go getTemp(messages)
    //go messageHandler(messages)
    

    for {

        currentTime := time.Now() //.Add(time.Hour * 26) debug setting NOTE some function's notably SleepTillNextTarget still works off of the system time
        hourInt := currentTime.Hour()
        //messageHandler(messages) // check that the temp hasn't been user altered
        // messageHandler is a endless loop oops
        var temptype string
        temptype = calcTimeCheck(hourInt, morningHour, nightHour, verbose)
        if temptype == "night" || temptype == "morning" {
            changeTemp(nightTemp,tempCmd, messagesChangeTemp, verbose)
        } else if temptype == "day"{
            changeTemp(dayTemp, tempCmd, messagesChangeTemp, verbose)
        }
        if verbose == true{
            fmt.Println("HourInt:", hourInt)
        }
        SleepTillNextTarget(hourInt, morningHour, nightHour, currentTime, morningMin, nightMin, verbose)
    }
}
func calcTimeCheck(hourInt int, morningHour int, nightHour int, verbose bool ) (string) {
    var temptype string
    var verboseTempChangeString string
    if hourInt < morningHour {
        verboseTempChangeString = "changing temp to morning temp"
        temptype = "morning"
    } else if hourInt >= nightHour {
        verboseTempChangeString = "changing temp to sleep temp"
        temptype = "night"
    } else {
        verboseTempChangeString = "changing temp to day time temp"
        temptype = "day"
    }
    if verbose == true{
        fmt.Println(verboseTempChangeString)
    }
    return temptype
}

func SleepTillNextTarget(hourInt int, morningHour int, nightHour int,cTime time.Time, morningMin int, nightMin int, verbose bool){
    // cTime is currentTime
    //cTime := time.Now()//.Add(time.Hour * 5)
    var targetTime time.Time
    if hourInt < morningHour {
        targetTime = getTargetTime(cTime, morningHour, morningMin, false)
    } else if hourInt >= nightHour {
        targetTime = getTargetTime(cTime, morningHour, morningMin, true)
    } else if hourInt < nightHour && hourInt >= morningHour {
        targetTime = getTargetTime(cTime, nightHour, nightMin, false)
    } 

    if verbose == true{
        fmt.Println(cTime, "\n", targetTime)
    }
    time.Sleep(time.Until(targetTime))
}

func getTargetTime(cTime time.Time, targetHour int, targetMin int, addDay bool) (time.Time){
    var targetTime time.Time
    if addDay == true{ 
    //                      year          month          day          hour                   min         sec nanosec timezone
    targetTime = time.Date(cTime.Year(), cTime.Month(), cTime.Day() + 1, targetHour, targetMin, 00, 0, cTime.Location())
    } else if addDay == false{ 
    //                      year          month          day                    hour         min sec nanosec timezone
    targetTime = time.Date(cTime.Year(), cTime.Month(), cTime.Day(), targetHour, targetMin, 00, 0, cTime.Location())
    }
    return targetTime
}


func messageHandler(messages chan string){
    for message := range messages {
        fmt.Println("temp", message)
        /*
        messageInt, err := strconv.Atoi(message)
        if err != nil {
            panic(err)
        } */
        //if message != "6500\n" ||  message != "3300\n\n" {
        if  message == "6500\n" || message == "3300\n" {
            fmt.Println("vaid temp")
        }else {
            fmt.Println("temp tampered with exiting...")
            //os.Exit(0)
        }
    }
}

func getTemp(m chan string) {
	execute("wl-gammarelay-rs watch {t}", m)	
}

func sendMessage(message chan string, sentMessage string) {
    message <- sentMessage 
}

func execute(cmd string, m chan string) (err error) {
	if cmd == "" {
		return errors.New("No command provided")
	}

	cmdArr := strings.Split(cmd, " ")
	name := cmdArr[0]

	args := []string{}
	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}

	command := exec.Command(name, args...)
	command.Env = os.Environ()

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()
	stdoutReader := bufio.NewReader(stdout)

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()
	stderrReader := bufio.NewReader(stderr)

	if err := command.Start(); err != nil {
		return err
	}

	go handleReader(stdoutReader, m)
	go handleReader(stderrReader, m)

	if err := command.Wait(); err != nil {
		return err
	}

    return err
}

func handleReader(reader *bufio.Reader, m chan string) error {
	for {

		str, err := reader.ReadString('\n')
		if len(str) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

        //fmt.Print(str)
        m <- str

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

