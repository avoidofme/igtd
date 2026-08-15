
# igtd

A high-performance Command Line Interface (CLI) utility written in Go that applies an automated Dracula theme color mapping profile over raw graphic assets. It isolates primary color channels (Red, Green, Blue dominance) using density masks and smooth blending calculations to deliver standard Dracula aesthetic profiles.

## Features

- **Accurate Profile Mapping:** Replicates Python PIL-based image filters in native Go configurations.
- **Smart Arguments:** Powered by Cobra for robust input verification and dynamic error handling.
- **Auto Directory Resolution:** Automatically creates missing target folders and maps input names safely.

## Prerequisites

Make sure you have [Go](https://go.dev) installed on your system.

## Installation

1. Clone or navigate to your local project directory.
2. Install the tracking dependencies:

   ```bash
   go mod download
   ```

3. Compile the application binary layout target:

   ```bash
   go build -o igdt main.go
   ```

## Usage

Run the compiled executable file directly from your terminal.

### Basic Syntax

```bash
./igtd <path-to-image> [flags]
```

### Examples

**Process an image into the current directory:**
*This creates an automatically prefixed file named `igdt_image.jpg` in your current folder.*

```bash
./igtd image.jpg
```

**Save the processed output directly into a specific file route:**

```bash
./igtd image.jpg -o output.jpg
```

**Output directly into a distinct folder path:**
*If the directory does not exist, the tool automatically builds it for you.*

```bash
./igtd image.jpg ./output_folder
```
