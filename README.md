# Usage Billing System

This is a simple Go program that processes usage records for resources and generates billing statements based on the usage duration.

## Overview

The program defines two structs:

1. `UsageRecord`: Represents a record of usage of a resource, containing fields like the event type (create/destroy), resource ID, account ID, and timestamp.

2. `BillingStatement`: Represents a billing statement for a resource, containing fields like the resource ID, account ID, and usage duration in minutes.

The main function `processUsage` takes an array of `UsageRecord` structs as input and returns an array of `BillingStatement` structs. It does this by:

1. Storing the "create" events in a map, using a key composed of the account ID and resource ID.
2. For each "destroy" event, it looks up the corresponding "create" event in the map.
3. If a matching "create" event is found, it calculates the duration between the "create" and "destroy" timestamps.
4. It then generates a `BillingStatement` struct with the resource ID, account ID, and calculated duration.

The `main` function demonstrates the usage of `processUsage` by providing a sample set of `UsageRecord` structs and printing the resulting `BillingStatement` structs.

## Usage

1. Clone the repository or download the source code.
2. Navigate to the project directory.
3. Build the program using `go build`.
4. Run the executable.

The program will output the billing statements based on the sample usage records provided in the `main` function.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
