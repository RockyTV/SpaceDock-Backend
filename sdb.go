/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package main

import (
    "os"
    "flag"
    "fmt"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    _ "github.com/KSP-SpaceDock/SpaceDock-Backend/routes"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/tools"
)

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the app package
*/
func main() {
    args := os.Args[1:]

    helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
    setupCommand := flag.NewFlagSet("setup", flag.ExitOnError)
    migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)

    // Setup subcommand flags
    dummyData := setupCommand.Bool("dummy", true, "Populates the database with dummy data")

    flag.Usage = func() {
        fmt.Printf("usage: sdb [command] [options]\n\n")
        fmt.Printf("SpaceDock backend application for handling database operations and http routes.\n\n")
        fmt.Printf("Use \"sdb help <command>\" for more information about a command.\n\n")
        fmt.Printf("    Commands:\n\n")
        fmt.Printf("        migrate     converts a pre-split SpaceDock database to the new backend database format\n")
        fmt.Printf("        setup       populates the database with dummy data and an administrator account\n\n")
        fmt.Printf("If no subcommand is specified, the backend application will run.\n")
    }

   if len(args) > 0 {
        switch args[0] {
        case "setup":
            setupCommand.Parse(args[1:])
        case "migrate":
            migrateCommand.Parse(args[1:])
        case "help":
            helpCommand.Parse(args[1:])
        default:
            flag.Usage()
            os.Exit(1)
        }
    } else {
        app.Run()
    }

    if helpCommand.Parsed() {
        helpArgs := helpCommand.Args()
        helpCommand.Usage = func() {
            // Default help text for help subcommand
            defaultUsage := func() {
                fmt.Printf("usage: sdb help <command>\n\n")
                fmt.Printf("    Commands:\n\n")
                fmt.Printf("        migrate     converts a pre-split SpaceDock database to the new backend database format\n")
                fmt.Printf("        setup       populates the database with dummy data and an administrator account\n\n")
            }

            // Check if we passed a valid subcommand as argument to the help command.
            // If we didn't, print the default help text.
            if len(helpArgs)  == 1 {
                switch helpArgs[0] {
                case "migrate":
                    fmt.Printf("usage: sdb migrate\n\n")
                    fmt.Printf("The migrate subcommand will convert an old database to the new backend format.\n")
                case "setup":
                    fmt.Printf("usage: sdb setup [-dummy=true|false]\n\n")
                    fmt.Printf("The setup subcommand will add an administrator account, a normal user,\n")
                    fmt.Printf("publisher, game, game admin and some dummy mods to the database.\n\n")
                    fmt.Printf("If you set the dummy flag to false, only an admin account will be added.\n")
                default:
                    defaultUsage()
                }
            } else if len(helpArgs) == 0 {
                defaultUsage()
            }
        }

        helpCommand.Usage()
        os.Exit(1)
    }

    if setupCommand.Parsed() {
        tools.Setup(*dummyData)
    }

    if migrateCommand.Parsed() {
        tools.MigrateDB()
    }
}

