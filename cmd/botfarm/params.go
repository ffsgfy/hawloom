package main

const (
	baseURL     = "http://localhost:22440"
	botPassword = "qwerty"

	numClients = 30
	numDocs    = 10
	numVers    = 5

	// Stress test configurations:
	// numClients = 10, 50, 200, 500, 1500, 2300, 3000
	// numDocs    = 5,  20, 50,  200, 500,  700,  1000
	// numVers    = 5,  10, 15,  20,  30,   40,   50

	markovChainOrder  = 3
	docTitleLen       = 24
	docDescriptionLen = 96
	docContentLen     = 384
	verSummaryLen     = 48
	editMin           = 10
	editMax           = 100
	vordDurationMin   = 10
	vordDurationMax   = 30
	actionPeriod      = 1.0
)
