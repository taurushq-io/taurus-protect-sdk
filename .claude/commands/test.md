# Tester Mode

You are in TESTER mode. Create comprehensive tests for the code.

## Your Task

Create tests covering:

1. **Happy path**: Expected inputs → expected outputs
2. **Edge cases**: Empty, null, boundary values, unicode
3. **Error conditions**: Invalid inputs, failures, exceptions

## Test Structure

Use **Arrange-Act-Assert**:
```
Arrange: Set up test data and preconditions
Act:     Execute the code under test  
Assert:  Verify the results
```

## Naming Convention

`test_[unit]_[scenario]_[expectedResult]`

Examples:
- `test_calculateTotal_withEmptyCart_returnsZero`
- `test_validateEmail_withInvalidFormat_throwsError`

## Edge Cases Checklist

**Strings**: empty, whitespace only, unicode, very long, special chars
**Numbers**: zero, negative, max/min, NaN, Infinity
**Collections**: empty, single element, many elements, duplicates, null elements
**General**: null/undefined, concurrent access, timeouts

## Output Format

```
## Test Plan
| Test Case | Type | Priority |
|-----------|------|----------|
| [description] | Unit/Integration | High/Medium/Low |

## Test Implementation
[Complete test code]

## Running Tests
[Command to run]

## Coverage Notes
- Covered: [what's tested]
- Not covered: [what's intentionally skipped]
```

## Code to Test

$ARGUMENTS
