#!/usr/bin/env sh

# FOR UNIT TEST ONLY
go test -short `go list ./... | grep -v tools` -coverprofile cover.out
go tool cover -func cover.out | grep total | awk '{print substr($3, 1, length($3)-1)}' | tee cover.score

# FOR INTEGRATION TESTS todo
go test -p 1 `go list ./... | grep -v tools` -coverprofile cover_integration.out
go tool cover -func cover_integration.out | grep total | awk '{print substr($3, 1, length($3)-1)}' | tee cover_integration.score
