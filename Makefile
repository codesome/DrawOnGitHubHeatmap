build: 
	go build ./cmd/name_in_git_heatmap/

update-vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor