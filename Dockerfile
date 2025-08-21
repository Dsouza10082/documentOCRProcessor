# Use golang and the build env
FROM golang:1.24.4 AS build

WORKDIR /go/src

ARG TARGETOS
ARG TARGETARCH

ENV GOOS=${TARGETOS:-linux}
ENV GOARCH=${TARGETARCH:-arm64}

# Install required packages first
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/go/bin:${PATH}"


COPY . .

COPY go.mod               ./go/src
COPY go.sum               ./go/src
COPY main.go              ./go/src

RUN go get
RUN CGO_ENABLED=0 go build -o ocr-processor

FROM golang:1.24.4-alpine

WORKDIR /go/bin

RUN apk update && apk add --no-cache \
    python3 \
    py3-pip \
    tesseract-ocr \
    tesseract-ocr-data-eng \
    tesseract-ocr-data-por \
    tesseract-ocr-data-spa \
    tesseract-ocr-data-fra \
    tesseract-ocr-data-deu \
    gcc \
    g++ \
    musl-dev \
    python3-dev \
    libffi-dev \
    jpeg-dev \
    zlib-dev \
    freetype-dev \
    lcms2-dev \
    openjpeg-dev \
    tiff-dev \
    tk-dev \
    tcl-dev \
    postgresql-dev \
    postgresql-client \
    && rm -rf /var/cache/apk/*

ENV TESSDATA_PREFIX=/usr/share/tessdata
ENV PYTHONUNBUFFERED=1
    
# Create virtual environment and install Python packages
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
    PyPDF2 \
    pdfplumber \
    pytesseract \
    pillow \
    pgai \
    psycopg2-binary \
    asyncpg

COPY --from=build /go/src/ocr-processor .

RUN adduser -D ocr-processor && \
# Add permissions for the user to access the files
    chown -R ocr-processor /go/bin/ocr-processor && \
    mkdir -p /home/ocr-processor/plugin && chown -R 799:799 /home/ocr-processor/plugin && \
    mkdir -p /home/ocr-processor/learning && chown -R 799:799 /home/ocr-processor/learning && \
    chmod +x /go/bin/ocr-processor

RUN mkdir /files

COPY ./learning/.  /home/ocr-processor/learning

# Set the entrypoint
ENTRYPOINT ["/go/bin/ocr-processor"]