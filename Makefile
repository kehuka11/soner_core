.PHONY: test-backend test-frontend test

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm run test:e2e
