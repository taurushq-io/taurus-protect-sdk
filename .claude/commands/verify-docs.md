# Verify Documentation Against Code

Compare existing documentation with the actual codebase to find inconsistencies.

## Instructions

1. **Read all documentation** in `docs/` and `README.md`
2. **Compare against the specs (OpenAPI, proto)** to identify:
   - Missing documentation for public APIs
   - Outdated method signatures or parameters
   - Deprecated methods still documented as current
   - Code examples that no longer compile
   - Incorrect class/method names
   - Incorrect lifecycles (e.g., requests statuses not matching enums from OpenAPI)
   - Missing new features or classes
3**Compare against source code** to identify:
    - Missing documentation for public APIs
    - Outdated method signatures or parameters
    - Deprecated methods still documented as current
    - Code examples that no longer compile
    - Incorrect class/method names
    - Incorrect lifecycles (e.g., requests statuses not matching enums from OpenAPI)
    - Missing new features or classes

4**Generate a report** with:
    - ✅ Documentation that matches code
    - ⚠️ Documentation that needs updates
    - ❌ Missing documentation
    - 🗑️ Documentation for removed code

5**For each issue**, provide:
    - File location
    - Description of the problem
    - Suggested fix, propose technical solution if applicable

6**Fix all issues** to ensure documentation is up-to-date and accurate

## Output format:


$ARGUMENTS