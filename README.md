# RF Simulator

A Go-based simulator for understanding RF principles, TDMA frames, and modulation schemes.

## Prerequisites

- Go 1.22.4 or higher
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/yourusername/rf-simulator.git
cd rf-simulator
```

2. Build the project
```bash
make build
```

3. Run the simulator
```bash
make run
```

## Project Structure

- `cmd/`: Application entry points
- `internal/`: Private application code
  - `tdma/`: TDMA frame and burst implementations
  - `modulation/`: Modulation schemes and related code
  - `utils/`: Common utilities
- `pkg/`: Public packages that can be imported by other projects
- `Makefile`: Build and development commands

