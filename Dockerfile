FROM node:20-alpine AS frontend
WORKDIR /media_tracker
COPY public public
COPY templates templates
COPY package.json package-lock.json ./

RUN npm install
RUN npx tailwindcss -i ./public/tailwind.css -o ./public/dist/output.css

FROM golang:1.22.2-alpine AS backend
WORKDIR /media_tracker

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /media_tracker/public/dist ./public/dist
COPY --from=frontend /media_tracker/public/index.js ./public/index.js
COPY --from=frontend /media_tracker/public/tailwind.css ./public/tailwind.css


ENV CGO_ENABLED=1
RUN go build -o media_tracker ./cmd/media_tracker

ENTRYPOINT ["./media_tracker"]
CMD ["--mode=release"]
