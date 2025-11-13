# Ant Solver: A word unscramber

[![Go version](https://img.shields.io/badge/go-1.x-blue.svg)](#)
[![License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](#)
[![Status](https://img.shields.io/badge/status-live-success.svg)](https://ant-solver.onrender.com)

A lightweight **Golang** server that generates valid **subsequential anagrams** —
words that can be formed by taking letters from an input string **in order**, possibly skipping some.

**Live demo:** [https://ant-solver.onrender.com](https://ant-solver.onrender.com)

---

## Table of Contents

- [Overview](#overview)
- [Example](#example)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API](#api)
- [Design](#design)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

`ant-solver` computes all **valid dictionary words** that can be formed as subsequences of a given input word.
Unlike full anagrams (which must use all letters exactly once), **subsequential anagrams** maintain the **relative letter order** but may skip letters.

For example, from `"tamponh"` you can form `"phantom"` by choosing letters in order (t → a → m → p → o → n → h).

---

## Example

Try it live:

**Request:**
```

[https://ant-solver.onrender.com/api/unscrambe/?q=tamponh](https://ant-solver.onrender.com/api/unscrambe?q=tamponh)

````

**Response:**
```json
{
  "3": ["tam","pat","ton","tho","nth","mho","hop","poh","tao","ham","pom","pho","mot","oat","amp","nam","mat","hot","nap","pan","pah","hoa","han","mop","tom","ant","hat","pot","map","pam","moa","poa","top","not","hap","ohm","opt","tap","nat","tan","mon","hon","noh","apt","man","mna"],
  "4": ["mano","tamp","moat","math","atop","oath","phot","phon","atom","tanh","than","pont","thon","opah","path","phat","moan","mona","noma","pant","moth"],
  "5": ["toman","panto","month","manto"],
  "6": ["potman","tampon","topman"],
  "7": ["phantom"]
}
````

Each key corresponds to the **word length**, and values are the list of subsequential anagrams of that length.

---

## Features

* Finds **all** valid subsequential anagrams from a dictionary
* Fast and memory-efficient Go implementation
* Simple HTTP API
* Easily deployable (e.g. Render, Fly.io, Railway)
* No external dependencies beyond the Go standard library

---

## Installation

```bash
git clone https://github.com/thomasbui93/ant-solver.git
cd ant-solver
go build -o ant-solver .
```

---

## Usage

Run the server locally:

```bash
./ant-solver
```

By default, it starts on port `3000`.

Example request:

```bash
curl "http://localhost:3000/?q=example"
```

---

## API

### `GET /`

**Query parameters:**

| Name | Type   | Required | Description                                          |
| ---- | ------ | -------- | ---------------------------------------------------- |
| `q`  | string | ✅        | Input string to generate subsequential anagrams from |

**Example:**

```
/?q=planet
```

**Response:**

```json
{
  "3": ["ant", "tan", "pan"],
  "4": ["plan", "neat"],
  "5": ["plane", "plant"]
}
```

---

## Design

### Algorithm

1. Load a dictionary of valid English words.
2. For each dictionary word, check if it can be **subsequenced** from the input string:

   * e.g., `"phantom"` is valid for `"tamponh"` because each letter appears in order.
3. Group results by word length for readability.

### Project Structure

```
internal/          # Core logic (dictionary loading, matcher)
assets/            # Word lists / static data
```

---

## Contributing

Contributions are welcome!

* Open an issue for bugs or ideas
* Submit a pull request
* Improve performance or add new dictionary sets

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE).

---

Built with ❤️ in Go.

> “Find the hidden words your letters can tell.”
