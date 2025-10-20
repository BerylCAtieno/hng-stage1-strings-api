# String Analyzer API

A RESTful API service that analyzes strings and stores their computed properties. Built with Go and SQLite.

## Features

- **String Analysis**: Compute multiple properties for any string including length, palindrome check, character frequency, word count, and SHA-256 hash
- **RESTful Endpoints**: Full CRUD operations for analyzed strings
- **Advanced Filtering**: Query strings using structured filters or natural language
- **Persistent Storage**: SQLite database for reliable data persistence
- **Error Handling**: Comprehensive error responses with appropriate HTTP status codes

## Table of Contents

- [Installation](#installation)
- [Setup](#setup)
- [Running Locally](#running-locally)
- [API Endpoints](#api-endpoints)
- [Query Parameters](#query-parameters)
- [Examples](#examples)
- [Dependencies](#dependencies)
- [Environment Variables](#environment-variables)
- [Testing](#testing)
- [Project Structure](#project-structure)

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Clone Repository

```bash
git clone https://github.com/BerylCAtieno/string-analyzer.git
cd string-analyzer
```

## Setup

### Install Dependencies

```bash
go mod download
go mod tidy
```

This will install all required dependencies specified in `go.mod`, including:
- `modernc.org/sqlite` - SQLite driver for Go

### Initialize Database

The database is automatically initialized when the application starts. It creates a `./data/strings.db` SQLite database file.

## Running Locally

### Start the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Verify Server is Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## API Endpoints

### 1. Create/Analyze String

**Endpoint**: `POST /strings`

**Request**:
```json
{
  "value": "hello world"
}
```

**Response** (201 Created):
```json
{
  "id": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
  "value": "hello world",
  "properties": {
    "length": 11,
    "is_palindrome": false,
    "unique_characters": 8,
    "word_count": 2,
    "sha256_hash": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
    "character_frequency_map": {
      "h": 1,
      "e": 1,
      "l": 3,
      "o": 2,
      "w": 1,
      "r": 1,
      "d": 1,
      " ": 1
    }
  },
  "created_at": "2025-08-27T10:00:00Z"
}
```

**Error Responses**:
- `400 Bad Request` - Missing or empty "value" field
- `409 Conflict` - String already exists in the system
- `422 Unprocessable Entity` - Invalid data type (value must be string)

---

### 2. Get Specific String

**Endpoint**: `GET /strings/{string_value}`

**Example**: `GET /strings/hello%20world`

**Response** (200 OK):
```json
{
  "id": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
  "value": "hello world",
  "properties": { /* ... */ },
  "created_at": "2025-08-27T10:00:00Z"
}
```

**Error Responses**:
- `404 Not Found` - String does not exist in the system

---

### 3. Get All Strings with Filtering

**Endpoint**: `GET /strings`

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `is_palindrome` | boolean | Filter for palindromic strings (true/false) |
| `min_length` | integer | Minimum string length (inclusive) |
| `max_length` | integer | Maximum string length (inclusive) |
| `word_count` | integer | Exact number of words |
| `contains_character` | string | Single character to search for (case-insensitive) |

**Example**: `GET /strings?is_palindrome=true&min_length=5&max_length=20`

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": "abc123...",
      "value": "racecar",
      "properties": { /* ... */ },
      "created_at": "2025-08-27T10:00:00Z"
    }
  ],
  "count": 1,
  "filters_applied": {
    "is_palindrome": true,
    "min_length": 5,
    "max_length": 20
  }
}
```

**Error Responses**:
- `400 Bad Request` - Invalid query parameter values or types
- `400 Bad Request` - min_length greater than max_length

---

### 4. Natural Language Filtering

**Endpoint**: `GET /strings/filter-by-natural-language`

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `query` | string | Natural language query to filter strings |

**Supported Query Patterns**:

| Query | Converts To |
|-------|------------|
| "all single word palindromic strings" | `word_count=1, is_palindrome=true` |
| "strings longer than 10 characters" | `min_length=11` |
| "strings shorter than 5 characters" | `max_length=4` |
| "strings at least 8 characters" | `min_length=8` |
| "strings at most 20 characters" | `max_length=20` |
| "exactly 15 character strings" | `min_length=15, max_length=15` |
| "two word palindromes" | `word_count=2, is_palindrome=true` |
| "strings containing the letter z" | `contains_character=z` |
| "palindromic strings with the first vowel" | `is_palindrome=true, contains_character=a` |

**Example**: `GET /strings/filter-by-natural-language?query=single%20word%20palindromes`

**Response** (200 OK):
```json
{
  "data": [ /* matching strings */ ],
  "count": 3,
  "interpreted_query": {
    "original": "single word palindromes",
    "parsed_filters": {
      "word_count": 1,
      "is_palindrome": true
    }
  }
}
```

**Error Responses**:
- `400 Bad Request` - Missing or empty query parameter
- `400 Bad Request` - Unable to parse natural language query

---

### 5. Delete String

**Endpoint**: `DELETE /strings/{string_value}`

**Example**: `DELETE /strings/hello%20world`

**Response** (204 No Content): Empty response body

**Error Responses**:
- `404 Not Found` - String does not exist in the system

---

## Query Parameters

### is_palindrome
Filters strings that read the same forwards and backwards (case-insensitive).

```bash
curl "http://localhost:8080/strings?is_palindrome=true"
```

### min_length / max_length
Filters strings by character length. Can be used together or separately.

```bash
curl "http://localhost:8080/strings?min_length=5&max_length=20"
```

### word_count
Filters strings with an exact number of whitespace-separated words.

```bash
curl "http://localhost:8080/strings?word_count=2"
```

### contains_character
Filters strings containing a specific character (case-insensitive).

```bash
curl "http://localhost:8080/strings?contains_character=a"
```

### Combining Filters
All filters can be combined in a single request:

```bash
curl "http://localhost:8080/strings?is_palindrome=true&min_length=5&word_count=1&contains_character=a"
```

## Examples

### Example 1: Create and Retrieve a Palindrome

```bash
# Create
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value":"racecar"}'

# Get specific string
curl http://localhost:8080/strings/racecar

# Get all palindromes
curl "http://localhost:8080/strings?is_palindrome=true"
```

### Example 2: Use Natural Language Filtering

```bash
# Find all single-word palindromic strings
curl "http://localhost:8080/strings/filter-by-natural-language?query=single%20word%20palindromes"

# Find strings longer than 10 characters
curl "http://localhost:8080/strings/filter-by-natural-language?query=strings%20longer%20than%2010%20characters"

# Find palindromes containing the letter 'a'
curl "http://localhost:8080/strings/filter-by-natural-language?query=palindromes%20containing%20the%20letter%20a"
```

### Example 3: Complex Filtering

```bash
# Find strings with 2-5 characters that contain the letter 'e'
curl "http://localhost:8080/strings?min_length=2&max_length=5&contains_character=e"

# Find all 3-word strings that are NOT palindromes
curl "http://localhost:8080/strings?word_count=3&is_palindrome=false"
```

### Example 4: Delete a String

```bash
curl -X DELETE http://localhost:8080/strings/hello%20world
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `modernc.org/sqlite` | Latest | SQLite database driver for Go |

All dependencies are specified in `go.mod`. Install them with:

```bash
go mod download
```

## Environment Variables

Currently, the application does not require any environment variables. Configuration is handled through defaults:

| Configuration | Default | Description |
|---------------|---------|-------------|
| Database Path | `./data/strings.db` | SQLite database file location |
| Server Port | `8080` | HTTP server port |

To modify these, edit the source code:
- Database path: `database/database.go` - `InitDB()` function
- Server port: `main.go` - `http.ListenAndServe()` call

## Testing

### Manual Testing with curl

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Create a String
```bash
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value":"test"}'
```

#### Get All Strings
```bash
curl http://localhost:8080/strings
```

#### Get Specific String
```bash
curl http://localhost:8080/strings/test
```

#### Filter Strings
```bash
curl "http://localhost:8080/strings?is_palindrome=true&min_length=5"
```

#### Natural Language Query
```bash
curl "http://localhost:8080/strings/filter-by-natural-language?query=palindromes"
```

#### Delete String
```bash
curl -X DELETE http://localhost:8080/strings/test
```

### Testing Tips

- Use URL encoding for spaces: `%20`
- Use URL encoding for special characters as needed
- Test with both valid and invalid inputs to verify error handling
- Verify database persistence by restarting the server


## Key Features 

### String Analysis
Each string is analyzed to compute:
- **length**: Total character count
- **is_palindrome**: Whether it reads the same forwards and backwards (case-insensitive)
- **unique_characters**: Count of distinct characters
- **word_count**: Number of whitespace-separated words
- **sha256_hash**: Unique identifier and SHA-256 hash of the string
- **character_frequency_map**: Occurrence count for each character

### Database
- SQLite database for persistent storage
- Automatic schema creation on first run
- ACID transactions for data integrity
- Located at `./data/strings.db`

### Filtering
- **Query Filters**: Combine multiple filters for complex queries
- **Natural Language**: Parse human-readable queries and convert to filters
- **Validation**: Prevent conflicting filters (e.g., min > max)
- **Error Handling**: Clear error messages for invalid parameters

## Troubleshooting

### Server won't start
- Check if port 8080 is already in use
- Ensure `./data` directory can be created
- Verify SQLite driver is installed: `go mod download`

### Database errors
- Check disk space
- Ensure write permissions in project directory
- Delete `./data/strings.db` and restart to reset database

### Filter returns no results
- Verify strings exist: `GET /strings`
- Check filter syntax
- Use `/health` endpoint to verify server is running

