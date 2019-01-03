build: 
	go build ./cmd/DrawOnGitHubHeatmap/

update-vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor