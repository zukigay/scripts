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
    for {
        hour := time.Now().Format("3")
        ampm := time.Now().Format("pm")
        hourInt, err := strconv.Atoi(hour)
        if err != nil {
            panic(err)
        }
        //messageHandler(messages) // check that the temp hasn't been user altered
        // messageHandler is a endless loop oops
        var temptype string
        if hourInt < 10 && ampm == "am" {
            fmt.Println("changing temp to morning temp")
            temptype = "night"
        } else if ampm == "am" && hourInt == 12 { // midnight
            fmt.Println("changing temp to sleep temp")
            temptype = "night"
        } else if ampm == "pm" && hourInt > 21 {
            fmt.Println("changing temp to sleep temp")
            temptype = "night"
        } else {
            fmt.Println("changing temp to day time temp")
            temptype = "day"
        }
        if temptype == "night"{
            changeTemp("3300", messagesChangeTemp)
        }
        fmt.Print(hourInt, ampm)

        //time.Sleep(60 * time.Second)
        time.Sleep(1 * time.Second)
    }
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

