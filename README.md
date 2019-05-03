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
Transactions:
	Monthly salary - 2000 - [salary]
	Train ticket - -430 - [commute]
End balance: 2570
```

## Storage
The save files are current stored in plain text under `~/finance_planner`.

This may be changed in the future but it is fine for now.