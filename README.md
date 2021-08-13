[![Go Report Card](https://goreportcard.com/badge/github.com/jrmanes/salary-calculator)](https://goreportcard.com/report/github.com/jrmanes/salary-calculator)

# Salary Calculator

Is an API which calculates the net salary in Spain.

## Post Request
Example:

```json
{
  "yearly_gross_salary": 50000,
  "payments_per_year": 12,
  "age": "28",
  "professional_category": "A",
  "contract_type": "A",
  "family_situation": "",
  "children_exclusive": false,
  "children_younger_than_25": 0,
  "children_younger_than_3": 0
}
```

POST request to: http://localhost:8080/api/v1/


The response will be something like:

```json
{
  "gross_salary_per_year": 50000,
  "gross_salary_per_month": 4166,
  "net_salary_per_year": 31500,
  "net_salary_per_month": 2649,
  "irpf_applied": 37,
  "irpf_retention_applied_per_year": 18500,
  "irpf_retention_applied_per_month": 1517,
  "retention_ss": 3200
}
```

----
Jose Ramón Mañes
----
