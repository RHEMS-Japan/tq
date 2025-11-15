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

## 5. Operators

`tq` supports all jq operators for arithmetic, comparison, logical operations, and more.

### Arithmetic Operators

#### Addition, subtraction, multiplication, division

```bash
# Basic arithmetic
$ tq '.age + 5' testdata/sample.toon
35

$ tq '.age - 5' testdata/sample.toon
25

$ tq '.age * 2' testdata/sample.toon
60

$ tq '.age / 2' testdata/sample.toon
15

# Modulo (remainder)
$ tq '.age % 7' testdata/sample.toon
2
```

#### String concatenation

```bash
# Concatenate strings
$ tq '.name + " (" + .email + ")"' testdata/sample.toon
John Doe (john@example.com)

# Build formatted strings
$ tq '.employees[0] | .name + " - " + .role' testdata/company.toon
Alice Smith - Engineer
```

#### Array concatenation

```bash
# Combine arrays
$ echo '{"a":[1,2],"b":[3,4]}' | tq '.a + .b' --json
[1,2,3,4]

# Add element to array
$ echo '[1,2,3]' | tq '. + [4,5]' --json
[1,2,3,4,5]
```

#### Object merge

```bash
# Merge two objects
$ echo '{"a":1,"b":2}' | tq '. + {c:3}' --json
{"a":1,"b":2,"c":3}

# Override values
$ echo '{"a":1,"b":2}' | tq '. + {b:10}' --json
{"a":1,"b":10}
```

### Comparison Operators

#### Numeric comparison

```bash
# Greater than
$ tq '.employees[] | select(.salary > 90000) | .name' testdata/company.toon
Alice Smith
Charlie Brown
Diana Prince

# Less than or equal
$ tq '.employees[] | select(.salary <= 87000) | .name' testdata/company.toon
Bob Johnson
Eve Wilson

# Range check
$ tq '.users[] | select(.age >= 25 and .age <= 30)' testdata/users.toon
name: Alice
age: 25
email: alice@example.com
---
name: Bob
age: 30
email: bob@example.com
```

#### Equality comparison

```bash
# String equality
$ tq '.employees[] | select(.role == "Engineer") | .name' testdata/company.toon
Alice Smith
Diana Prince

# Boolean check
$ tq '.employees[] | select(.active == true) | .name' testdata/company.toon
Alice Smith
Bob Johnson
Charlie Brown
Eve Wilson

# Not equal
$ tq '.employees[] | select(.role != "Manager") | .name' testdata/company.toon
Alice Smith
Bob Johnson
Diana Prince
Eve Wilson
```

### Logical Operators

#### AND operator

```bash
# Multiple conditions
$ tq '.employees[] | select(.salary > 90000 and .active == true) | .name' testdata/company.toon
Alice Smith
Charlie Brown

# Complex filters
$ tq '.employees[] | select(.role == "Engineer" and .salary > 90000 and .active) | .name' testdata/company.toon
Alice Smith
```

#### OR operator

```bash
# Any of multiple conditions
$ tq '.employees[] | select(.role == "Manager" or .salary > 95000) | .name' testdata/company.toon
Charlie Brown
Diana Prince

# Multiple role check
$ tq '.employees[] | select(.role == "Engineer" or .role == "Designer") | .name' testdata/company.toon
Alice Smith
Bob Johnson
Diana Prince
Eve Wilson
```

#### NOT operator

```bash
# Negate boolean
$ tq '.employees[] | select(.active | not) | .name' testdata/company.toon
Diana Prince

# Combine with other conditions
$ tq '.employees[] | select(.role == "Engineer" and (.active | not)) | .name' testdata/company.toon
Diana Prince
```

### Alternative Operator (//)

The alternative operator `//` returns the right side if the left side is `null` or `false`.

```bash
# Provide default value
$ echo '{"a":null,"b":"value"}' | tq '.a // "default"' --json
"default"

$ echo '{"a":"value","b":"fallback"}' | tq '.a // "default"' --json
"value"

# Chain alternatives
$ echo '{"a":null,"b":null,"c":"final"}' | tq '.a // .b // .c // "none"' --json
"final"

# Use with missing fields
$ echo '{"name":"John"}' | tq '.age // 0' --json
0
```

### Combining Operators

Build complex queries by combining multiple operators:

```bash
# Calculate adjusted salaries
$ tq '.employees[] | {name, oldSalary: .salary, newSalary: (.salary * 1.1)}' testdata/company.toon
name: Alice Smith
oldSalary: 95000
newSalary: 104500
---
name: Bob Johnson
oldSalary: 85000
newSalary: 93500
---
name: Charlie Brown
oldSalary: 110000
newSalary: 121000
---
name: Diana Prince
oldSalary: 98000
newSalary: 107800
---
name: Eve Wilson
oldSalary: 87000
newSalary: 95700

# Filter and calculate
$ tq '.employees[] | select(.salary > 90000 and .active) | {name, bonus: (.salary * 0.1)}' testdata/company.toon
name: Alice Smith
bonus: 9500
---
name: Charlie Brown
bonus: 11000

# Average salary calculation
$ tq '([.employees[].salary] | add) / ([.employees[]] | length)' testdata/company.toon
95000
```

## 6. Built-in Functions

`tq` supports a wide range of jq built-in functions for working with arrays, objects, strings, and types.

### Array Functions

#### Get length of array, object, or string

```bash
# Array length
$ tq '.users | length' testdata/users.toon
3

# Object key count
$ tq '.employees[0] | length' testdata/company.toon
5

# String length
$ tq '.name | length' testdata/sample.toon
8
```

#### Reverse an array

```bash
$ tq '.users | reverse | map(.name)' testdata/users.toon --json
["Charlie","Bob","Alice"]
```

#### Sort arrays

```bash
# Simple sort (ascending)
$ tq '[5,2,8,1,9] | sort' --json <<< '[5,2,8,1,9]'
[1,2,5,8,9]

# Sort by field
$ tq '.employees | sort_by(.salary) | map({name, salary})' testdata/company.toon
[5]{name,salary}:
  Bob Johnson,85000
  Eve Wilson,87000
  Alice Smith,95000
  Diana Prince,98000
  Charlie Brown,110000
```

#### Remove duplicates

```bash
$ tq 'unique' --json <<< '[1,2,2,3,1,3]'
[1,2,3]
```

#### Sum, min, and max

```bash
# Sum all salaries
$ tq '[.employees[].salary] | add' testdata/company.toon
475000

# Find minimum salary
$ tq '[.employees[].salary] | min' testdata/company.toon
85000

# Find maximum salary
$ tq '[.employees[].salary] | max' testdata/company.toon
110000
```

#### Get first and last elements

```bash
# First employee
$ tq '.employees | first | .name' testdata/company.toon
Alice Smith

# Last employee
$ tq '.employees | last | .name' testdata/company.toon
Eve Wilson
```

#### Flatten nested arrays

```bash
$ tq 'flatten' --json <<< '[[1,2],[3,4],[5]]'
[1,2,3,4,5]

# Deep flatten
$ tq 'flatten' --json <<< '[[1,[2,3]],[[4]]]'
[1,2,3,4]
```

#### Group by expression

```bash
# Group employees by role
$ tq '.employees | group_by(.role) | map({role: .[0].role, employees: map(.name)})' testdata/company.toon --json
[{"role":"Designer","employees":["Bob Johnson","Eve Wilson"]},{"role":"Engineer","employees":["Alice Smith","Diana Prince"]},{"role":"Manager","employees":["Charlie Brown"]}]

# Count by role
$ tq '.employees | group_by(.role) | map({role: .[0].role, count: length})' testdata/company.toon
[3]{role,count}:
  Designer,2
  Engineer,2
  Manager,1
```

### Object Functions

#### Get object keys

```bash
# Get all field names
$ tq '.employees[0] | keys' testdata/company.toon --json
["active","id","name","role","salary"]

# Sort keys of object
$ tq 'keys' testdata/sample.toon --json
["active","age","email","name"]
```

#### Check if key exists

```bash
# Check if 'active' field exists
$ tq '.employees[0] | has("active")' testdata/company.toon
true

# Check if 'manager' field exists
$ tq '.employees[0] | has("manager")' testdata/company.toon
false
```

#### Convert to entries

```bash
# Convert object to key-value pairs
$ tq 'to_entries | .[0]' testdata/sample.toon
key: name
value: John Doe

# Transform using entries
$ tq '.employees[0] | to_entries | map(select(.key != "id"))' testdata/company.toon --json
[{"key":"name","value":"Alice Smith"},{"key":"role","value":"Engineer"},{"key":"salary","value":95000},{"key":"active","value":true}]
```

### Type Functions

#### Get type of value

```bash
# Check types
$ tq '.name | type' testdata/sample.toon
"string"

$ tq '.age | type' testdata/sample.toon
"number"

$ tq '.active | type' testdata/sample.toon
"boolean"

$ tq '.users | type' testdata/users.toon
"array"

$ tq 'type' testdata/company.toon
"object"
```

#### Type conversions

```bash
# Convert string to number
$ tq 'tonumber' --json <<< '"42"'
42

# Convert number to string
$ tq 'tostring' --json <<< '42'
"42"

# Use in transformations
$ tq '.employees | map({name, salary: (.salary | tostring)})' testdata/company.toon --json
[{"name":"Alice Smith","salary":"95000"},{"name":"Bob Johnson","salary":"85000"},{"name":"Charlie Brown","salary":"110000"},{"name":"Diana Prince","salary":"98000"},{"name":"Eve Wilson","salary":"87000"}]
```

### String Functions

#### String matching

```bash
# Check if string starts with prefix
$ tq '.name | startswith("John")' testdata/sample.toon
true

# Check if string ends with suffix
$ tq '.email | endswith("@example.com")' testdata/sample.toon
true

# Check if string contains substring
$ tq '.name | contains("Doe")' testdata/sample.toon
true
```

#### Split and join

```bash
# Split email by @
$ tq '.email | split("@")' testdata/sample.toon --json
["john","example.com"]

# Join array elements
$ tq '[.employees[].name] | join(", ")' testdata/company.toon
"Alice Smith, Bob Johnson, Charlie Brown, Diana Prince, Eve Wilson"
```

### Combining Functions

Build powerful transformations by combining multiple functions:

```bash
# Get unique roles, sorted
$ tq '[.employees[].role] | unique | sort' testdata/company.toon --json
["Designer","Engineer","Manager"]

# Calculate total salary by role
$ tq '.employees | group_by(.role) | map({role: .[0].role, total: ([.[].salary] | add)})' testdata/company.toon
[3]{role,total}:
  Designer,172000
  Engineer,193000
  Manager,110000

# Find highest paid employee
$ tq '.employees | sort_by(.salary) | reverse | first | {name, salary}' testdata/company.toon
name: Charlie Brown
salary: 110000

# Get average salary
$ tq '([.employees[].salary] | add) / (.employees | length)' testdata/company.toon
95000
```

## 7. Advanced Examples

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
