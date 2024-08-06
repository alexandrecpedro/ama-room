package main

import (
	"fmt"
	// "os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	// (1) Load .env variables
	if err := godotenv.Load(); err != nil {
		// panic(err)
		panic(fmt.Sprintf("Error while loading .env: %v", err))
	}

	// (2) Verify the env variables
	// fmt.Println("WSRS_DATABASE_HOST:", os.Getenv("WSRS_DATABASE_HOST"))
	// fmt.Println("WSRS_DATABASE_PORT:", os.Getenv("WSRS_DATABASE_PORT"))
	// fmt.Println("WSRS_DATABASE_NAME:", os.Getenv("WSRS_DATABASE_NAME"))
	// fmt.Println("WSRS_DATABASE_USER:", os.Getenv("WSRS_DATABASE_USER"))
	// fmt.Println("WSRS_DATABASE_PASSWORD:", os.Getenv("WSRS_DATABASE_PASSWORD"))

	// (3) Set and execute the following command: tern migrate
	// Native package for Go = os/exec
	// exec.Command(Command_to_be_executed, args)
	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)
	// if err := cmd.Run(); err != nil {
	// 	panic(err)
	// }

	// (4) Capture and print the command output for diagnostics
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erro ao executar tern migrate: %v\nSaída: %s\n", err, output)
		panic(err)
	}

	// fmt.Printf("Saída: %s\n", output)
}
