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
}

type langfuseCtxValue string

const (
	langfuseCtxValuePath langfuseCtxValue = "langfusePath"
)

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

func (l *Langfuse) Trace(ctx context.Context, t *model.Trace) (context.Context, error) {
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
	ctx = context.WithValue(ctx, langfuseCtxValuePath, path.Path{path.Element{Type: path.Trace, ID: t.ID}})

	l.observer.Dispatch(event)

	return ctx, nil
}

//nolint:dupl
func (l *Langfuse) Generation(ctx context.Context, g *model.Generation) (context.Context, error) {
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
		return nil, err
	}
	g.TraceID = traceID

	event.Body = g

	if p, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		if p.Last().Type != path.Trace {
			g.ParentObservationID = p.Last().ID
		}
		p.Push(path.Generation, g.ID)
		ctx = context.WithValue(ctx, langfuseCtxValuePath, p)
	}

	l.observer.Dispatch(event)

	return ctx, nil
}

//nolint:dupl
func (l *Langfuse) GenerationEnd(ctx context.Context, g *model.Generation) (context.Context, error) {
	if p, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		e, pathErr := p.PopIf(path.Generation)
		if pathErr != nil {
			return nil, fmt.Errorf("invalid path: %w", pathErr)
		}
		ctx = context.WithValue(ctx, langfuseCtxValuePath, p)

		g.ID = e.ID
		g.TraceID = p.At(0).ID
	}

	if g.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeGenerationUpdate,
		Timestamp: time.Now().UTC(),
	}

	event.Body = g

	l.observer.Dispatch(event)

	return ctx, nil
}

func (l *Langfuse) Score(ctx context.Context, s *model.Score) (context.Context, error) {
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
		return nil, err
	}
	s.TraceID = traceID

	event.Body = s
	l.observer.Dispatch(event)

	return ctx, nil
}

//nolint:dupl
func (l *Langfuse) Span(ctx context.Context, s *model.Span) (context.Context, error) {
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
		return nil, err
	}
	s.TraceID = traceID

	event.Body = s
	if p, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		if p.Last().Type != path.Trace {
			s.ParentObservationID = p.Last().ID
		}
		p.Push(path.Span, s.ID)
		ctx = context.WithValue(ctx, langfuseCtxValuePath, p)
	}

	l.observer.Dispatch(event)

	return ctx, nil
}

//nolint:dupl
func (l *Langfuse) SpanEnd(ctx context.Context, s *model.Span) (context.Context, error) {
	if p, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		e, pathErr := p.PopIf(path.Span)
		if pathErr != nil {
			return nil, fmt.Errorf("invalid path: %w", pathErr)
		}
		ctx = context.WithValue(ctx, langfuseCtxValuePath, p)

		s.ID = e.ID
		s.TraceID = p.At(0).ID
	}

	if s.ID == "" {
		return nil, fmt.Errorf("span ID is required")
	}

	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	event := model.IngestionEvent{
		ID:        uuid.New().String(),
		Type:      model.IngestionEventTypeSpanUpdate,
		Timestamp: time.Now().UTC(),
	}

	event.Body = s

	l.observer.Dispatch(event)

	return ctx, nil
}

func (l *Langfuse) Event(ctx context.Context, e *model.Event) (context.Context, error) {
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
		return nil, err
	}
	e.TraceID = traceID

	if p, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		if p.Last().Type != path.Trace {
			e.ParentObservationID = p.Last().ID
		}
	}

	event.Body = e
	l.observer.Dispatch(event)

	return ctx, nil
}

func (l *Langfuse) extractTraceID(ctx context.Context, traceName string) (string, error) {
	extractedTraceID := ""
	if path, ok := ctx.Value(langfuseCtxValuePath).(path.Path); ok {
		extractedTraceID = path.At(0).ID
		return extractedTraceID, nil
	}

	ctxTrace, errTrace := l.Trace(
		ctx,
		&model.Trace{
			Name: traceName,
		},
	)
	if errTrace != nil {
		return "", errTrace
	}

	if path, ok := ctxTrace.Value(langfuseCtxValuePath).(path.Path); ok {
		extractedTraceID = path.At(0).ID
		return extractedTraceID, nil
	}

	return extractedTraceID, nil
}

func (l *Langfuse) Flush(ctx context.Context) {
	l.observer.Wait(ctx)
}
