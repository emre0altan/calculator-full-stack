# Prompt Log

User prompts received after the request to start this log are recorded below. Related follow-up requests may be combined into one entry.

## 1. Go calculator REST API

> create a small REST api project in Go for calculator operations. api will include addition, subtraction, multiplication, division, exponentiation, square root and percentage. place the project under a new backend folder at root folder of the project

## 2. Organize the Go project

> place the types, functions and tests as in an ideal folder organization of a go project

## 3. Validate calculator edge cases

> validate inputs and handle these edge cases:
>
> 1) addition overflow for add endpoint,
> 2) multiplication overflow for multiply endpoint
> 3) zero division for divide endpoint

## 4. Separate calculator operations

> separate operations into their own functions in Calculate function

## 5. Division overflow edge case

> check edge case of dividing huge number with a very small floating point which would cause overflow

## 6. Add modulus endpoint

> add modulus operator endpoint as well with input validation, edge case handling and tests

## 7. Create a React calculator UI

> create a simple React typescript calculator UI that consumes the api, dont overengineer it, keep it simple. it should work both on large screens and small mobile screens, must be responsive. project must be placed in web folder at root project. before any coding, suggest the component and state structure that actually needed and continue after approval

## 8. Use one calculation endpoint

> reduce api surface to one /calculate endpoint instead of having one endpoint per operation and send operation type as a parameter

## 9. Run web and backend together in Docker

> create a dockerfile to run both web and backend projects together

## 10. Improve and test the React app

> show the selected operator between the number inputs in the react app, use % of for percentage, add meaningful web tests, and improve the web folder and code organization where needed. keep it simple

## 11. Validate backend inputs and update documentation

> add ajv-like schema validation for the backend inputs, then update the root, web and backend readme files with setup instructions, api examples, design decisions, how to run tests and how to create coverage reports
