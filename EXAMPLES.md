# tq Examples

This document provides practical examples of using `tq` to query and transform TOON files.

## Sample Data

All examples use the sample files in the `testdata/` directory:
- `sample.toon` - Simple user object
- `users.toon` - Array of users
- `company.toon` - Company data with employees and departments
- `products.toon` - Product catalog

## 1. Element Extraction

Extract specific fields or array elements from TOON data.

### Extract a single field

```bash
$ tq '.name' testdata/sample.toon
John Doe
```

### Extract nested field

```bash
$ tq '.address.city' testdata/company.toon
San Francisco
```

### Extract array element by index

```bash
$ tq '.users[0]' testdata/users.toon
name: Alice
age: 25
email: alice@example.com
```

### Extract specific field from array element

```bash
$ tq '.employees[2].name' testdata/company.toon
Charlie Brown
```

### Extract all elements from an array

```bash
$ tq '.users[]' testdata/users.toon
name: Alice
age: 25
email: alice@example.com
---
name: Bob
age: 30
email: bob@example.com
---
name: Charlie
age: 35
email: charlie@example.com
```

### Extract specific field from all array elements

```bash
$ tq '.employees[].name' testdata/company.toon
Alice Smith
---
Bob Johnson
---
Charlie Brown
---
Diana Prince
---
Eve Wilson
```

## 2. Filtering

Filter data based on conditions using `select()`.

### Filter by numeric comparison

```bash
# Find employees with salary > 90000
$ tq '.employees[] | select(.salary > 90000)' testdata/company.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
---
id: 3
name: Charlie Brown
role: Manager
salary: 110000
active: true
---
id: 4
name: Diana Prince
role: Engineer
salary: 98000
active: false
```

### Filter by string equality

```bash
# Find employee named "Alice Smith"
$ tq '.employees[] | select(.name == "Alice Smith")' testdata/company.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
```

### Filter by boolean value

```bash
# Find active employees
$ tq '.employees[] | select(.active == true)' testdata/company.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
---
id: 2
name: Bob Johnson
role: Designer
salary: 85000
active: true
---
id: 3
name: Charlie Brown
role: Manager
salary: 110000
active: true
---
id: 5
name: Eve Wilson
role: Designer
salary: 87000
active: true
```

### Complex filter with multiple conditions

```bash
# Find active engineers
$ tq '.employees[] | select(.active == true and .role == "Engineer")' testdata/company.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
```

### Filter products in stock

```bash
$ tq '.products[] | select(.inStock == true)' testdata/products.toon
id: 101
name: Laptop
price: 1299.99
inStock: true
tags:
- electronics
- computers
---
id: 102
name: Mouse
price: 29.99
inStock: true
tags:
- electronics
- accessories
---
id: 104
name: Monitor
price: 399.99
inStock: true
tags:
- electronics
- displays
```

## 3. Output Formatting

Control how data is formatted in the output.

### Default output (TOON format)

```bash
$ tq '.' testdata/sample.toon
name: John Doe
age: 30
email: john@example.com
active: true
```

### Output as JSON

```bash
$ tq --json '.' testdata/sample.toon
{"name":"John Doe","age":30,"email":"john@example.com","active":true}
```

### Pretty-print JSON (using jq directly)

```bash
$ tq --json '.' testdata/users.toon | jq '.'
{
  "users": [
    {
      "name": "Alice",
      "age": 25,
      "email": "alice@example.com"
    },
    ...
  ]
}
```

## 4. Data Transformation

Transform data into new structures.

### Extract specific fields only

```bash
$ tq '.users | map({name, email})' testdata/users.toon
[3]{name,email}:
  Alice,alice@example.com
  Bob,bob@example.com
  Charlie,charlie@example.com
```

### Rename fields

```bash
$ tq '.users | map({fullName: .name, contact: .email})' testdata/users.toon
[3]{fullName,contact}:
  Alice,alice@example.com
  Bob,bob@example.com
  Charlie,charlie@example.com
```

### Calculate new fields

```bash
# Add 10% raise to all salaries
$ tq '.employees | map({name, oldSalary: .salary, newSalary: (.salary * 1.1)})' testdata/company.toon
[5]{name,oldSalary,newSalary}:
  Alice Smith,95000,104500
  Bob Johnson,85000,93500
  Charlie Brown,110000,121000
  Diana Prince,98000,107800
  Eve Wilson,87000,95700
```

### Combine multiple fields

```bash
$ tq '.employees | map({name, info: "\(.role) - $\(.salary)"})' testdata/company.toon --json
{"name":"Alice Smith","info":"Engineer - $95000"}
{"name":"Bob Johnson","info":"Designer - $85000"}
{"name":"Charlie Brown","info":"Manager - $110000"}
{"name":"Diana Prince","info":"Engineer - $98000"}
{"name":"Eve Wilson","info":"Designer - $87000"}
```

### Group and aggregate

```bash
# Count employees by role
$ tq '[.employees[] | .role] | group_by(.) | map({role: .[0], count: length})' testdata/company.toon --json
[{"role":"Designer","count":2},{"role":"Engineer","count":2},{"role":"Manager","count":1}]
```

## Advanced Examples

### Pipe multiple operations

```bash
# Get names of active engineers with salary > 90000
$ tq '.employees[] | select(.active == true and .role == "Engineer" and .salary > 90000) | .name' testdata/company.toon
Alice Smith
```

### Array slicing

```bash
# Get first 2 employees
$ tq '.employees[:2]' testdata/company.toon --json
[{"id":1,"name":"Alice Smith","role":"Engineer","salary":95000,"active":true},{"id":2,"name":"Bob Johnson","role":"Designer","salary":85000,"active":true}]
```

### Sorting

```bash
# Sort employees by salary (descending)
$ tq '.employees | sort_by(.salary) | reverse | map({name, salary})' testdata/company.toon
[5]{name,salary}:
  Charlie Brown,110000
  Diana Prince,98000
  Alice Smith,95000
  Eve Wilson,87000
  Bob Johnson,85000
```

### Keys and values

```bash
# Get all field names from first employee
$ tq '.employees[0] | keys' testdata/company.toon --json
["active","id","name","role","salary"]
```

### Check if field exists

```bash
# Check which objects have an 'active' field
$ tq '.employees[] | select(has("active"))' testdata/company.toon
# (returns all employees since they all have 'active' field)
```

## Tips

1. **Use `--json` for easier piping to other tools**: When you need to process the output with other JSON tools, use `--json` flag.

2. **Combine with standard Unix tools**:
   ```bash
   tq '.employees[].name' testdata/company.toon | wc -l  # Count employees
   ```

3. **Test filters incrementally**: Build complex queries step by step:
   ```bash
   tq '.employees[]' testdata/company.toon                    # First, see all employees
   tq '.employees[] | select(.active)' testdata/company.toon  # Then add filter
   tq '.employees[] | select(.active) | .name' testdata/company.toon  # Finally, extract name
   ```

4. **Use jq documentation**: Since `tq` uses `jq` under the hood, you can refer to [jq's manual](https://stedolan.github.io/jq/manual/) for advanced filtering and transformation syntax.

## See Also

- [jq Manual](https://stedolan.github.io/jq/manual/) - Complete jq documentation
- [TOON Format Specification](https://github.com/toon-format/spec) - TOON format details
- [tq GitHub Repository](https://github.com/RHEMS-japan/tq) - Source code and issues
