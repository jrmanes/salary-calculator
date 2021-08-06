package salary

import (
	"encoding/json"
	"log"
	"net/http"
)

type NetSalary struct {
	// GrossSalaryPerYear salary gross per year
	GrossSalaryPerYear float64 `json:"gross_salary_per_year"`
	// GrossSalaryPerMonth salary gross per month
	GrossSalaryPerMonth float64 `json:"gross_salary_per_month"`

	// NetSalaryPerYear net salary, calculated by year
	NetSalaryPerYear float64 `json:"net_salary_per_year"`
	// NetSalaryPerMonth net salary, calculated by month
	NetSalaryPerMonth float64 `json:"net_salary_per_month"`

	// IRPFApplied percentage of IRPF to apply
	IRPFApplied float64 `json:"irpf_applied"`
	// IRPFRetentionApplied percentage of IRPF retention per year
	IRPFRetentionAppliedPerYear float64 `json:"irpf_retention_applied_per_year"`
	// IRPFRetentionApplied percentage of IRPF retention per year
	IRPFRetentionAppliedPerMonth float64 `json:"irpf_retention_applied_per_month"`

	// RetentionSS quantity in euros to discount depending on the contract type
	RetentionSS float64 `json:"retention_ss"`
}

type Salary struct {
	// YearlyGrossSalary total salary gross per year
	YearlyGrossSalary float64 `json:"yearly_gross_salary"`
	// PaymentsPerYear number of payments received per year
	PaymentsPerYear float64 `json:"payments_per_year"`
	// Age years of the user
	Age string `json:"age"`
	// ProfessionalCategory category in which the user is, can be found in the contract
	ProfessionalCategory string `json:"professional_category"`
	// ContractType Type of the contract, A=Indefinido, B=Temporal
	ContractType string `json:"contract_type"`
	// FamilySituation
	FamilySituation string `json:"family_situation"`
	// ChildrenExclusive is has children
	ChildrenExclusive bool `json:"children_exclusive"`
	// ChildrenYoungerThan25 number of children under 25 years
	ChildrenYoungerThan25 int `json:"children_younger_than_25"`
	// ChildrenYoungerThan3 number of children under 3 years
	ChildrenYoungerThan3 int `json:"children_younger_than_3"`
}

// CalculateSalaryHandler will add a user into the database
func (s *Salary) CalculateSalaryHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println("err", err)
	}

	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// calculate the net salary and get it into var n
	n := s.CalculateNetSalary()

	w.Write(n)
}

// CalculateNetSalary is where we calculate the salary net using other methods
func (s *Salary) CalculateNetSalary() []byte {
	log.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	grossSalary := s.YearlyGrossSalary
	log.Println("salary gross per year", grossSalary)
	irpf := s.CalculateIRPF()

	grossPerMonth := s.SplitSalaryByPayments()

	s.RestRangesOfIRPF()

	// create a new NetSalary struct
	n := NetSalary{
		GrossSalaryPerYear:           grossSalary,
		GrossSalaryPerMonth:          grossPerMonth,
		NetSalaryPerYear:             s.YearlyGrossSalary - s.RestIRPPerYear(),
		NetSalaryPerMonth:            s.SplitSalaryByPayments() - s.RestIRPPerMonth() - s.RestCotizationBase(),
		IRPFApplied:                  irpf,
		IRPFRetentionAppliedPerMonth: s.RestIRPPerMonth(),
		IRPFRetentionAppliedPerYear:  s.RestIRPPerYear(),
		RetentionSS:                  s.RestCotizationBase(),
	}

	log.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	log.Println("netSalary:", n)
	log.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	netS, err := json.Marshal(n)
	if err != nil {
		log.Println("err", err)
	}

	return netS
}

// RestCotizationBase rest the Seguridad Social percentage
// if the contract is "indefinido/fijo" it has 6.35, otherwise 6.35
func (s *Salary) RestCotizationBase() float64 {
	salaryToFloat := s.MonthlyGrossSalary()
	var retention float64

	switch base := s.ContractType; {
	case base == "A":
		// contract "fijo"
		retention = 6.35
	case base == "B":
		// contract "temporal"
		retention = 6.4
	default:
		retention = 6.4
		log.Println("Contract Type is not A or B... set default value")
	}

	// return the retention apply
	return salaryToFloat / 100.00 * retention
}

// RestIRPPerMonth rest the irpf percentage from the gross salary
func (s *Salary) RestIRPPerMonth() float64 {
	return s.SplitSalaryByPayments() / 100 * s.CalculateIRPF()
}

// RestIRPPerYear rest the irpf from the gross year salary
func (s *Salary) RestIRPPerYear() float64 {
	return s.YearlyGrossSalary / 100 * s.CalculateIRPF()
}

// SplitSalaryByPayments returns the salary per number of payments
func (s *Salary) SplitSalaryByPayments() float64 {
	return s.YearlyGrossSalary / s.PaymentsPerYear
}

// RestRangesOfIRPF rest the ranges of IRPF
func (s *Salary) RestRangesOfIRPF() float64 {
	salaryTotal := s.YearlyGrossSalary
	salary := s.ToFloat()

	var (
		toDiscount   float64
		currentRange float64
		tmpDisc      float64
	)

	log.Println("\n.................................")

	if salaryTotal < 12450.0 || salaryTotal > 12450.0 {
		log.Println("< 12450.0 || salaryTotal > 12450.0")
		//salary = salaryTotal - (12450.0 * 0.19)
		currentRange = 12450.0
		if s.ToFloat() < 12450.0 {
			currentRange = 12459.0
		}
		log.Println("tmpdiscrange", currentRange)
		currentRange = currentRange * 0.19
		log.Println("currentRange with discount", currentRange)
		log.Println("toDiscount before:", toDiscount)
		toDiscount += currentRange

		log.Println("toDiscountafter:", toDiscount)
		log.Println("toDiscount", toDiscount)
		log.Println(".................................")
	}
	if salaryTotal >= 12450.0 {
		log.Println("12450.0")
		//salary = salaryTotal - ((20200.0 - 12450.0) * 0.24)
		//toDiscount += (20200.0 - 12450.0) * 0.24
		currentRange = 20200.0 - 12450.0
		if s.ToFloat() < 20200.0 {
			currentRange = s.ToFloat() - 12450.0
		}
		log.Println("currentRange", currentRange)
		currentRange = currentRange * 0.24
		log.Println("currentRange with discount", currentRange)
		log.Println("toDiscount before:", toDiscount)
		toDiscount += currentRange

		log.Println("toDiscountafter:", toDiscount)
		log.Println("toDiscount:", toDiscount)
		log.Println(".................................")
	}
	if salaryTotal >= 20200.0 {
		log.Println("20200.0")
		//salary = salaryTotal - ((35200.0 - 20200.0) * 0.30)
		//toDiscount +=  (35200.0 - 20200.0) * 0.30
		currentRange = 35200.0 - 20200.0
		if s.ToFloat() < 35200.0 {
			currentRange = s.ToFloat() - 20200.0
		}
		log.Println("tmpdiscrange", currentRange)
		currentRange = currentRange * 0.30
		log.Println("currentRange with discount", currentRange)
		log.Println("toDiscount before:", toDiscount)
		toDiscount += currentRange

		log.Println("toDiscountafter:", toDiscount)
		log.Println("toDiscount:", toDiscount)
		log.Println(".................................")
	}
	if salaryTotal >= 35200.0 {
		log.Println("35200.0")
		//salary = salaryTotal - ((60000.0 - 35200.0) * 0.37)
		//toDiscount += (60000.0 - 35200.0) * 0.37
		currentRange = 60000.0 - 35200.0
		if s.ToFloat() < 60000.0 {
			currentRange = s.ToFloat() - 35200.0
		}
		log.Println("tmpdiscrange", currentRange)
		currentRange = currentRange * 0.37
		log.Println("currentRange with discount", currentRange)
		log.Println("toDiscount before:", toDiscount)
		toDiscount += currentRange

		log.Println("toDiscountafter:", toDiscount)
		log.Println("toDiscount:", toDiscount)
		log.Println(".................................")
	}
	if salaryTotal >= 60000.0 {
		log.Println("60000.0")
		salary = salaryTotal - ((300000.0 - 60000.0) * 0.45)
		toDiscount += float64(s.YearlyGrossSalary) - salary
		log.Println("tmpDisc before:", tmpDisc)
		tmpDisc += toDiscount
		log.Println("tmpDisc after:", tmpDisc)

		log.Println("tmpDisc:", tmpDisc)
		log.Println("toDiscount:", toDiscount)
		log.Println(".................................")
	}
	if salaryTotal > 300000.0 {
		log.Println("300000.0")
		salary = salaryTotal - (300000.0 * 0.47)
		toDiscount += float64(s.YearlyGrossSalary) - salary
		log.Println("tmpDisc before:", tmpDisc)
		tmpDisc += toDiscount
		log.Println("tmpDisc after:", tmpDisc)

		log.Println("tmpDisc:", tmpDisc)
		log.Println("toDiscount:", toDiscount)
		log.Println(".................................")
	}

	log.Println("salary total: ", float64(s.YearlyGrossSalary), " - ", toDiscount, " --> ",
		(salary - toDiscount), " per month", ((salary - toDiscount) / 12))
	log.Println(".................................")

	salary = s.ToFloat() - toDiscount

	log.Println("tmpDisc", tmpDisc)
	log.Println("toDiscount", toDiscount)
	log.Println("SALARY:", salary)
	log.Println("SALARY with discount", s.ToFloat()-toDiscount)
	log.Println("salary - discount per MONTH:", s.YearlyGrossSalary, "-", "toDiscount", toDiscount, "=", ((s.ToFloat() - toDiscount) / 12))
	log.Println("SALARY AFTER DISCOUNTS", salary)
	log.Println("##############################")

	return salary
}

// CalculateIRPF check the percentage of irpf need depending on your salary range
func (s *Salary) CalculateIRPF() float64 {

	switch gross := s.YearlyGrossSalary; {
	case gross < 12450:
		return 19
	case gross >= 12450 && gross < 20200:
		return 24
	case gross >= 20200 && gross <= 35200:
		return 30
	case gross >= 35200 && gross <= 60000:
		return 37
	case gross >= 60000 && gross <= 300000:
		return 45
	case gross > 300000:
		return 47
	default:
		return 19
	}

	return 19
}

// MonthlyGrossSalary returns the salary casted to float64
func (s *Salary) MonthlyGrossSalary() float64 {
	return s.YearlyGrossSalary / s.PaymentsPerYear
}

// ToFloat returns the salary casted to float64
func (s *Salary) ToFloat() float64 {
	return float64(s.YearlyGrossSalary)
}
