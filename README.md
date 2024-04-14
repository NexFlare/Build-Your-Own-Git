# Build Your Own Git

This project is a simplified version of Git, written in Go. It explores the internal mechanics of Git by implementing core functionalities from scratch. This mini-Git is designed for educational purposes to provide insight into how Git manages and stores data, along with the basic functionalities that make Git a powerful tool for version control.

## Features

- **git init**: Initialize a new Git repository.
- **cat-file**: Display the type, size, and content of object files.
- **hash-object**: Compute object ID and optionally creates a blob from a file.
- **ls-tree**: Used to inspect a tree object (currently supports --name-only flag)
- **write-tree**: Create a tree object from the staging area.
- **commit-tree**: Create a commit object from a tree object.

## Installation

```bash
# Clone the repository
git clone https://github.com/NexFlare/build-git-go

# Change directory
cd build-git-go

# Build the project (Go compiler required)
go build
```
