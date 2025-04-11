package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func run(client *Client, docs []string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for {
		delay := time.Millisecond * time.Duration(max(rand.Float64()+actionPeriod-0.5, 0)*1000)
		ctxlog.Info(ctx, "sleeping", "delay", delay)
		time.Sleep(delay)
		if ctx.Err() != nil {
			ctxlog.Info(ctx, "exiting", "err", ctx.Err())
			break
		}

		docPath := docs[rand.IntN(len(docs))]
		docID, content, err := client.GetDoc(docPath)
		if err != nil {
			ctxlog.Error2(ctx, "failed to get doc", err, "path", docPath)
			continue
		}
		ctxlog.Info(ctx, "got doc", "doc_id", docID)

		vers, err := client.GetVerList(docID)
		if err != nil {
			ctxlog.Error2(ctx, "failed to get ver list", err, "doc_id", docID)
			continue
		}
		ctxlog.Info(ctx, "got ver list", "count", len(vers))

		actions := []func() error{}
		canDelete := []string{}
		canUnvote := []string{}
		canVote := []string{}

		for i := range vers {
			verRow := &vers[i]
			if verRow.author == client.username {
				canDelete = append(canDelete, verRow.path)
			}
			if verRow.hasVote {
				canUnvote = append(canUnvote, verRow.path)
			} else {
				canVote = append(canVote, verRow.path)
			}
		}

		if len(vers) < numVers {
			actions = append(actions, func() error {
				ctxlog.Info(ctx, "action: new ver")
				return client.NewVer(docID, content)
			})
		}

		if len(canDelete) > 0 {
			actions = append(actions, func() error {
				path := canDelete[rand.IntN(len(canDelete))]
				ctxlog.Info(ctx, "action: delete ver", "path", path)
				return client.DeleteVer(path)
			})
		}

		if len(canUnvote) > 0 {
			actions = append(actions, func() error {
				path := canUnvote[rand.IntN(len(canUnvote))]
				ctxlog.Info(ctx, "action: unvote ver", "path", path)
				return client.VerUnvote(path)
			})
		} else if len(canVote) > 0 {
			actions = append(actions, func() error {
				path := canVote[rand.IntN(len(canVote))]
				ctxlog.Info(ctx, "action: vote ver", "path", path)
				return client.VerVote(path)
			})
		}

		if len(actions) > 0 {
			if err := actions[rand.IntN(len(actions))](); err != nil {
				ctxlog.Error2(ctx, "action failed", err)
			}
		} else {
			ctxlog.Warn(ctx, "no actions")
		}
	}
}

type clientPayload struct {
	err    error
	index  int
	client *Client
}

func initClient(index int, gen *TextGenerator, out chan<- *clientPayload) {
	payload := clientPayload{}
	username := fmt.Sprintf("bot-%d", index)
	client := NewClient(username, gen.Clone())
	if err := client.Auth(); err != nil {
		payload.err = err
	} else {
		payload.index = index
		payload.client = client
	}
	out <- &payload
}

func main() {
	ctxlog.SetDefault(ctxlog.New(os.Stdout, ctxlog.INFO))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	sourcePath := flag.String("source", "", "markov chain source text path")
	flag.Parse()

	source, err := os.ReadFile(*sourcePath)
	if err != nil {
		panic(err)
	}

	chain := NewMarkovChain(string(source), markovChainOrder)
	gen := NewTextGenerator(chain)

	clientsOut := make(chan *clientPayload)
	clients := []*Client{}
	client0 := (*Client)(nil)

	for i := range numClients {
		go initClient(i, gen, clientsOut)
	}
	for range numClients {
		payload := <-clientsOut
		if payload == nil {
			panic("client payload is nil")
		}
		if payload.err != nil {
			panic(payload.err)
		} else {
			ctxlog.Info(ctx, "client initialized", "bot", payload.client.username)
			clients = append(clients, payload.client)
			if payload.index == 0 {
				client0 = payload.client
			}
		}
	}
	close(clientsOut)

	if len(clients) < numClients || client0 == nil {
		panic("not all clients initialized")
	}

	docs, err := client0.GetDocList(client0.username)
	if err != nil {
		panic(err)
	}
	for len(docs) < numDocs {
		if path, err := client0.NewDoc(); err == nil {
			docs = append(docs, path)
			ctxlog.Info(ctx, "doc created", "path", path)
		} else {
			panic(err)
		}
	}
	docs = safeSlice(docs, 0, numDocs)

	wg := sync.WaitGroup{}
	for _, client := range clients {
		wg.Add(1)
		go run(client, docs, &wg, ctxlog.With(ctx, "bot", client.username))
	}

	<-ctx.Done()
	wg.Wait()
}
