package calendar

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
)

type GoogleClient struct {
	service *calendar.Service
}

func NewGoogleClient(ctx context.Context) (*GoogleClient, error) {
	service, err := GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleClient{service: service}, nil
}

func (c *GoogleClient) FetchMeetings(calendarID string, start, end time.Time) ([]Event, error) {
	events, err := c.service.Events.List(calendarID).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	meetings := make([]Event, 0, len(events.Items))
	for _, event := range events.Items {
		// Skip all-day events (they don't have DateTime)
		if event.Start.DateTime == "" {
			continue
		}

		startTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			continue
		}

		endTime, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			continue
		}

		meeting := Event{
			Type:        EventTypeMeeting,
			Title:       event.Summary,
			Description: event.Description,
			Start:       startTime,
			End:         endTime,
		}

		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (c *GoogleClient) FetchTodaysMeetings(calendarID string) ([]Event, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return c.FetchMeetings(calendarID, startOfDay, endOfDay)
}

func (c *GoogleClient) CreateEvent(calendarID string, event Event) error {
	calEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Start: &calendar.EventDateTime{
			DateTime: event.Start.Format(time.RFC3339),
			TimeZone: event.Start.Location().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: event.End.Format(time.RFC3339),
			TimeZone: event.End.Location().String(),
		},
	}

	_, err := c.service.Events.Insert(calendarID, calEvent).Do()
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}
