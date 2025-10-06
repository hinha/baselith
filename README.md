# baselith

Baselith is a cross-platform command-line database migration management tool built with Go and Cobra. It works on Linux, macOS, and Windows, and provides a unified approach to managing database schema changes across different database systems.

## Features

- Cross-platform support (Linux, macOS, Windows)
- Database migration management for PostgreSQL and MySQL
- Support for both transactional and non-transactional migrations
- XML-based migration configuration
- Support for both direct database connection parameters and YAML configuration files

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/hinha/baselith.git
cd baselith
```

2. Build the application:
```bash
cd cmd && go build -o baselith .
```

### From Release

Download the pre-built binary for your platform from the [Releases](https://github.com/hinha/baselith/releases) page.

## Usage

### Basic Commands

Run the application without arguments to see the welcome message:
```bash
./baselith
```

Check the version:
```bash
./baselith version
```

### Database Migration Commands

Execute the full migration to latest version:
```bash
./baselith --host=localhost --port=5432 --user=user --password=password --dbname=postgres --schema=public --driver=postgresql
```

Migration subcommands (default: "up"):
- `up` - Run all pending migrations
- `down` - Rollback the last migration or to a specified ID
- `to` - Migrate to a specific migration ID
- `redo` - Rollback and re-apply the latest migration

### Example Usage

Execute migrations with PostgreSQL:
```bash
./baselith --host=localhost --port=5432 --user=user --password=password --dbname=postgres --schema=public --driver=postgresql
```

Execute migrations with MySQL:
```bash
./baselith --host=localhost --port=3306 --user=user --password=password --dbname=mydb --driver=mysql
```

Use YAML configuration file:
```bash
./baselith --config=path/to/config.yaml
```

Specify migration subcommand:
```bash
./baselith --host=localhost --port=5432 --user=user --password=password --dbname=postgres --schema=public --driver=postgresql --s=down
```

Migrate to a specific version:
```bash
./baselith --host=localhost --port=5432 --user=user --password=password --dbname=postgres --schema=public --driver=postgresql --s=to --to=002
```

Get help:
```bash
./baselith --help
```

## Configuration

Baselith supports configuration via command-line flags or YAML file. Database connection parameters include:

- `--driver` - Database driver (postgres, mysql, etc.) [default: "postgres"]
- `--host` - Database host [default: "localhost"]
- `--port` - Database port [default: 5432]
- `--dbname` - Database name
- `--user` - Database user
- `--password` - Database password
- `--schema` - Database schema (for PostgreSQL) [default: "public"]
- `--s` - Subcommand to execute: up, down, to, redo [default: "up"]
- `--to` - Target migration ID for 'to' or 'down' subcommands
- `--config` - Path to configuration file
- `--yaml` - Output YAML configuration

## Migration Files

Baselith uses XML-based migration files where you can define:

- SQL migrations with up/down scripts
- Transactional and non-transactional migrations
- Metadata including author, labels, and migration types
- Sequential migration IDs for ordered execution


## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.