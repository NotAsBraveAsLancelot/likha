
# Likha - The Tiny Test Data Generator

Likha is a high-performance, thread-safe Go CLI application designed to generate massive amounts of test data in various formats (CSV, JSON, XML, YAML). Built for simplicity and extreme scalability, Likha can generate billions of records with configurable fields and sophisticated value generation strategies.

## Features

- **High Performance**: Thread-safe and memory-efficient, capable of generating billions of records
- **Multiple Output Formats**: Support for CSV, JSON, XML, and YAML
- **Flexible Value Generation**: Six different generator types for maximum flexibility
- **Intuitive Scaling**: Use human-readable suffixes (10k, 10m, 10b) for record counts
- **Progress Tracking**: Real-time progress bar
- **Rich Configuration**: YAML-based configuration with extensive customization options

## Installation

```bash
go install likha@latest
```

Or build from source:

```bash
git clone likha.git
cd likha
go build -o likha
```

## Quick Start

1. Create a configuration file (`config.yaml`):

```yaml
fields:
  - name: "id"
    generator:
      type: "builtin"
      settings:
        function: "random_int"
        min: 1
        max: 1000
  - name: "name"
    generator:
      type: "list"
      settings:
        values: ["Alice", "Bob", "Charlie", "Diana"]

output:
  type: "json"
  file: "output.json"
  settings:
    pretty: true
```

2. Generate data:

```bash
likha --config config.yaml --count 1000
```

## Usage

```bash
likha [flags]

Flags:
  -c, --config string   Path to configuration file (default "config.yaml")
  -n, --count string    Number of records to generate (supports k, m, b suffixes)
  -h, --help           Help for likha
  -v, --version        Version information
```

### Count Format

Likha supports human-readable count formats:
- `100` - 100 records
- `10k` - 10,000 records
- `5m` - 5,000,000 records
- `2b` - 2,000,000,000 records

## Configuration

Likha uses YAML configuration files with two main sections: `fields` and `output`.

### Fields Configuration

Each field has a `name` and a `generator` with specific `type` and `settings`.

#### Generator Types

##### 1. Simple Value Generator
Uses a static value for all records:

```yaml
- name: "constant_field"
  generator:
    type: "simple"
    settings:
      value: "fixed_value"
```

##### 2. List Generator
Randomly selects from a predefined list:

```yaml
- name: "status"
  generator:
    type: "list"
    settings:
      values: ["active", "inactive", "pending"]
```

##### 3. Builtin Generator
Uses built-in functions for common data types:

```yaml
- name: "timestamp"
  generator:
    type: "builtin"
    settings:
      function: "random_isodate"
      start_date: "2022-01-01T00:00:00Z"
      end_date: "2023-12-31T23:59:59Z"
```

**Available builtin functions:**
- `random_epoch` - Unix timestamp
- `random_isodate` - ISO 8601 date string
- `random_string` - Random string with configurable length and character set
- `random_int` - Random integer within range
- `random_decimal` - Random decimal with configurable precision

##### 4. Expression Generator
Uses template expressions with field references and builtin functions:

```yaml
- name: "email"
  generator:
    type: "expression"
    settings:
      expression: "#name-$random_int(10,99)@example.com"
```

**Expression syntax:**
- `#field_name` - Reference to another field's value
- `$function_name(args)` - Call builtin function
- String concatenation with literal text

##### 5. Custom Generator
Executes external binary for each record:

```yaml
- name: "custom_data"
  generator:
    type: "custom"
    settings:
      binary_path: "/path/to/custom/generator"
      args: ["--format", "json"]
```

##### 6. ForeignKey Generator
Conditionally generates values based on another field:

```yaml
- name: "login_device"
  generator:
    type: "foreignkey"
    source_field: "status"
    map:
      "active":
        type: "list"
        settings:
          values: ["mobile", "desktop"]
      "inactive":
        type: "simple"
        settings:
          value: null
```

### Output Configuration

Configure output format and file settings:

```yaml
output:
  type: "json"  # csv, json, xml, yaml
  file: "output.json"
  settings:
    # JSON settings
    pretty: true

    # CSV settings
    # delimiter: ","
    # include_headers: true

    # XML settings
    # root_node: "records"

    # YAML settings
    # indent: 2
```

## Performance Considerations

- **Memory Usage**: Likha processes records in batches to maintain low memory footprint
- **Thread Safety**: All generators are thread-safe for concurrent execution
- **I/O Optimization**: Buffered writes minimize disk I/O overhead
- **Progress Tracking**: Non-blocking progress updates don't impact generation speed

## Examples

See the `examples/` directory for complete configuration examples including:
- Basic user data generation
- E-commerce product catalogs
- Log file simulation
- Complex relational data
- Custom format examples

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
