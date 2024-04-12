package langfuse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/henomis/langfuse-go/internal/pkg/api"
	"github.com/henomis/langfuse-go/internal/pkg/observer"
	"github.com/henomis/langfuse-go/model"
)

const (
	defaultFlushInterval = 500 * time.Millisecond
)

type Langfuse struct {
	flushInterval time.Duration
	client        *api.Client
	observer      *observer.Observer[model.IngestionEvent]
}

func New(ctx context.Context) *Langfuse {
	client := api.New()

	l := &Langfuse{
		flushInterval: defaultFlushInterval,
		client:        client,
		observer: observer.NewObserver(
			ctx,
			func(ctx context.Context, events []model.IngestionEvent) {
				err := ingest(ctx, client, events)
				if err != nil {
					fmt.Println(err)
				}
			},
		),
	}

	return l
}

func (l *Langfuse) WithFlushInterval(d time.Duration) *Langfuse {
	l.flushInterval = d
	return l
}

func ingest(ctx context.Context, client *api.Client, events []model.IngestionEvent) error {
	req := api.Ingestion{
		Batch: events,
	}

	res := api.IngestionResponse{}
	return client.Ingestion(ctx, &req, &res)
}

func (l *Langfuse) Trace(t *model.Trace) (*model.Trace, error) {
	t.ID = buildID(&t.ID)
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeTraceCreate,
			Timestamp: time.Now().UTC(),
			Body:      t,
		},
	)
	return t, nil
}

func (l *Langfuse) Generation(g *model.Generation, parentID *string) (*model.Generation, error) {
	if g.TraceID == "" {
		traceID, err := l.createTrace(g.Name)
		if err != nil {
			return nil, err
		}

		g.TraceID = traceID
	}

	g.ID = buildID(&g.ID)

	if parentID != nil {
		g.ParentObservationID = *parentID
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeGenerationCreate,
			Timestamp: time.Now().UTC(),
			Body:      g,
		},
	)
	return g, nil
}

func (l *Langfuse) GenerationEnd(g *model.Generation) (*model.Generation, error) {
	if g.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeGenerationUpdate,
			Timestamp: time.Now().UTC(),
			Body:      g,
		},
	)

	return g, nil
}

func (l *Langfuse) Score(s *model.Score) (*model.Score, error) {
	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}
	s.ID = buildID(&s.ID)

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeScoreCreate,
			Timestamp: time.Now().UTC(),
			Body:      s,
		},
	)
	return s, nil
}

func (l *Langfuse) Span(s *model.Span, parentID *string) (*model.Span, error) {
	if s.TraceID == "" {
		traceID, err := l.createTrace(s.Name)
		if err != nil {
			return nil, err
		}

		s.TraceID = traceID
	}

	s.ID = buildID(&s.ID)

	if parentID != nil {
		s.ParentObservationID = *parentID
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeSpanCreate,
			Timestamp: time.Now().UTC(),
			Body:      s,
		},
	)

	return s, nil
}

func (l *Langfuse) SpanEnd(s *model.Span) (*model.Span, error) {
	if s.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeSpanUpdate,
			Timestamp: time.Now().UTC(),
			Body:      s,
		},
	)

	return s, nil
}

func (l *Langfuse) Event(e *model.Event, parentID *string) (*model.Event, error) {
	if e.TraceID == "" {
		traceID, err := l.createTrace(e.Name)
		if err != nil {
			return nil, err
		}

		e.TraceID = traceID
	}

	e.ID = buildID(&e.ID)

	if parentID != nil {
		e.ParentObservationID = *parentID
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        uuid.New().String(),
			Type:      model.IngestionEventTypeEventCreate,
			Timestamp: time.Now().UTC(),
			Body:      e,
		},
	)

	return e, nil
}

func (l *Langfuse) createTrace(traceName string) (string, error) {
	trace, errTrace := l.Trace(
		&model.Trace{
			Name: traceName,
		},
	)
	if errTrace != nil {
		return "", errTrace
	}

	return trace.ID, fmt.Errorf("unable to get trace ID")
}

func (l *Langfuse) Flush(ctx context.Context) {
	l.observer.Wait(ctx)
}

func buildID(id *string) string {
	if id == nil {
		return uuid.New().String()
	} else if *id == "" {
		return uuid.New().String()
	}

	return *id
}
