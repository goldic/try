# try

A lightweight error-handling library for Go that simplifies error checking and handling of panics. Inspired by the `try-catch` pattern from other languages, this library allows you to streamline your code and reduce repetitive `if err != nil` blocks.

## Overview

Go’s traditional approach to error handling, while explicit and predictable, can lead to verbose and cluttered code. Repeatedly checking for errors across multiple function calls can make code harder to read and maintain. Moreover, handling panics—whether from your own code or third-party libraries—requires extra care, and often leads to even more boilerplate code.

The **`try`** library addresses these issues by providing concise, readable functions to handle common error patterns and panic recovery.

## Features

- **Simplifies error handling** by automatically panicking on errors and recovering them in a structured way.
- **Ensures clean code** without repetitive `if err != nil` checks.
- **Easy panic recovery** with built-in `Catch` to convert panics into returned errors.
- **Custom error handling** with `Handle` for advanced error processing.

## Installation

```bash
go get github.com/goldic/try
```

## Example Usage

Without **`try`**, error handling in Go looks like this:

```go
func LoadJSON(rawURL string) (map[string]any, error) {
    resp, err := http.Get(rawURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch data: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }
	
    var result map[string]any
    if err = json.Unmarshal(data, &result); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
    }
	
    return result, nil
}
```

With **`try`**, the code becomes much cleaner and easier to follow:

```go
func LoadJSON(rawURL string) (result map[string]any, err error) {
    defer try.Catch(&err)

    resp := try.Val(http.Get(rawURL))
    defer resp.Body.Close()

    try.Require(resp.StatusCode == http.StatusOK, "unexpected status code")

    data := try.Val(io.ReadAll(resp.Body)
    try.Check(json.Unmarshal(data, &result))
    return
}
```

## Functions

### `try.Val(value T, err error) T`

Handles function calls that return an error. If the function returns an error, `Val` will panic.

```go
data := try.Val(io.ReadAll(resp.Body))
```

- **When to use:** For calls where you want the error to be handled automatically by the library.

### `try.Val2(v1 T1, v2 T2, err error) (T1, T2)`

```go
buf := bufio.NewReader(os.Stdin)
line, isPrefix := try.Val2(buf.ReadLine())
```

### `try.Check(error)`

A simpler form of `Val`, `Check` takes only the error argument, and if the error is not `nil`, it panics.

```go
try.Check(json.Unmarshal(data, v))
```

- **When to use:** For quick, inline error handling when you don't need to capture the result, but just want to verify that an error didn’t occur.

### `try.Require(ok bool, err any)`

Ensures that a condition is met. If the condition is false, it will panic with the provided message.

```go
try.Require(resp.StatusCode == http.StatusOK, "unexpected status code")
```

- **When to use:** For conditions that must be true for the program to continue.

### `try.Catch(*error)`

Catch from any panic and converts it to an error. This function is typically used with `defer` to ensure the function returns an error instead of crashing.

```go
func Foo() (err error) {
    defer try.Catch(&err)
    ....
}
```

- **When to use:** In functions where you want to ensure panics are caught and returned as errors.

### `try.Handle(handler func(error))`

Custom panic handler that allows you to log or process the error in a specific way before recovering.

```go
defer try.Handle(func(err error) {
    log.Printf("An error occurred: %v", err)
})
```

- **When to use:** When you want to do custom logging or processing of errors when a panic occurs.

## Why Use `try`?

- **Cleaner code:** Focus on your core logic instead of writing repetitive error checks.
- **Improved readability:** Your code becomes more concise and easier to understand.
- **Better panic management:** Automatically handle panics and convert them to errors, ensuring your application stays stable.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

