package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/obanoff/gator/internal/commands"
	"github.com/obanoff/gator/internal/config"
	"github.com/obanoff/gator/internal/database"
	"github.com/obanoff/gator/internal/middleware"

	_ "github.com/lib/pq"
)

func main() {
	state := &config.State{
		Logger: config.NewLogger(log.LstdFlags),
	}

	cfg, err := config.Read()
	if err != nil {
		state.Logger.Error(err)
		os.Exit(1)
	}

	state.Config = &cfg

	db, err := sql.Open("postgres", state.Config.DBUrl)
	if err != nil {
		state.Logger.Error(err)
		os.Exit(1)
	}

	state.DB = config.Database{
		Queries: database.New(db),
		DB:      db,
	}

	cmds := commands.NewCommands()
	cmds.Register("login", commands.LoginHandler)
	cmds.Register("register", commands.RegisterHandler)
	cmds.Register("reset", commands.ResetHandler)
	cmds.Register("users", commands.UsersHandler)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(commands.AddFeedHandler))
	cmds.Register("agg", commands.AggHandler)
	cmds.Register("feeds", commands.FeedsHandler)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(commands.FollowHandler))
	cmds.Register("following", middleware.MiddlewareLoggedIn(commands.FollowingHandler))
	cmds.Register("unfollow", middleware.MiddlewareLoggedIn(commands.UnfollowHandler))
	cmds.Register("browse", middleware.MiddlewareLoggedIn(commands.BrowserHandler))

	err = cmds.Run(state, os.Args[1:])
	if err != nil {
		state.Logger.Error(err)
		os.Exit(1)
	}

}
