# Survey CLI Tool

A Go-based CLI tool for exploring and analyzing survey data from Excel files. Supports advanced querying, filtering, and interactive exploration via a REPL interface.

## Setup

1. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd <repo-folder>
   ```

2. **Install Go (if not already):**
   - [Download Go](https://golang.org/dl/) and follow the installation instructions for your OS.

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Build the project:**
   ```sh
   go build
   ```

## Running the CLI/REPL

Ensure your survey Excel file (e.g., `so_2024_raw.xlsx`) is in the project directory.

Start the REPL:
```sh
./air-research-task3-vscode
```

You will see a prompt:
```
Survey CLI. Type 'help', 'list', ...
(survey)>
```

## Running Tests

To run all tests:
```sh
go test ./...
```

## Available Commands

- **list**
  - Lists all questions in the survey.

- **search <string>**
  - Searches for questions containing the given string in their key, text, or options.

- **responses <ResponseQuery>**
  - Shows responses for the given query. See below for ResponseQuery syntax.

- **subset <ResponseQuery>**
  - Creates a subset of responses matching the query.

- **analyze <question-key>**
  - Shows the distribution of answers for a single- or multi-choice question, including counts, percentages, and an ASCII bar graph.

- **clear**
  - Clears the screen.

- **quit**
  - Exits the REPL.

## ResponseQuery Syntax

A `ResponseQuery` string allows you to select specific columns (keys) and/or a range of responses. Syntax:

```
keys:<key1>,<key2>,...;range:[<start>..<end>]
```

- **keys:** Comma-separated list of question keys to include. Keys can contain spaces and may be quoted with single or double quotes if they contain commas or special characters.
  - Examples:
    - `keys:id,name`
    - `keys:'Q1: Age',"Q2: Gender"`
    - `keys:foo bar, baz qux`
    - `keys:'foo,bar',plain`

- **range:** Selects a range of responses by index or symbolic endpoints.
  - Format: `range:[<start>..<end>]`
  - `<start>` and `<end>` can be:
    - `first` (first response)
    - `last` (last response)
    - An integer index (e.g., `0`, `5`)
    - With offsets: `first+2`, `last-3`
  - Examples:
    - `range:[first..last]` (all responses)
    - `range:[0..9]` (first 10 responses)
    - `range:[first+5..last-2]`

- **Combined:**
  - `keys:id,name;range:[0..9]`
  - `keys:'Q1: Age';range:[first+10..first+19]`

- **Omitting keys or range:**
  - If `keys` is omitted, all columns are included.
  - If `range` is omitted, all responses are included.

## Example Usage

- List all questions:
  ```
  list
  ```
- Search for questions containing "age":
  ```
  search age
  ```
- Show responses for the first 5 responses and only for 'id' and 'Q1: Age':
  ```
  responses keys:id,'Q1: Age';range:[0..4]
  ```
- Analyze the distribution of answers for question 'Q2: Gender':
  ```
  analyze Q2: Gender
  ```
- Create a subset for responses where 'Q2: Gender' is 'Female':
  ```
  subset keys:id,'Q2: Gender';range:[first..last]
  ```

## Subset Command Details

The `subset` (or `subsets`, `sub`) command allows you to filter responses by a specific question and option, and further refine the output using the ResponseQuery syntax.

### Usage

```
subset <question-key> <option> [<ResponseQuery>]
```

- `<question-key>`: The key of the question to filter on.
- `<option>`: The answer option to match (case-insensitive).
- `[<ResponseQuery>]` (optional): Further restricts which columns and rows are shown, using the same syntax as the `responses` command.

### Examples

- Show all responses where `Q2: Gender` is `Female`:
  ```
  subset 'Q2: Gender' Female
  ```
- Show only the `id` and `Q1: Age` columns for responses where `Q2: Gender` is `Male`, for the first 10 matches:
  ```
  subset 'Q2: Gender' Male keys:id,'Q1: Age';range:[0..9]
  ```
- Show all columns for responses where `Q3: Country` is `Germany`:
  ```
  subset 'Q3: Country' Germany keys:*
  ```

### Notes
- The subset command supports the full ResponseQuery syntax for advanced filtering and column selection.
- If you specify `keys:*` in the query, all columns will be shown.
- If you specify additional keys, they will be shown in addition to the filter question.
- The range selector in the query allows you to limit the number of matching responses displayed.

---

For further details, see the code or use the `help` command in the REPL.
