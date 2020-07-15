package rediscli_test

import (
	"testing"

	"github.com/luddd3/rediscli"
	"github.com/stretchr/testify/assert"
)

func beforeEach() *rediscli.Client {
	client, err := rediscli.New("localhost:6379")
	client.Debug = false

	if err != nil {
		panic(err)
	}

	err = client.Auth("redispass")
	if err != nil {
		panic(err)
	}
	return client
}

func TestPing(t *testing.T) {
	client := beforeEach()
	results, err := client.Ping()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", results)
}

func TestGetMiss(t *testing.T) {
	client := beforeEach()
	results, err := client.Get("unknown")
	assert.Nil(t, err)
	assert.Equal(t, nil, results)
}
