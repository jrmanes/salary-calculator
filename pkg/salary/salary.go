package salary

import (
	"encoding/json"
	"log"
	"net/http"
)

type NetSalary struct {
	// GrossSalaryPerYear salary gross per year
	GrossSalaryPerYear int `json:"gross_salary_per_year"`
	// GrossSalaryPerMonth salary gross per month
	GrossSalaryPerMonth int `json:"gross_salary_per_month"`

	// NetSalaryPerYear net salary, calculated by year
	NetSalaryPerYear int `json:"net_salary_per_year"`
	// NetSalaryPerMonth net salary, calculated by month
	NetSalaryPerMonth int `json:"net_salary_per_month"`

	// IRPFApplied percentage of IRPF to apply
	IRPFApplied int `json:"irpf_applied"`
	// IRPFRetentionApplied percentage of IRPF retention per year
	IRPFRetentionAppliedPerYear int `json:"irpf_retention_applied_per_year"`
	// IRPFRetentionApplied percentage of IRPF retention per year
	IRPFRetentionAppliedPerMonth int `json:"irpf_retention_applied_per_month"`

	// RetentionSS quantity in euros to discount depending on the contract type
	RetentionSS float64 `json:"retention_ss"`
}

type Salary struct {
	// YearlyGrossSalary total salary gross per year
	YearlyGrossSalary int `json:"yearly_gross_salary"`
	// PaymentsPerYear number of payments received per year
	PaymentsPerYear int `json:"payments_per_year"`
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

	n := s.CalculateNetSalary()

	w.Write(n)
}

// CalculateNetSalary is where we calculate the salary net using other methods
func (s *Salary) CalculateNetSalary() []byte {
	log.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	grossSalary := s.YearlyGrossSalary
	log.Println("salary gross per year", grossSalary)

	irpf := s.CalculateIRPF()
	log.Println("IRPF:", irpf)

	grossPerMonth := s.SplitSalaryByPayments()

	log.Println("salary gross per month:", grossPerMonth)
	log.Println("Discount after apply IRPF per YEAR:", s.RestIRPPerYear())
	log.Println("Discount after apply IRPF per MONTH:", s.RestIRPPerMonth())
	log.Println("Discount Base cotizacion SS:", s.RestCotizationBase())
	log.Println("salary after ranges discounts", s.RestRangesOfIRPF())

	// create a new NetSalary struct
	n := NetSalary{
		GrossSalaryPerYear:           grossSalary,
		GrossSalaryPerMonth:          grossPerMonth,
		NetSalaryPerYear:             s.YearlyGrossSalary - s.RestIRPPerYear(),
		NetSalaryPerMonth:            s.SplitSalaryByPayments() - s.RestIRPPerMonth(),
		IRPFApplied:                  irpf,
		IRPFRetentionAppliedPerMonth: s.RestIRPPerMonth(),
		IRPFRetentionAppliedPerYear:  s.RestIRPPerYear(),
		RetentionSS:                  s.RestCotizationBase(),
	}

	log.Println("netSalary:", n)
	log.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	netS, err := json.Marshal(n)
	if err != nil {
		log.Println("err", err)
	}

	return netS
}

// RestCotizationBase rest the Seguridad Social percentage
func (s *Salary) RestCotizationBase() float64 {
	salaryToFloat := float64(s.YearlyGrossSalary)
	var retention float64

	switch base := s.ContractType; {
	case base == "A":
		retention = 6.4
	case base == "B":
		retention = 6.35
	default:
		retention = 6.4
		log.Println("Conntract Type is not A or B... set default value")
	}

	return salaryToFloat / 100.00 * retention
}

// RestIRPPerMonth rest the irpf percentage from the gross salary
func (s *Salary) RestIRPPerMonth() int {
	return s.SplitSalaryByPayments() / 100 * s.CalculateIRPF()
}

// RestIRPPerYear rest the irpf from the gross year salary
func (s *Salary) RestIRPPerYear() int {
	return s.YearlyGrossSalary / 100 * s.CalculateIRPF()
}

// SplitSalaryByPayments returns the salary per number of payments
func (s *Salary) SplitSalaryByPayments() int {
	return s.YearlyGrossSalary / s.PaymentsPerYear
}

// RestRangesOfIRPF rest the ranges of IRPF
func (s *Salary) RestRangesOfIRPF() float64 {
	// TODO: calculate the discounts per each range and apply it to the gross salary

	salaryTotal := float64(s.YearlyGrossSalary)
	salary := s.ToFloat()
	log.Println("##############################")

	if salaryTotal < 12450.0 || salaryTotal > 12450.0 {
		salary = salaryTotal - (12450.0 * 0.19)
		log.Println("salary first range: ", salary)
	}
	if salaryTotal >= 12450.0 {
		salary = salaryTotal - ((20200.0 - 12450.0) * 0.24)
		log.Println("salary second range after discounts: ", salary)
	}
	if salaryTotal >= 20200.0 {
		salary = salaryTotal - ((35200.0 - 20200.0) * 0.30)
		log.Println("salary third after discounts: ", salary)
	}
	if salaryTotal >= 35200.0 {
		salary = salaryTotal - ((60000.0 - 35200.0) * 0.37)
		log.Println("salary fourth range: ", salary)
	}
	if salaryTotal >= 60000.0 {
		salary = salaryTotal - ((300000.0 - 60000.0) * 0.45)
		log.Println("salary fifth range: ", salary)
	}
	if salaryTotal > 300000.0 {
		salary = salaryTotal - (300000.0 * 0.47)
		log.Println("salary last range: ", salary)
	}

	log.Println("SALARY salaryTotal", salaryTotal)
	log.Println("SALARY AFTER DISCOUNTS", salary)
	log.Println("##############################")

	return salary
}

// CalculateIRPF check the percentage of irpf need depending on your salary range
func (s *Salary) CalculateIRPF() int {

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

// ToFloat returns the salary casted to float64
func (s *Salary) ToFloat() float64 {
	return float64(s.YearlyGrossSalary)
}
