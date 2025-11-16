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

### Object Construction

#### Basic field selection

```bash
# Select only specific fields
$ tq '.employees[0] | {name, role}' testdata/company.toon
name: Alice Smith
role: Engineer

# Select from multiple objects
$ tq '.employees[] | {name, salary}' testdata/company.toon
name: Alice Smith
salary: 95000
---
name: Bob Johnson
salary: 85000
---
name: Charlie Brown
salary: 110000
---
name: Diana Prince
salary: 98000
---
name: Eve Wilson
salary: 87000
```

#### Renaming and computing fields

```bash
# Rename fields
$ tq '.employees[0] | {employeeName: .name, position: .role}' testdata/company.toon
employeeName: Alice Smith
position: Engineer

# Add computed fields
$ tq '.employees[0] | {name, role, annualBonus: (.salary * 0.1)}' testdata/company.toon
name: Alice Smith
role: Engineer
annualBonus: 9500
```

#### Nested object construction

```bash
# Create nested structure
$ tq '.employees[0] | {profile: {name, role}, compensation: {salary, active}}' testdata/company.toon --json
{"profile":{"name":"Alice Smith","role":"Engineer"},"compensation":{"salary":95000,"active":true}}

# Mix nested and flat fields
$ tq '.employees[0] | {id, person: {name, role}, status: .active}' testdata/company.toon --json
{"id":1,"person":{"name":"Alice Smith","role":"Engineer"},"status":true}
```

### Array Construction

#### Building arrays from fields

```bash
# Create array from multiple fields
$ tq '.employees[0] | [.name, .role, .salary]' testdata/company.toon --json
["Alice Smith","Engineer",95000]

# Array of all names
$ tq '[.employees[].name]' testdata/company.toon --json
["Alice Smith","Bob Johnson","Charlie Brown","Diana Prince","Eve Wilson"]

# Array of salaries
$ tq '[.employees[].salary]' testdata/company.toon --json
[95000,85000,110000,98000,87000]
```

#### Conditional array construction

```bash
# Build array with filter
$ tq '[.employees[] | select(.salary > 90000) | .name]' testdata/company.toon --json
["Alice Smith","Charlie Brown","Diana Prince"]

# Complex filtered array
$ tq '[.employees[] | select(.active and .salary > 85000) | {name, salary}]' testdata/company.toon --json
[{"name":"Alice Smith","salary":95000},{"name":"Charlie Brown","salary":110000},{"name":"Eve Wilson","salary":87000}]
```

#### Using range()

```bash
# Generate range
$ tq '[range(5)]' --json <<< 'null'
[0,1,2,3,4]

# Range with start and end
$ tq '[range(3;7)]' --json <<< 'null'
[3,4,5,6]

# Use range for indexing
$ tq '.employees | [range(3)] | map(.employees[.] | .name)' testdata/company.toon --json
["Alice Smith","Bob Johnson","Charlie Brown"]
```

### Complex Transformation Patterns

#### Restructure data completely

```bash
# Transform employee data to summary format
$ tq '.employees | map({employee: .name, details: {position: .role, compensation: .salary, status: (if .active then "Active" else "Inactive" end)}})' testdata/company.toon --json
[{"employee":"Alice Smith","details":{"position":"Engineer","compensation":95000,"status":"Active"}},{"employee":"Bob Johnson","details":{"position":"Designer","compensation":85000,"status":"Active"}},{"employee":"Charlie Brown","details":{"position":"Manager","compensation":110000,"status":"Active"}},{"employee":"Diana Prince","details":{"position":"Engineer","compensation":98000,"status":"Inactive"}},{"employee":"Eve Wilson","details":{"position":"Designer","compensation":87000,"status":"Active"}}]
```

#### Create lookup tables

```bash
# Create role-based lookup
$ tq '.employees | group_by(.role) | map({role: .[0].role, members: [.[].name], avgSalary: (([.[].salary] | add) / length)})' testdata/company.toon --json
[{"role":"Designer","members":["Bob Johnson","Eve Wilson"],"avgSalary":86000},{"role":"Engineer","members":["Alice Smith","Diana Prince"],"avgSalary":96500},{"role":"Manager","members":["Charlie Brown"],"avgSalary":110000}]
```

#### Combine multiple transformations

```bash
# Filter, transform, and aggregate
$ tq '[.employees[] | select(.active)] | {activeCount: length, totalSalary: ([.[].salary] | add), avgSalary: (([.[].salary] | add) / length), employees: [.[].name]}' testdata/company.toon --json
{"activeCount":4,"totalSalary":377000,"avgSalary":94250,"employees":["Alice Smith","Bob Johnson","Charlie Brown","Eve Wilson"]}
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

#### String interpolation

Build strings with embedded expressions:

```bash
# Basic interpolation
$ tq '.employees[0] | "\(.name) works as a \(.role)"' testdata/company.toon
Alice Smith works as a Engineer

# Multiple fields
$ tq '.employees[0] | "\(.name) - \(.role) - $\(.salary)"' testdata/company.toon
Alice Smith - Engineer - $95000

# With computation
$ tq '.employees[0] | "Annual bonus: $\(.salary * 0.1)"' testdata/company.toon
Annual bonus: $9500

# Formatted output
$ tq '.employees[] | "\(.name) (\(.role)): \(if .active then "Active" else "Inactive" end)"' testdata/company.toon
Alice Smith (Engineer): Active
Bob Johnson (Designer): Active
Charlie Brown (Manager): Active
Diana Prince (Engineer): Inactive
Eve Wilson (Designer): Active
```

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

# Multiple checks
$ tq '.employees[] | select(.name | startswith("Alice")) | .role' testdata/company.toon
Engineer
```

#### Split and join

```bash
# Split email by @
$ tq '.email | split("@")' testdata/sample.toon --json
["john","example.com"]

# Split by space
$ tq '.name | split(" ")' testdata/sample.toon --json
["John","Doe"]

# Join array elements
$ tq '[.employees[].name] | join(", ")' testdata/company.toon
"Alice Smith, Bob Johnson, Charlie Brown, Diana Prince, Eve Wilson"

# Join with newline
$ tq '[.users[].name] | join("\n")' testdata/users.toon
"Alice\nBob\nCharlie"
```

#### Case conversion

```bash
# Convert to uppercase
$ tq '.name | ascii_upcase' testdata/sample.toon
JOHN DOE

$ tq '.employees[].name | ascii_upcase' testdata/company.toon
ALICE SMITH
---
BOB JOHNSON
---
CHARLIE BROWN
---
DIANA PRINCE
---
EVE WILSON

# Convert to lowercase
$ tq '.name | ascii_downcase' testdata/sample.toon
john doe

# Use in transformation
$ tq '.employees[] | {name, upperName: (.name | ascii_upcase)}' testdata/company.toon
name: Alice Smith
upperName: ALICE SMITH
---
name: Bob Johnson
upperName: BOB JOHNSON
---
name: Charlie Brown
upperName: CHARLIE BROWN
---
name: Diana Prince
upperName: DIANA PRINCE
---
name: Eve Wilson
upperName: EVE WILSON
```

#### Trimming strings

```bash
# Remove prefix
$ tq '.email | ltrimstr("john@")' testdata/sample.toon
example.com

# Remove suffix
$ tq '.email | rtrimstr("example.com")' testdata/sample.toon
john@

# Chain trimming
$ tq '.email | ltrimstr("john@") | rtrimstr(".com")' testdata/sample.toon
example
```

#### Regular expressions

```bash
# Test if pattern matches
$ tq '.email | test("@example.com$")' testdata/sample.toon
true

# Test with pattern
$ tq '.employees[].name | test("^A")' testdata/company.toon
true
---
false
---
false
---
false
---
false

# Filter using regex
$ tq '.employees[] | select(.name | test("^[AB]")) | .name' testdata/company.toon
Alice Smith
Bob Johnson

# Extract with match
$ tq '.email | match("(.+)@(.+)") | .captures[0].string' testdata/sample.toon
john
```

#### String replacement

```bash
# Replace first occurrence
$ tq '.employees[0].role | sub("Engineer"; "Senior Engineer")' testdata/company.toon
Senior Engineer

# Replace all occurrences
$ tq '.name | gsub(" "; "_")' testdata/sample.toon
John_Doe

# Remove characters
$ tq '.name | gsub("[aeiou]"; "")' testdata/sample.toon
Jhn D

# Format names
$ tq '.employees[] | {name, slug: (.name | ascii_downcase | gsub(" "; "-"))}' testdata/company.toon
name: Alice Smith
slug: alice-smith
---
name: Bob Johnson
slug: bob-johnson
---
name: Charlie Brown
slug: charlie-brown
---
name: Diana Prince
slug: diana-prince
---
name: Eve Wilson
slug: eve-wilson
```

#### Format functions

```bash
# URL encode
$ tq '.name | @uri' testdata/sample.toon
John%20Doe

# Base64 encode
$ tq '.name | @base64' testdata/sample.toon
Sm9obiBEb2U=

# JSON encode (escape for JSON)
$ tq '.name | @json' testdata/sample.toon
"John Doe"

# Use in templates
$ tq '.employees[0] | "https://example.com/profile/\(.name | @uri)"' testdata/company.toon
https://example.com/profile/Alice%20Smith
```

#### Combining string functions

```bash
# Complex transformation
$ tq '.employees[] | {name, email: ((.name | ascii_downcase | gsub(" "; ".")) + "@company.com")}' testdata/company.toon
name: Alice Smith
email: alice.smith@company.com
---
name: Bob Johnson
email: bob.johnson@company.com
---
name: Charlie Brown
email: charlie.brown@company.com
---
name: Diana Prince
email: diana.prince@company.com
---
name: Eve Wilson
email: eve.wilson@company.com

# Parse and format
$ tq '.email | split("@") | "\(.[0]) at \(.[1])"' testdata/sample.toon
john at example.com

# Generate slugs
$ tq '.employees[] | {name, slug: (.name | ascii_downcase | gsub("[^a-z0-9]+"; "-") | ltrimstr("-") | rtrimstr("-"))}' testdata/company.toon
name: Alice Smith
slug: alice-smith
---
name: Bob Johnson
slug: bob-johnson
---
name: Charlie Brown
slug: charlie-brown
---
name: Diana Prince
slug: diana-prince
---
name: Eve Wilson
slug: eve-wilson
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

## 7. Conditionals

`tq` supports full conditional logic with if-then-else expressions, including elif and nested conditions.

### Basic if-then-else

```bash
# Simple conditional
$ tq '.employees[0] | if .salary > 90000 then "high salary" else "normal salary" end' testdata/company.toon
high salary

# Numeric comparison
$ tq '.users[] | if .age >= 30 then .name else empty end' testdata/users.toon
Bob
---
Charlie

# String comparison
$ tq '.employees[] | if .role == "Manager" then .name else empty end' testdata/company.toon
Charlie Brown

# Boolean check
$ tq '.employees[] | if .active then "\(.name) is active" else "\(.name) is inactive" end' testdata/company.toon
Alice Smith is active
---
Bob Johnson is active
---
Charlie Brown is active
---
Diana Prince is inactive
---
Eve Wilson is active
```

### elif (else if) Chains

```bash
# Grade by score
$ tq '.employees[] | {name, level: (if .salary >= 100000 then "Senior" elif .salary >= 90000 then "Mid" else "Junior" end)}' testdata/company.toon
name: Alice Smith
level: Mid
---
name: Bob Johnson
level: Junior
---
name: Charlie Brown
level: Senior
---
name: Diana Prince
level: Mid
---
name: Eve Wilson
level: Junior

# Multiple conditions
$ tq '.employees[] | if .salary > 100000 then "A" elif .salary > 95000 then "B" elif .salary > 85000 then "C" else "D" end' testdata/company.toon
B
---
D
---
A
---
B
---
C
```

### Nested Conditionals

```bash
# Nested if-then-else
$ tq '.employees[] | if .active then (if .salary > 90000 then "active-high" else "active-normal" end) else "inactive" end' testdata/company.toon
active-high
---
active-normal
---
active-high
---
inactive
---
active-normal

# Complex nesting with object construction
$ tq '.employees[] | {name, category: (if .role == "Engineer" then (if .salary > 95000 then "Senior Engineer" else "Engineer" end) elif .role == "Manager" then "Management" else .role end)}' testdata/company.toon
name: Alice Smith
category: Engineer
---
name: Bob Johnson
category: Designer
---
name: Charlie Brown
category: Management
---
name: Diana Prince
category: Senior Engineer
---
name: Eve Wilson
category: Designer
```

### Conditionals with Logical Operators

```bash
# AND operator
$ tq '.employees[] | if .active and .salary > 90000 then .name else empty end' testdata/company.toon
Alice Smith
---
Charlie Brown

# OR operator
$ tq '.employees[] | if .role == "Engineer" or .role == "Manager" then {name, role} else empty end' testdata/company.toon
name: Alice Smith
role: Engineer
---
name: Charlie Brown
role: Manager
---
name: Diana Prince
role: Engineer

# Combined logic
$ tq '.employees[] | if (.role == "Engineer" and .salary > 90000) or .role == "Manager" then "\(.name) - Leadership Track" else "\(.name) - Individual Contributor" end' testdata/company.toon
Alice Smith - Leadership Track
---
Bob Johnson - Individual Contributor
---
Charlie Brown - Leadership Track
---
Diana Prince - Leadership Track
---
Eve Wilson - Individual Contributor
```

### Conditionals in Object Construction

```bash
# Add conditional fields
$ tq '.employees[] | {name, role, status: (if .active then "✓ Active" else "✗ Inactive" end)}' testdata/company.toon
name: Alice Smith
role: Engineer
status: ✓ Active
---
name: Bob Johnson
role: Designer
status: ✓ Active
---
name: Charlie Brown
role: Manager
status: ✓ Active
---
name: Diana Prince
role: Engineer
status: ✗ Inactive
---
name: Eve Wilson
role: Designer
status: ✓ Active

# Multiple conditional fields
$ tq '.employees[] | {name, tier: (if .salary >= 100000 then "A" else "B" end), employment: (if .active then "Current" else "Former" end)}' testdata/company.toon
name: Alice Smith
tier: B
employment: Current
---
name: Bob Johnson
tier: B
employment: Current
---
name: Charlie Brown
tier: A
employment: Current
---
name: Diana Prince
tier: B
employment: Former
---
name: Eve Wilson
tier: B
employment: Current
```

### Conditionals in String Interpolation

```bash
# Status messages
$ tq '.employees[] | "\(.name): \(if .active then "Currently employed as \(.role)" else "No longer with company" end)"' testdata/company.toon
Alice Smith: Currently employed as Engineer
---
Bob Johnson: Currently employed as Designer
---
Charlie Brown: Currently employed as Manager
---
Diana Prince: No longer with company
---
Eve Wilson: Currently employed as Designer

# Format with symbols
$ tq '.employees[] | "\(if .active then "✓" else "✗" end) \(.name) - \(.role)"' testdata/company.toon
✓ Alice Smith - Engineer
---
✓ Bob Johnson - Designer
---
✓ Charlie Brown - Manager
---
✗ Diana Prince - Engineer
---
✓ Eve Wilson - Designer

# Complex formatting
$ tq '.employees[] | "[\(if .salary >= 100000 then "HIGH" elif .salary >= 90000 then "MID" else "STD" end)] \(.name) - $\(.salary)"' testdata/company.toon
[MID] Alice Smith - $95000
---
[STD] Bob Johnson - $85000
---
[HIGH] Charlie Brown - $110000
---
[MID] Diana Prince - $98000
---
[STD] Eve Wilson - $87000
```

### Practical Examples

#### Classify and report

```bash
# Employee performance tiers
$ tq '.employees | map({name, performance: (if .salary > 95000 and .active then "Exceeds Expectations" elif .salary > 85000 and .active then "Meets Expectations" elif .active then "Developing" else "Inactive" end)})' testdata/company.toon
[5]{name,performance}:
  Alice Smith,Meets Expectations
  Bob Johnson,Developing
  Charlie Brown,Exceeds Expectations
  Diana Prince,Inactive
  Eve Wilson,Meets Expectations
```

#### Generate reports

```bash
# Summary with conditionals
$ tq '{totalEmployees: (.employees | length), active: ([.employees[] | select(.active)] | length), highEarners: ([.employees[] | select(.salary > 95000)] | length), status: (if ([.employees[] | select(.active)] | length) > 3 then "Fully Staffed" else "Hiring" end)}' testdata/company.toon
totalEmployees: 5
active: 4
highEarners: 3
status: Fully Staffed
```

#### Filter and transform

```bash
# Conditional transformation
$ tq '[.employees[] | if .active then {name, role, status: "active", adjustedSalary: (.salary * 1.1)} else {name, role, status: "inactive", adjustedSalary: .salary} end]' testdata/company.toon --json
[{"name":"Alice Smith","role":"Engineer","status":"active","adjustedSalary":104500},{"name":"Bob Johnson","role":"Designer","status":"active","adjustedSalary":93500},{"name":"Charlie Brown","role":"Manager","status":"active","adjustedSalary":121000},{"name":"Diana Prince","role":"Engineer","status":"inactive","adjustedSalary":98000},{"name":"Eve Wilson","role":"Designer","status":"active","adjustedSalary":95700}]
```

## 8. Error Handling

`tq` provides robust error handling with optional access, try-catch, and alternative operators.

### Optional Access Operator (?)

The `?` operator returns `null` instead of throwing an error when accessing missing fields or invalid indices.

```bash
# Missing field
$ tq '.employees[0].department?' testdata/company.toon
null

# Array out of bounds
$ tq '.employees[99]?' testdata/company.toon
null

# Nested optional access
$ tq '.employees[0] | .metadata?.created?' testdata/company.toon
null

# Use with data that exists
$ tq '.employees[0].name?' testdata/company.toon
Alice Smith
```

### try-catch Expressions

Catch errors and provide fallback values:

```bash
# Basic try-catch
$ tq 'try .employees[0].invalid catch "field not found"' testdata/company.toon
null

# Division by zero
$ tq 'try (.employees[0].salary / 0) catch "division error"' testdata/company.toon
division error

# Type conversion error
$ tq 'try (.name | tonumber) catch 0' testdata/sample.toon
0

# Try-catch with computation
$ tq '.employees[] | {name, bonus: (try (.salary * 0.1) catch 0)}' testdata/company.toon
name: Alice Smith
bonus: 9500
---
name: Bob Johnson
bonus: 8500
---
name: Charlie Brown
bonus: 11000
---
name: Diana Prince
bonus: 9800
---
name: Eve Wilson
bonus: 8700
```

### Alternative Operator (//) in Depth

Provide default values for null or false fields:

```bash
# Simple default
$ tq '.employees[] | {name, department: (.department // "Unassigned")}' testdata/company.toon
name: Alice Smith
department: Unassigned
---
name: Bob Johnson
department: Unassigned
---
name: Charlie Brown
department: Unassigned
---
name: Diana Prince
department: Unassigned
---
name: Eve Wilson
department: Unassigned

# Chain multiple alternatives
$ tq '.employees[0] | .manager // .supervisor // "No manager assigned"' testdata/company.toon
No manager assigned

# With computations
$ tq '.employees[] | {name, status: (if .active then "Active" else "Inactive" end), team: (.team // "General")}' testdata/company.toon
name: Alice Smith
status: Active
team: General
---
name: Bob Johnson
status: Active
team: General
---
name: Charlie Brown
status: Active
team: General
---
name: Diana Prince
status: Inactive
team: General
---
name: Eve Wilson
status: Active
team: General
```

### Combining Error Handling Techniques

```bash
# Optional access with alternative
$ tq '.employees[] | {name, email: (.email? // "no-email@company.com")}' testdata/company.toon
name: Alice Smith
email: no-email@company.com
---
name: Bob Johnson
email: no-email@company.com
---
name: Charlie Brown
email: no-email@company.com
---
name: Diana Prince
email: no-email@company.com
---
name: Eve Wilson
email: no-email@company.com

# Try-catch with alternative
$ tq '.employees[] | {name, level: ((try .level catch null) // "Standard")}' testdata/company.toon
name: Alice Smith
level: Standard
---
name: Bob Johnson
level: Standard
---
name: Charlie Brown
level: Standard
---
name: Diana Prince
level: Standard
---
name: Eve Wilson
level: Standard

# Nested optional access with alternatives
$ tq '.employees[] | {name, contact: (.contact?.email? // .contact?.phone? // "No contact info")}' testdata/company.toon
name: Alice Smith
contact: No contact info
---
name: Bob Johnson
contact: No contact info
---
name: Charlie Brown
contact: No contact info
---
name: Diana Prince
contact: No contact info
---
name: Eve Wilson
contact: No contact info
```

### Practical Error Handling Examples

#### Safe data extraction

```bash
# Extract fields safely with defaults
$ tq '.employees | map({name, role, salary: (.salary // 0), active: (.active // false)})' testdata/company.toon
[5]{name,role,salary,active}:
  Alice Smith,Engineer,95000,true
  Bob Johnson,Designer,85000,true
  Charlie Brown,Manager,110000,true
  Diana Prince,Engineer,98000,false
  Eve Wilson,Designer,87000,true
```

#### Graceful degradation

```bash
# Build robust queries that don't fail on missing data
$ tq '.employees[] | "\(.name) - \(.role // "Unknown Role") - \((.salary // 0) | if . > 0 then "Salary: $\(.)" else "Salary not disclosed" end)"' testdata/company.toon
Alice Smith - Engineer - Salary: $95000
---
Bob Johnson - Designer - Salary: $85000
---
Charlie Brown - Manager - Salary: $110000
---
Diana Prince - Engineer - Salary: $98000
---
Eve Wilson - Designer - Salary: $87000
```

#### Safe aggregations

```bash
# Calculate with error handling
$ tq '{total: ([.employees[]? | .salary? // 0] | add), count: ([.employees[]?] | length), average: (([.employees[]? | .salary? // 0] | add) / ([.employees[]?] | length))}' testdata/company.toon
total: 475000
count: 5
average: 95000
```

## 9. Variables

Variables allow you to bind values for reuse in complex queries using the `as $var` syntax.

### Basic Variable Binding

```bash
# Store a value and reuse it
$ tq '.age as $a | {name, age: $a, next_year: ($a + 1)}' testdata/sample.toon
name: John Doe
age: 30
next_year: 31

# Calculate percentage
$ tq '.employees[0] | .salary as $s | {name, salary: $s, monthly: ($s / 12)}' testdata/company.toon
name: Alice Smith
salary: 95000
monthly: 7916.666666666667
```

### Multiple Variables

```bash
# Use multiple variables in calculation
$ echo '{"price":100,"quantity":5,"discount":0.1}' | tq '.price as $p | .quantity as $q | .discount as $d | {subtotal: ($p * $q), discount_amount: ($p * $q * $d), total: ($p * $q * (1 - $d))}' --json
{
  "discount_amount": 50,
  "subtotal": 500,
  "total": 450
}

# Chained variable assignments
$ echo '{"salary":100000}' | tq '.salary as $s | ($s * 0.1) as $bonus | ($s * 0.05) as $tax | {salary: $s, bonus: $bonus, tax: $tax, net: ($s + $bonus - $tax)}' --json
{
  "bonus": 10000,
  "net": 105000,
  "salary": 100000,
  "tax": 5000
}
```

### Variables in Filtering

```bash
# Use variable in select
$ tq '.employees[] | .salary as $s | select($s > 90000) | {name, salary: $s}' testdata/company.toon
name: Alice Smith
salary: 95000
---
name: Charlie Brown
salary: 110000
---
name: Diana Prince
salary: 98000

# Filter with computed variable
$ tq '.users[] | .age as $a | select($a >= 25 and $a <= 30) | {name, age: $a}' testdata/users.toon
name: Alice
age: 25
---
name: Bob
age: 30
```

### Variables Across Nested Iterations

```bash
# Preserve parent context in nested iteration
$ echo '{"users":[{"name":"Alice","scores":[85,90,95]},{"name":"Bob","scores":[70,75,80]}]}' | tq '.users[] | .name as $username | .scores[] | {user: $username, score: .}' --json
{
  "score": 85,
  "user": "Alice"
}
{
  "score": 90,
  "user": "Alice"
}
{
  "score": 95,
  "user": "Alice"
}
{
  "score": 70,
  "user": "Bob"
}
{
  "score": 75,
  "user": "Bob"
}
{
  "score": 80,
  "user": "Bob"
}

# Complex nested example with categories
$ echo '[{"category":"Electronics","items":[{"name":"Laptop","price":1200},{"name":"Mouse","price":25}]},{"category":"Books","items":[{"name":"Novel","price":15},{"name":"Textbook","price":80}]}]' | tq '.[] | .category as $cat | .items[] | {category: $cat, item: .name, price: .price}' --json
{
  "category": "Electronics",
  "item": "Laptop",
  "price": 1200
}
{
  "category": "Electronics",
  "item": "Mouse",
  "price": 25
}
{
  "category": "Books",
  "item": "Novel",
  "price": 15
}
{
  "category": "Books",
  "item": "Textbook",
  "price": 80
}
```

### Variables with Aggregations

```bash
# Calculate percentages relative to total
$ echo '[5,10,15,20]' | tq '[.[]] | add as $total | . | map({value: ., percentage: ((. / $total) * 100)})' --json
[
  {
    "percentage": 10,
    "value": 5
  },
  {
    "percentage": 20,
    "value": 10
  },
  {
    "percentage": 30,
    "value": 15
  },
  {
    "percentage": 40,
    "value": 20
  }
]

# Normalize values
$ echo '[100,200,300]' | tq '[.[]] | max as $max | . | map(. / $max)' --json
[
  0.3333333333333333,
  0.6666666666666666,
  1
]

# Group statistics
$ tq '.employees | length as $count | [.[].salary] | add as $total | {employee_count: $count, total_salary: $total, average_salary: ($total / $count)}' testdata/company.toon
employee_count: 5
total_salary: 475000
average_salary: 95000
```

### Variables with Conditionals

```bash
# Use variable in conditional
$ tq '.employees[] | .salary as $s | {name, salary: $s, level: (if $s >= 100000 then "senior" elif $s >= 90000 then "mid" else "junior" end)}' testdata/company.toon
name: Alice Smith
salary: 95000
level: mid
---
name: Bob Johnson
salary: 85000
level: junior
---
name: Charlie Brown
salary: 110000
level: senior
---
name: Diana Prince
salary: 98000
level: mid
---
name: Eve Wilson
salary: 87000
level: junior

# Conditional calculation with variable
$ echo '{"price":100,"discount":0.2,"premium":true}' | tq '.price as $p | .discount as $d | .premium as $prem | if $prem then ($p * (1 - $d) * 0.9) else ($p * (1 - $d)) end' --json
72
```

### Practical Examples

#### Calculate employee bonuses based on performance tiers

```bash
$ tq '.employees[] | .salary as $s | .active as $act | select($act) | {name, salary: $s, bonus: (if $s >= 100000 then ($s * 0.15) elif $s >= 90000 then ($s * 0.12) else ($s * 0.10) end)}' testdata/company.toon
name: Alice Smith
salary: 95000
bonus: 11400
---
name: Bob Johnson
salary: 85000
bonus: 8500
---
name: Charlie Brown
salary: 110000
bonus: 16500
---
name: Eve Wilson
salary: 87000
bonus: 8700
```

#### Generate report with computed metrics

```bash
$ echo '{"revenue":1000000,"costs":750000,"employees":50}' | tq '.revenue as $r | .costs as $c | .employees as $e | {revenue: $r, costs: $c, profit: ($r - $c), margin: ((($r - $c) / $r) * 100), revenue_per_employee: ($r / $e), profit_per_employee: (($r - $c) / $e)}' --json
{
  "costs": 750000,
  "margin": 25,
  "profit": 250000,
  "profit_per_employee": 5000,
  "revenue": 1000000,
  "revenue_per_employee": 20000
}
```

#### Create index mapping

```bash
$ echo '["apple","banana","cherry"]' | tq '. as $items | $items | to_entries | map({(.value): .key}) | add' --json
{
  "apple": 0,
  "banana": 1,
  "cherry": 2
}
```

## 10. Recursive Descent

The recursive descent operator `..` allows you to search through all levels of nested structures.

### Basic Recursive Descent

```bash
# Get all values recursively
$ echo '{"a":1,"b":{"c":2,"d":{"e":3}}}' | tq '..' --json
{
  "a": 1,
  "b": {
    "c": 2,
    "d": {
      "e": 3
    }
  }
}
1
{
  "c": 2,
  "d": {
    "e": 3
  }
}
2
{
  "e": 3
}
3

# Find all numbers in nested structure
$ echo '{"a":1,"b":{"c":2,"d":"text","e":{"f":3}}}' | tq '.. | select(type == "number")' --json
1
2
3
```

### Find Specific Field Values

```bash
# Find all 'name' fields recursively
$ echo '{"name":"John","user":{"name":"Alice","profile":{"name":"Bob"}}}' | tq '.. | .name? // empty'
"John"
"Alice"
"Bob"

# Find specific field value anywhere in structure
$ tq '.. | select(.name? == "Alice Smith")' testdata/company.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
```

### Type-Based Searches

```bash
# Find all strings
$ echo '{"a":"hello","b":{"c":"world","d":123}}' | tq '.. | select(type == "string")' --json
"hello"
"world"

# Find all objects
$ echo '{"a":{"b":1},"c":{"d":2}}' | tq '.. | select(type == "object")' --json
{
  "a": {
    "b": 1
  },
  "c": {
    "d": 2
  }
}
{
  "b": 1
}
{
  "d": 2
}

# Find all arrays
$ echo '{"a":[1,2],"b":{"c":[3,4]}}' | tq '.. | select(type == "array")' --json
[
  1,
  2
]
[
  3,
  4
]
```

### Conditional Recursive Searches

```bash
# Find all numbers greater than a threshold
$ echo '{"a":5,"b":{"c":15,"d":{"e":25}},"f":8}' | tq '.. | select(type == "number" and . > 10)' --json
15
25

# Find strings starting with prefix
$ echo '{"a":"test1","b":{"c":"other","d":"test2","e":{"f":"test3"}}}' | tq '.. | select(type == "string" and startswith("test"))' --json
"test1"
"test2"
"test3"

# Find objects that have a specific field
$ echo '{"users":[{"id":1,"name":"Alice"},{"name":"Bob"},{"id":2,"name":"Charlie"}]}' | tq '.. | select(type == "object" and has("id"))' --json
{
  "id": 1,
  "name": "Alice"
}
{
  "id": 2,
  "name": "Charlie"
}
```

### Practical Examples

#### Find all email addresses in nested structure

```bash
$ echo '{"user":{"email":"user@example.com","profile":{"contact":{"email":"alt@example.com"}}},"admin":{"email":"admin@example.com"}}' | tq '.. | .email? // empty'
"user@example.com"
"alt@example.com"
"admin@example.com"
```

#### Extract all numeric values for analysis

```bash
$ tq '.. | select(type == "number")' testdata/company.toon | head -10
2020
---
1
---
95000
---
2
---
85000
---
```

#### Find all active users at any level

```bash
$ echo '{"dept1":{"users":[{"name":"A","active":true},{"name":"B","active":false}]},"dept2":{"teams":{"alpha":{"users":[{"name":"C","active":true}]}}}}' | tq '.. | select(.active? == true)' --json
{
  "active": true,
  "name": "A"
}
{
  "active": true,
  "name": "C"
}
```

#### Collect all unique field names

```bash
$ echo '{"a":1,"b":{"c":2,"d":{"e":3}}}' | tq '[.. | objects | keys] | flatten | unique' --json
[
  "a",
  "b",
  "c",
  "d",
  "e"
]
```

#### Find maximum value anywhere in structure

```bash
$ echo '{"metrics":{"cpu":{"max":95,"avg":70},"memory":{"max":88,"avg":65}}}' | tq '[.. | select(type == "number")] | max' --json
95
```

## 11. Advanced Examples

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
