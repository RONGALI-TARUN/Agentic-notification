# Copilot Instructions

## Karpathy Guidelines

These guidelines bias toward caution over speed and are intended to reduce common LLM coding mistakes.

### 1. Think Before Coding
- State assumptions explicitly.
- Surface tradeoffs instead of silently choosing one path.
- Ask when the request is ambiguous.
- Prefer the simpler interpretation when multiple approaches are reasonable.

### 2. Simplicity First
- Implement the minimum necessary to solve the request.
- Avoid speculative features, abstractions, or configuration beyond the ask.
- Keep changes surgical and easy to review.

### 3. Surgical Changes
- Touch only the files and logic required by the request.
- Do not refactor unrelated code.
- Match the existing style and structure.
- Remove only the imports, variables, or functions introduced by your own changes when they become unused.

### 4. Goal-Driven Execution
- Define clear success criteria before implementing.
- Verify the result with tests or build commands when applicable.
- Prefer a short plan for multi-step work.

### Working style
- Avoid over-engineering.
- Leave extension points only when they are necessary for future growth and explicitly relevant to the request.
- Keep the implementation scoped to the current milestone unless the user explicitly asks otherwise.
