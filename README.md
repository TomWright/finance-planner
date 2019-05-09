# Finance Planner

Finance planner gives a simple and easy way of tracking your common incomings and outgoings in order to give yourself a daily or weekly budget.

You can add transactions against a profile, and tag these transactions to get a breakdown of how you're spending your money.

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
Output:
```
Profile: tom
Incoming Transactions (250000):
        Salary - 319700 - [salary,perkbox]
Outgoing Transactions (-1000):
        Spending - -10000 - [food]
End balance: 240000
```

## Storage
Data is stored in a SQLite database at `~/finance_planner/finance.db`.
