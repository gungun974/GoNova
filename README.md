# GoNova

A Go CLI to help create my highly opinionated Golang project templates.

## Features

- Default project file structure
- Pre-built module that you can modify to suit your needs
- Pretty logger
- Support for Postgres and SQLite with migrations
- Vite integration
- Templ and HTMX support
- Linting for Golang and TypeScript
- Custom HTTP response struct model and error handling
- Nix development environment and build

## Usage

### Init
```
gonova init [name] [flags]

Flags:
  -h, --help      help for init
      --nix       Init with nix module
      --no-git    Don't init the project with Git
      --postgre   Init with postgre module
      --sqlite    Init with sqlite module
```

### Dev

```
make dev
```

### Build

```
make build
```

## Todo

- [ ] Add Mailer Module
- [ ] Add Scheduler Module
- [ ] Add Ginko unit tests
- [ ] Add e2e tests starter
- [ ] Extract HTMX, Templ, and JS components into their own modules
