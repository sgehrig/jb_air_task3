# Air Research Task 3

A CLI tool for analyzing survey data with flexible querying and reporting.

## Setup

1. **Clone the repository:**
   ```shell
   git clone <repo-url>
   cd air-research-task3
   ```

2. **Install Go (if not already installed):**
   - [Download Go](https://golang.org/dl/)

3. **Install dependencies:**
   - This project uses only the Go standard library.

## Running Tests

To run all tests:
```shell
go test ./survey
```

## Running the REPL Interface

Start the CLI REPL:
```shell
go run main.go
```

You will enter an interactive prompt where you can use the commands described below.

---

## CLI Commands

### `list`
List all available questions in the survey schema.

### `search <query>`
Search for questions or responses matching the given query string.

- `<query>`: A string to search for in question keys, text, or response values.

### `responses [<ResponseQuery>]`
Show survey responses. Optionally filter and limit the output using a ResponseQuery string.

- `<ResponseQuery>`: (optional) See "ResponseQuery String" below for syntax.

### `subset <ResponseQuery>`
Show a subset of responses as specified by the ResponseQuery.

- `<ResponseQuery>`: Required. See "ResponseQuery String" below for syntax.

### `analyze <question_key>`
Show the distribution of answers for a single or multi-choice question, including counts, percentages, and an ASCII bar graph.

- `<question_key>`: The key of the question to analyze. Must be a single or multi-choice question.

### `clear`
Clear the screen.

### `quit`
Exit the REPL.

---

## ResponseQuery String

The `ResponseQuery` string is used to filter and select specific keys and ranges of responses. It is used in the `responses` and `subset` commands.

**Syntax:**
- `keys:<key1>,<key2>,...;range:[<start>..<end>]`
- Sections can be separated by `;` or newlines.
- Both `keys` and `range` are optional for `responses`, but required for `subset`.

**Examples:**
- `keys:name,email;range:[first+1..last-2]`
- `keys: id \n range: [0..5]`
- `range:[first..last]`
- `keys:x,y,z`
- `keys:'foo,bar', "baz qux", plain, 'with ''quote'''`

**Range Endpoints:**
- `first`, `last` (optionally with +N or -N, e.g., `first+2`, `last-1`)
- Integer index (e.g., `0`, `5`)

**Quoted Keys:**
- Use single or double quotes for keys containing commas, spaces, or quotes.
- Escaped quotes: `''` for single, `\"` for double.

**Behavior:**
- If no `keys` are specified, all keys are included.
- If no `range` is specified, the full range (`first..last`) is used.

---

## Example Usage

### Show all responses:
```
responses
```

### Show only the first 3 responses with keys "name" and "email":
```
responses keys:name,email;range:[first..2]
```

### Show a subset of responses (IDs 2 to 4):
```
subset range:[2..4]
```

### Analyze answer distribution for a question:
```
analyze favorite_color
```

---

## Output Example for `analyze` Command

```
Distribution for [favorite_color] (SC):
  red                   10   25.0% |███████        |
  blue                  20   50.0% |████████████   |
  green                 10   25.0% |███████        |
  (n/a)                 0    0.0%  |               |
         Total: 40
```

---

## Notes

- Only single and multi-choice questions can be analyzed with the `analyze` command.
- The tool is designed for extensibility and easy integration with survey data in Go.
- All commands and their arguments are case-sensitive and must be entered as shown.
- The `ResponseQuery` parser supports robust quoting and range selection as tested in `survey/response_query_test.go`.
