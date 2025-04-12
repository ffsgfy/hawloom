package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	username string
	client   *http.Client
	gen      *TextGenerator
	timer    Timer
}

func NewClient(username string, gen *TextGenerator) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		username: username,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Jar:       jar,
		},
		gen:   gen,
		timer: NewTimer(),
	}
}

func (c *Client) Auth() error {
	vals := url.Values{}
	vals.Set("name", c.username)
	vals.Set("password", botPassword)
	vals.Set("password-re", botPassword)

	th := c.timer.Start("register")
	resp, err := c.client.PostForm(baseURL+"/auth/register", vals)
	if _, err = readResponse(resp, err, th); err != nil {
		return err
	}

	th = c.timer.Start("login")
	resp, err = c.client.PostForm(baseURL+"/auth/login", vals)
	if _, err = readResponse(resp, err, th); err != nil {
		return err
	}

	return nil
}

func findDocID(doc *goquery.Document) (string, error) {
	if docID, ok := doc.Find("#doc-id").Attr("value"); ok {
		return docID, nil
	} else {
		return "", fmt.Errorf("failed to find doc id")
	}
}

func (c *Client) NewDoc() (string, error) {
	title := strings.TrimSpace(c.gen.Reset().GenerateN(docTitleLen, true))
	description := strings.TrimSpace(c.gen.Reset().GenerateN(docDescriptionLen, true))
	content := strings.TrimSpace(c.gen.Reset().GenerateN(docContentLen, false))
	vordDuration := randRange(vordDurationMin, vordDurationMax)

	vals := url.Values{}
	vals.Set("title", title)
	vals.Set("description", description)
	vals.Set("content", content)
	vals.Set("vduration", fmt.Sprint(vordDuration))
	vals.Set("vmode", "selection")
	vals.Set("public", "true")
	vals.Set("majority", "false")

	th := c.timer.Start("new_doc")
	resp, err := c.client.PostForm(baseURL+"/doc/new", vals)
	doc, err := parseResponse(resp, err, th)
	if err != nil {
		return "", err
	}

	if docID, err := findDocID(doc); err == nil {
		return "/doc/" + docID, nil
	} else {
		return "", err
	}
}

func (c *Client) GetDoc(docPath string) (string, string, error) {
	th := c.timer.Start("get_doc")
	resp, err := c.client.Get(baseURL + docPath)
	doc, err := parseResponse(resp, err, th)
	if err != nil {
		return "", "", err
	}

	docID, err := findDocID(doc)
	if err != nil {
		return "", "", err
	}

	return docID, doc.Find("#txt-current").Text(), nil
}

func (c *Client) GetDocList(username string) ([]string, error) {
	th := c.timer.Start("get_doc_list")
	resp, err := c.client.Get(baseURL + "/user/" + username)
	doc, err := parseResponse(resp, err, th)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, a := range doc.Find("table > tbody > tr > td > a").EachIter() {
		if path, ok := a.Attr("href"); ok {
			paths = append(paths, path)
		} else {
			return nil, fmt.Errorf("failed to find doc path")
		}
	}

	return paths, nil
}

func editContent(content string, gen *TextGenerator) string {
	toDelete := randRange(editMin, editMax)
	toInsert := randRange(editMin, editMax)

	tokens := strings.Split(content, "")
	index := rand.IntN(len(tokens))
	prefix := safeSlice(tokens, 0, index)
	suffix := safeSlice(tokens, index+toDelete, len(tokens))
	gen.Seed(prefix)

	generated := make([]string, toInsert)
	for i := range generated {
		generated[i] = gen.Generate()
	}

	return strings.Join(slices.Concat(prefix, generated, suffix), "")
}

func (c *Client) NewVer(docID, content string) error {
	content = editContent(content, c.gen)
	summary := strings.TrimSpace(c.gen.Reset().GenerateN(verSummaryLen, true))

	vals := url.Values{}
	vals.Set("doc-id", docID)
	vals.Set("content", content)
	vals.Set("summary", summary)

	th := c.timer.Start("new_ver")
	resp, err := c.client.PostForm(baseURL+"/ver/new", vals)
	_, err = parseResponse(resp, err, th)

	return err
}

func (c *Client) DeleteVer(verPath string) error {
	req, err := http.NewRequest("DELETE", baseURL+verPath, nil)
	if err != nil {
		return err
	}
	th := c.timer.Start("delete_ver")
	resp, err := c.client.Do(req)
	_, err = readResponse(resp, err, th)

	return err
}

type verListRow struct {
	path    string
	author  string
	hasVote bool
}

func (c *Client) GetVerList(docID string) ([]verListRow, error) {
	th := c.timer.Start("get_ver_list")
	resp, err := c.client.Get(baseURL + "/ver/list?vord-num=-1&doc-id=" + docID)
	doc, err := parseResponse(resp, err, th)
	if err != nil {
		return nil, err
	}

	rows := []verListRow{}
	for _, tr := range doc.Find("table > tbody > tr").EachIter() {
		td := tr.Find("td")
		if td.Length() == 1 {
			continue // No data
		} else if td.Length() != 4 {
			return nil, fmt.Errorf("unexpected number of fields in ver list row: %d", td.Length())
		}

		path, ok := td.Eq(0).Find("span").Attr("hx-get")
		if !ok {
			return nil, fmt.Errorf("failed to find ver path")
		}
		author := td.Eq(1).Text()
		hasVote := td.Eq(3).Find("span").HasClass("underline")

		rows = append(rows, verListRow{
			path:    path,
			author:  author,
			hasVote: hasVote,
		})
	}

	return rows, nil
}

func (c *Client) VerVote(verPath string) error {
	th := c.timer.Start("ver_vote")
	resp, err := c.client.PostForm(baseURL+verPath+"/vote", nil)
	_, err = readResponse(resp, err, th)
	return err
}

func (c *Client) VerUnvote(verPath string) error {
	th := c.timer.Start("ver_unvote")
	resp, err := c.client.PostForm(baseURL+verPath+"/unvote", nil)
	_, err = readResponse(resp, err, th)
	return err
}
