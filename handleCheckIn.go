package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/techplexengineer/go-frc-attend/data"
	"log"
	"net/http"
	"time"
)

const (
	YMDDateFormat = "2006-01-02"
)

func (s Server) handleCheckIn() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			s.handleInternalError(err)(w, req)
			return
		}
		userID := req.FormValue("userid")

		user, err := s.db.GetUser(context.Background(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// user does not exist
				SetFlash(w, "User does not exist")
				http.Redirect(w, req, "/create/"+userID, http.StatusSeeOther)
				return
			}
			log.Print(fmt.Errorf("error getting user %s: %w", userID, err))
			s.handleInternalError(err)(w, req)
			return
		}

		count, err := s.db.IsUserCheckedIn(context.Background(), data.IsUserCheckedInParams{
			Userid: userID,
			Date:   time.Now().Format(YMDDateFormat),
		})
		if err != nil {
			err = fmt.Errorf("error IsUserCheckedIn %s: %w", userID, err)
			s.handleInternalError(err)(w, req)
			return
		}
		if count > 0 {
			SetFlash(w, fmt.Sprintf("Success %s %s already checked in", user.FirstName, user.LastName))
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		err = s.db.CheckinUser(context.Background(), data.CheckinUserParams{
			Userid: userID,
			Date:   time.Now().Format(YMDDateFormat),
		})
		if err != nil {
			err = fmt.Errorf("error CheckinUser %s: %w", userID, err)
			s.handleInternalError(err)(w, req)
			return
		}

		SetFlash(w, fmt.Sprintf("Success %s %s checked in", user.FirstName, user.LastName))
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
