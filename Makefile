PHONY: gomocks

gomocks:
	mockgen . Client > mocks/client.go