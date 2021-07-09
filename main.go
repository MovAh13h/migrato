package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)

	upDsn := upCmd.String("dsn", "", "dsn")
	downDsn := downCmd.String("dsn", "", "dsn")

	runDsn := runCmd.String("dsn", "", "dsn")
	runFile := runCmd.String("file", "", "file")

	newName := newCmd.String("name", "", "name")

	if len(os.Args) < 2 {
		panic("Expected subcommand")
	}

	switch os.Args[1] {
	case "up", "down":
		var dsn string
		var opType string

		if os.Args[1] == "up" {
			upCmd.Parse(os.Args[2:])
			dsn = *upDsn
			opType = "up"
		} else if os.Args[1] == "down" {
			downCmd.Parse(os.Args[2:])
			dsn = *downDsn
			opType = "down"
		}

		conn, err := sql.Open("mysql", dsn)

		if err != nil {
			panic(err)
		}

		defer conn.Close()

		tx, err := conn.Begin()

		if err != nil {
			panic(err)
		}

		defer tx.Rollback()

		files, err := ioutil.ReadDir(fmt.Sprintf("./migrations/%s", opType))

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if file.IsDir() == true {
				continue
			}

			filePath := fmt.Sprintf("./migrations/%s/%s", opType, file.Name())

			queryBytes, err := ioutil.ReadFile(filePath)

			if err != nil {
				tx.Rollback()
				panic(err)
			}

			query := string(queryBytes)

			_, err = tx.Exec(query)

			if err != nil {
				tx.Rollback()
				panic(err)
			}

			fmt.Printf("[OK] %s\n", file.Name())
		}

		tx.Commit()

		break

	case "run":
		runCmd.Parse(os.Args[2:])

		conn, err := sql.Open("mysql", *runDsn)

		if err != nil {
			panic(err)
		}

		defer conn.Close()

		queryBytes, err := ioutil.ReadFile(*runFile)

		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(string(queryBytes))

		if err != nil {
			panic(err)
		}

		fmt.Printf("[OK] %s", *runFile)

		break

	case "new":
		newCmd.Parse(os.Args[2:])

		currTime := time.Now().Unix()

		upFilePath := fmt.Sprintf("./migrations/up/%d_%s.up.sql", currTime, *newName)
		downFilePath := fmt.Sprintf("./migrations/down/%d_%s.down.sql", currTime, *newName)

		upFile, err := os.Create(upFilePath)

		if err != nil {
			panic(err)
		}

		defer upFile.Close()

		downFile, err := os.Create(downFilePath)

		if err != nil {
			panic(err)
		}

		defer downFile.Close()

		fmt.Printf("[OK] Created %s\n", upFilePath)
		fmt.Printf("[OK] Created %s\n", downFilePath)

		break

	case "init":
		err := os.MkdirAll("./migrations/up", 0755)

		if err != nil {
			panic(err)
		}

		err = os.MkdirAll("./migrations/down", 0755)

		if err != nil {
			panic(err)
		}

		break

	default:
		panic("No valid subcommand found")
	}
}
