# Finance Planner

Finance planner gives a simple and easy way of tracking your common incomings and outgoings in order to give yourself a daily or weekly budget.

You can add transactions against a profile, and tag these transactions to get a breakdown of how you're spending your money.

## Installation

Download a binary from the releases page and save it as `/usr/local/bin/finance`.

## Usage

### Add a transaction
Transactions have a few properties:
- Label: What is it for?
- Amount: Positive or negative number that will impact your ending balance.
- Tags: Comma separate list of tags, used to categorise the transactions.
```
finance add-transaction --profile=tom --label="Train ticket" --amount=-43500 --tags=commute,travel
```

### List your transactions
```
finance list-transactions --profile=tom
```

- Append `--in` to only show incoming transactions
- Append `--out` to only show outgoing transactions

### Update a transaction
`update-transaction` looks a lot like `add-transaction`, but with an added `id` argument.

Any given values will overwrite existing values on the transaction.

```
finance update-transaction --id="tra:11111111-1111-1111-1111-111111111111" --profile=tom --label="Train ticket" --amount=-43500 --tags=commute,travel
```

## Storage
Data is stored in a SQLite database at `~/finance_planner/finance.db`.
