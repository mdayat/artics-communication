package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dbutil"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	logger := log.With().Caller().Logger()

	env, err := configs.LoadEnv()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	ctx := context.TODO()
	db, err := configs.NewDb(ctx, env.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	defer db.Conn.Close()

	configs := configs.NewConfigs(env, db)
	err = dbutil.RetryableTxWithoutData(ctx, configs.Db.Conn, configs.Db.Queries, func(qtx *repository.Queries) error {
		// Seed "user" table
		users := []repository.BulkInsertUserParams{
			{
				ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Email:    "john@gmail.com",
				Name:     "John",
				Password: "John",
				Role:     dtos.UserRole,
			},
			{
				ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Email:    "anne@gmail.com",
				Name:     "Anne",
				Password: "Anne",
				Role:     dtos.UserRole,
			},
		}

		_, err = qtx.BulkInsertUser(ctx, users)
		if err != nil {
			return fmt.Errorf("failed to bulk insert user: %w", err)
		}

		// Seed "meeting_room" table
		meetingRooms := []repository.BulkInsertMeetingRoomParams{
			{
				ID:   pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Name: "Meeting Room A",
			},
			{
				ID:   pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Name: "Meeting Room B",
			},
			{
				ID:   pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Name: "Meeting Room C",
			},
		}

		_, err = qtx.BulkInsertMeetingRoom(ctx, meetingRooms)
		if err != nil {
			return fmt.Errorf("failed to bulk insert meeting room: %w", err)
		}

		// Seed "time_slot" table
		now := time.Now()
		timeSlots := []repository.BulkInsertTimeSlotParams{
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[0].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()), Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[0].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location()), Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[1].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 10, 30, 0, 0, now.Location()), Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[1].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()), Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[2].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 0, 0, now.Location()), Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				MeetingRoomID: meetingRooms[2].ID,
				StartDate:     pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()), Valid: true},
				EndDate:       pgtype.Timestamptz{Time: time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location()), Valid: true},
			},
		}

		_, err = qtx.BulkInsertTimeSlot(ctx, timeSlots)
		if err != nil {
			return fmt.Errorf("failed to bulk insert time slot: %w", err)
		}

		// Seed "reservation" table
		reservations := []repository.BulkInsertReservationParams{
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				UserID:        users[0].ID,
				MeetingRoomID: meetingRooms[0].ID,
				TimeSlotID:    timeSlots[0].ID,
				Canceled:      true,
				CanceledAt:    pgtype.Timestamptz{Time: now, Valid: true},
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				UserID:        users[1].ID,
				MeetingRoomID: meetingRooms[1].ID,
				TimeSlotID:    timeSlots[3].ID,
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				UserID:        users[0].ID,
				MeetingRoomID: meetingRooms[2].ID,
				TimeSlotID:    timeSlots[4].ID,
			},
			{
				ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
				UserID:        users[1].ID,
				MeetingRoomID: meetingRooms[2].ID,
				TimeSlotID:    timeSlots[5].ID,
			},
		}

		_, err = qtx.BulkInsertReservation(ctx, reservations)
		if err != nil {
			return fmt.Errorf("failed to bulk insert reservation: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to seed database")
	}
}
