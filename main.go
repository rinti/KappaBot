package main

import (
    "net"
    "net/textproto"
    "log"
    "bufio"
    "fmt"
    "strings"
    "io/ioutil"
    "encoding/json"
)

func main() {
 
    type Config struct {
        Name string `json:"name"`
        Network string `json:"network"`
        Port string `json:"port"`
        Channels []string `json:"channels"`
    }

    // Get config
    bs, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal("Could not open config")
    }

    // Load json config data
    var conf Config
    err = json.Unmarshal(bs, &conf)
    if err != nil {
        log.Fatal("Could not unmarshal JSON")
    }

    server := conf.Network + ":" + conf.Port
    fmt.Printf(server)
    connection, err := net.Dial("tcp", server)
    if err != nil {
        log.Fatal("Couldn't connect ", err)
    }
    connection.Write([]byte("NICK " + conf.Name + "\r\n"))
    connection.Write([]byte("USER kappot kappot kappot :The Kappa Bot\r\n"))
    defer connection.Close()

    reader := bufio.NewReader(connection)
    tp := textproto.NewReader(reader)

    for {
        line, err := tp.ReadLine()
        if err != nil {
            break   
        }

        fmt.Printf("<< %s\n", line)
    
        // We need to reply PONG to PINGs
        if(strings.Contains(line, "PING")) {
            pong_response := "PONG :" + strings.Split(line, ":")[1] + "\r\n"
            connection.Write([]byte(pong_response))
            fmt.Printf(">> %s\n", pong_response)
        }

        // If we're connected, join our channels
        if(strings.Contains(line, "MODE") && strings.Contains(line, conf.Name)) {
            for _, key := range conf.Channels {
                connection.Write([]byte("JOIN " + key + "\r\n"))
            }
        }
    }
}
