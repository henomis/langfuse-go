package langfuse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/henomis/langfuse-go/internal/pkg/api"
	"github.com/henomis/langfuse-go/internal/pkg/observer"
	"github.com/henomis/langfuse-go/internal/pkg/path"
	"github.com/henomis/langfuse-go/model"
)

const (
	defaultFlushInterval = 500 * time.Millisecond
)

type Langfuse struct {
	flushInterval time.Duration
	client        *api.Client
	observer      *observer.Observer[model.IngestionEvent]
	path          *path.Path
}

func New() *Langfuse {
	client := api.New()

	l := &Langfuse{
		flushInterval: defaultFlushInterval,
		client:        client,
		observer: observer.NewObserver(
			func(events []model.IngestionEvent) {
				err := ingest(client, events)
				if err != nil {
					fmt.Println(err)
				}
			},
		),
		path: &path.Path{},
	}

	return l
}

func (l *Langfuse) WithFlushInterval(d time.Duration) *Langfuse {
	l.flushInterval = d
	return l
}

func ingest(client *api.Client, events []model.IngestionEvent) error {
	req := api.Ingestion{
		Batch: events,
	}

	res := api.IngestionResponse{}
	return client.Ingestion(context.Background(), &req, &res)
}

func (l *Langfuse) Trace(ctx context.Context, t *model.Trace) error {
	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeTraceCreate,
		Timestamp: time.Now().UTC(),
	}
	if t.ID == "" {
		traceID := uuid.New().String()
		t.ID = traceID
	}

	event.Body = t
	l.path = &path.Path{path.Element{Type: path.Trace, ID: t.ID}}

	l.observer.Dispatch(event)

	return nil
}

//nolint:dupl
func (l *Langfuse) Generation(ctx context.Context, g *model.Generation) error {
	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeGenerationCreate,
		Timestamp: time.Now().UTC(),
	}
	if g.ID == "" {
		traceID := uuid.New().String()
		g.ID = traceID
	}

	traceID, err := l.extractTraceID(ctx, g.Name)
	if err != nil {
		return err
	}
	g.TraceID = traceID

	if l.path.Last().Type != path.Trace {
		g.ParentObservationID = l.path.Last().ID
	}

	event.Body = g

	l.path.Push(path.Generation, g.ID)

	l.observer.Dispatch(event)

	return nil
}

//nolint:dupl
func (l *Langfuse) GenerationEnd(ctx context.Context, g *model.Generation) error {
	generation, err := l.path.PopIf(path.Generation)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	g.ID = generation.ID
	g.TraceID = l.path.At(0).ID

	if g.ID == "" {
		return fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return fmt.Errorf("trace ID is required")
	}

	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeGenerationUpdate,
		Timestamp: time.Now().UTC(),
	}

	event.Body = g

	l.observer.Dispatch(event)

	return nil
}

func (l *Langfuse) Score(ctx context.Context, s *model.Score) error {
	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeScoreCreate,
		Timestamp: time.Now().UTC(),
	}
	if s.ID == "" {
		traceID := uuid.New().String()
		s.ID = traceID
	}

	traceID, err := l.extractTraceID(ctx, s.Name)
	if err != nil {
		return err
	}
	s.TraceID = traceID

	event.Body = s
	l.observer.Dispatch(event)

	return nil
}

//nolint:dupl
func (l *Langfuse) Span(ctx context.Context, s *model.Span) error {
	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeSpanCreate,
		Timestamp: time.Now().UTC(),
	}
	if s.ID == "" {
		traceID := uuid.New().String()
		s.ID = traceID
	}

	traceID, err := l.extractTraceID(ctx, s.Name)
	if err != nil {
		return err
	}
	s.TraceID = traceID

	if l.path.Last().Type != path.Trace {
		s.ParentObservationID = l.path.Last().ID
	}

	l.path.Push(path.Span, s.ID)

	event.Body = s

	l.observer.Dispatch(event)

	return nil
}

//nolint:dupl
func (l *Langfuse) SpanEnd(ctx context.Context, s *model.Span) error {
	span, err := l.path.PopIf(path.Span)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	s.ID = span.ID
	s.TraceID = l.path.At(0).ID

	if s.ID == "" {
		return fmt.Errorf("span ID is required")
	}

	if s.TraceID == "" {
		return fmt.Errorf("trace ID is required")
	}

	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeSpanUpdate,
		Timestamp: time.Now().UTC(),
	}

	event.Body = s

	l.observer.Dispatch(event)

	return nil
}

func (l *Langfuse) Event(ctx context.Context, e *model.Event) error {
	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeEventCreate,
		Timestamp: time.Now().UTC(),
	}
	if e.ID == "" {
		traceID := uuid.New().String()
		e.ID = traceID
	}

	traceID, err := l.extractTraceID(ctx, e.Name)
	if err != nil {
		return err
	}
	e.TraceID = traceID

	if l.path.Last().Type != path.Trace {
		e.ParentObservationID = l.path.Last().ID
	}

	event.Body = e
	l.observer.Dispatch(event)

	return nil
}

func (l *Langfuse) extractTraceID(ctx context.Context, traceName string) (string, error) {
	tracePath := l.path.At(0)
	if tracePath != nil {
		return tracePath.ID, nil
	}

	errTrace := l.Trace(
		ctx,
		&model.Trace{
			Name: traceName,
		},
	)
	if errTrace != nil {
		return "", errTrace
	}

	tracePath = l.path.At(0)
	if tracePath != nil {
		return tracePath.ID, nil
	}

	return "", fmt.Errorf("unable to get trace ID")
}

func (l *Langfuse) Flush(ctx context.Context) {
	l.observer.Wait(ctx)
}
