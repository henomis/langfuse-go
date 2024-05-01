package main

import (
	"context"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
)

func main() {
	l := langfuse.New(context.Background())

	trace, err := l.Trace(&model.Trace{Name: "test-trace"})
	if err != nil {
		panic(err)
	}

	span, err := l.Span(&model.Span{Name: "test-span", TraceID: trace.ID}, nil)
	if err != nil {
		panic(err)
	}

	generation, err := l.Generation(
		&model.Generation{
			TraceID: trace.ID,
			Name:    "test-generation",
			Model:   "gpt-3.5-turbo",
			ModelParameters: model.M{
				"maxTokens":   "1000",
				"temperature": "0.9",
			},
			Input: []model.M{
				{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				{
					"role":    "user",
					"content": "Please generate a summary of the following documents \nThe engineering department defined the following OKR goals...\nThe marketing department defined the following OKR goals...",
				},
			},
			Metadata: model.M{
				"key": "value",
			},
		},
		&span.ID,
	)
	if err != nil {
		panic(err)
	}

	_, err = l.Event(
		&model.Event{
			Name:    "test-event",
			TraceID: trace.ID,
			Metadata: model.M{
				"key": "value",
			},
			Input: model.M{
				"key": "value",
			},
			Output: model.M{
				"key": "value",
			},
		},
		&generation.ID,
	)
	if err != nil {
		panic(err)
	}

	generation.Output = model.M{
		"completion": "The Q3 OKRs contain goals for multiple teams...",
	}
	_, err = l.GenerationEnd(generation)
	if err != nil {
		panic(err)
	}

	_, err = l.Score(
		&model.Score{
			TraceID: trace.ID,
			Name:    "test-score",
			Value:   0.9,
		},
	)
	if err != nil {
		panic(err)
	}

	_, err = l.SpanEnd(span)
	if err != nil {
		panic(err)
	}

	l.Flush(context.Background())

}
