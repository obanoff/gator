package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/obanoff/gator/internal/client"
	"github.com/obanoff/gator/internal/config"
	"github.com/obanoff/gator/internal/database"
)

var LoginHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 2 {
		return errors.New("username not provided")
	}

	_, err := s.DB.Queries.GetUserByName(context.Background(), args[1])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(args[1])
	if err != nil {
		return err
	}

	s.Logger.Success(fmt.Sprintf("username set to %s", args[1]))

	return nil
})

var RegisterHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 2 {
		return errors.New("name not provided")
	}

	user, err := s.DB.Queries.CreateUser(context.Background(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: args[1],
	})
	if err != nil {
		return err
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	s.Logger.Success(user)

	return nil
})

var UsersHandler = CommandHandler(func(s *config.State, args []string) error {
	users, err := s.DB.Queries.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		var str string
		if user.Name == s.Config.Username {
			str = fmt.Sprintf(" * %s (current)\n", user.Name)
		} else {
			str = fmt.Sprintf(" * %s\n", user.Name)
		}

		s.Logger.Success(str)
	}

	return nil
})

var ResetHandler = CommandHandler(func(s *config.State, args []string) error {
	err := s.DB.Queries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
})

var AddFeedHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 3 {
		return errors.New("not enough arguments provided")
	}

	tx, err := s.DB.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.DB.Queries.WithTx(tx)

	feed, err := qtx.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   args[1],
		Url:    args[2],
		Name_2: s.Config.Username,
	})
	if err != nil {
		return err
	}

	_, err = qtx.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		Name: s.Config.Username,
		FeedID: sql.NullInt32{
			Int32: feed.ID,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err == nil {
		s.Logger.Success(fmt.Sprintf("%-v", feed))
	}

	return err
})

var AggHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 2 {
		return errors.New("time not provided")
	}

	dur, err := time.ParseDuration(args[1])
	if err != nil {
		return errors.New("invalid time format")
	}

	s.Logger.Warning(fmt.Sprintf("Collecting feeds every %v", args[1]))

	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		client.ScrapeFeeds(s)
	}

	return nil
})

var FeedsHandler = CommandHandler(func(s *config.State, args []string) error {
	feeds, err := s.DB.Queries.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		s.Logger.Success(fmt.Sprintf("%-v", feed))
	}

	return nil
})

var FollowHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 2 {
		return errors.New("url not provided")
	}

	feed, err := s.DB.Queries.GetFeedByUrl(context.Background(), args[1])
	if err != nil {
		return err
	}

	feedFollow, err := s.DB.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		Name: s.Config.Username,
		FeedID: sql.NullInt32{
			Int32: feed.ID,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	s.Logger.Success(fmt.Sprintf("Feed: %s, User: %s\n", feedFollow.FeedName, feedFollow.UserName))

	return nil
})

var FollowingHandler = CommandHandler(func(s *config.State, args []string) error {
	rows, err := s.DB.Queries.GetFeedFollowsForUser(context.Background(), s.Config.Username)
	if err != nil {
		return err
	}

	for _, row := range rows {
		s.Logger.Success(row.FeedName)
	}

	return nil
})

var UnfollowHandler = CommandHandler(func(s *config.State, args []string) error {
	if len(args) < 2 {
		return errors.New("url not provided")
	}

	return s.DB.Queries.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		Name: s.Config.Username,
		Url:  args[1],
	})
})

var BrowserHandler = CommandHandler(func(s *config.State, args []string) error {
	var limit int
	var err error

	if len(args) < 2 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid limit format")
		}
	}

	posts, err := s.DB.Queries.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		Name:  s.Config.Username,
		Limit: int32(limit),
	})

	for _, p := range posts {
		fmt.Printf(`
%s

%s

%s
%v

-----------------------------------------------------------------------------------------
`,
			p.Title, p.Description, p.Url, p.PublishedAt.Time)
	}

	return nil
})
