package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func randRange(valMin, valMax int) int {
	return valMin + rand.IntN(valMax-valMin+1)
}

func safeSlice[T any](s []T, low, high int) []T {
	low = min(max(low, 0), len(s))
	high = min(max(high, 0), len(s))
	return s[low:high]
}

func checkResponseOK(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("request status %d (%v): %s", resp.StatusCode, resp.Request.URL, string(body))
}

func readResponse(resp *http.Response, err error, th *timerHandle) ([]byte, error) {
	if th != nil {
		defer th.Stop()
	}
	if err != nil {
		return nil, fmt.Errorf("request error (%v): %w", resp.Request.URL, err)
	}
	defer resp.Body.Close()
	if err := checkResponseOK(resp); err != nil {
		return nil, err
	}
	if body, err := io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	} else {
		return body, nil
	}
}

func parseResponse(resp *http.Response, err error, th *timerHandle) (*goquery.Document, error) {
	body, err := readResponse(resp, err, th)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}
