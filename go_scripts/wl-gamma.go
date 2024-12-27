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

func changeTemp(temp string, m chan string){
    //tempString := "busctl --user set-property rs.wl-gammarelay / rs.wl.gammarelay Temperature q " + temp
    tempString := strings.Join([]string{"busctl --user set-property rs.wl-gammarelay / rs.wl.gammarelay Temperature q ", temp}, "")
    fmt.Println(tempString)
    execute(tempString, m)
}

func main() {
    messages := make(chan string)
    messagesChangeTemp := make(chan string)
    go getTemp(messages)
    //go messageHandler(messages)
    
    /* 
    for hour := time.Now().Format("3"); int(hour) > 20; hour = time.Now().Format("3") {
    for hour := time.Now().Format("3"); int(hour) > 20; hour = time.Now().Format("3") {
    for hour := time.Now().Format("3"); int(hour) > 20; hour = time.Now().Format("3") {
    for hour := time.Now().Format("3"); int(hour) > 20; hour = time.Now().Format("3") {
    for hour := time.Now().Format("3"); int(hour) > 20; hour = time.Now().Format("3") {
    */
    nightTemp := "3300"
    dayTemp := "6500"
    morningHour := 10
    nightHour := 21
    for {

        //currentTime := time.Now().Format(time.TimeOnly)
        currentTime := time.Now().Add(time.Hour * -4).Format(time.TimeOnly)

        splitTime := strings.Split(currentTime, ":")
        hourInt, err := strconv.Atoi(splitTime[0])
        if err != nil {
            panic(err)
        }       
        //messageHandler(messages) // check that the temp hasn't been user altered
        // messageHandler is a endless loop oops
        var temptype string
        temptype = calcTimeCheck(hourInt, morningHour, nightHour)
        if temptype == "night" || temptype == "morning" {
            changeTemp(nightTemp, messagesChangeTemp)
        } else if temptype == "day"{
            changeTemp(dayTemp, messagesChangeTemp)
        }
        fmt.Print(hourInt)
        SleepTillNextTarget(hourInt, morningHour, nightHour)
    }
}
func calcTimeCheck(hourInt int, morningHour int, nightHour int) (string) {
    var temptype string
    if hourInt < morningHour {
        fmt.Println("changing temp to morning temp")
        temptype = "morning"
    } else if hourInt >= nightHour {
        fmt.Println("changing temp to sleep temp")
        temptype = "night"
    } else {
        fmt.Println("changing temp to day time temp")
        temptype = "day"
    }
    return temptype
}

func SleepTillNextTarget(hourInt int, morningHour int, nightHour int){
    var sleepHours int
    if hourInt < morningHour {
        sleepHours = morningHour - hourInt
    } else if hourInt >= nightHour {
        sleepHours = hourInt - nightHour
    } else if hourInt < nightHour && hourInt >= morningHour {
        sleepHours = nightHour - hourInt
    } 
    fmt.Println("sleeping ", sleepHours, "hours")
    time.Sleep(60 * time.Hour * time.Duration(sleepHours)) // flawed since sleeping by hours and not a more procise merserment, leads to subpar reualts 

    currentTime := time.Now()//.Format(time.TimeOnly)
    targetTime := currentTime.Add(time.Hour * 5) // adds 5 hours to time
    formatedTime := targetTime.Format(time.TimeOnly)
    return formatedTime
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

