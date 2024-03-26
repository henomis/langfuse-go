package main

import (
	"context"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
)

func main() {
	ctx := context.Background()
	l := langfuse.New()

	ctx, err := l.Trace(
		ctx,
		&model.Trace{
			Name: "test-trace",
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Span(
		ctx,
		&model.Span{
			Name: "test-span",
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Generation(
		ctx,
		&model.Generation{
			Name:  "test-generation",
			Model: "gpt-3.5-turbo",
			ModelParameters: map[string]interface{}{
				"maxTokens":   "1000",
				"temperature": "0.9",
			},
			Input: []map[string]interface{}{
				{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				{
					"role":    "user",
					"content": "Please generate a summary of the following documents \nThe engineering department defined the following OKR goals...\nThe marketing department defined the following OKR goals...",
				},
			},
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Event(
		ctx,
		&model.Event{
			Name: "test-event",
			Metadata: map[string]interface{}{
				"key": "value",
			},
			Input: map[string]interface{}{
				"key": "value",
			},
			Output: map[string]interface{}{
				"key": "value",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.GenerationEnd(
		ctx,
		&model.Generation{
			Output: map[string]interface{}{
				"completion": "The Q3 OKRs contain goals for multiple teams...",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Score(
		ctx,
		&model.Score{
			Name:  "test-score",
			Value: 0.9,
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.SpanEnd(
		ctx,
		&model.Span{},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Span(
		ctx,
		&model.Span{
			Name: "test-span",
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Generation(
		ctx,
		&model.Generation{
			Name:  "test-generation",
			Model: "gpt-3.5-turbo",
			ModelParameters: map[string]interface{}{
				"maxTokens":   "1000",
				"temperature": "0.9",
			},
			Input: []map[string]interface{}{
				{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
				{
					"role":    "user",
					"content": "Please generate a summary of the following documents \nThe engineering department defined the following OKR goals...\nThe marketing department defined the following OKR goals...",
				},
			},
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.GenerationEnd(
		ctx,
		&model.Generation{
			Output: map[string]interface{}{
				"completion": "The Q3 OKRs contain goals for multiple teams...",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Score(
		ctx,
		&model.Score{
			Name:  "test-score",
			Value: 0.9,
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err = l.Event(
		ctx,
		&model.Event{
			Name: "test-event",
			Metadata: map[string]interface{}{
				"key": "value",
			},
			Input: map[string]interface{}{
				"key": "value",
			},
			Output: map[string]interface{}{
				"key": "value",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	_, err = l.SpanEnd(
		ctx,
		&model.Span{},
	)
	if err != nil {
		panic(err)
	}

	l.Flush(context.Background())

}
