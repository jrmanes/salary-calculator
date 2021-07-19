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

	//j, err := json.Marshal(s)
	//if err != nil {
	//	log.Println("err", err)
	//}
	//
	//log.Println("j", j)

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

	s.RestIRPPerYear()
	log.Println("Discount after apply IRPF per YEAR:", s.RestIRPPerYear())

	s.RestIRPPerMonth()
	log.Println("Discount after apply IRPF per MONTH:", s.RestIRPPerMonth())

	s.RestCotizationBase()
	log.Println("Discount Base cotizacion SS:", s.RestCotizationBase())

	// create a new NetSalary struct
	n := NetSalary{
		GrossSalaryPerYear:  grossSalary,
		GrossSalaryPerMonth: grossPerMonth,
		NetSalaryPerYear:    s.YearlyGrossSalary - s.RestIRPPerYear(),
		NetSalaryPerMonth:   s.SplitSalaryByPayments() - s.RestIRPPerMonth(),
		IRPFApplied:         irpf,
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
	toFloat := float64(s.SplitSalaryByPayments())
	var retention float64

	switch base := s.ContractType; {
	case base == "A":
		retention = 6.4
	case base == "B":
		retention = 6.35
	default:
		retention = 6.4
		log.Println("ContrackType is not A or B... set default value")
	}

	return toFloat / 100.00 * retention
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

// CalculateIRPF check the percentage of irpf need depending on your salary range
func (s *Salary) CalculateIRPF() int {

	switch gros := s.YearlyGrossSalary; {
	case gros < 12450:
		return 19
	case gros >= 12450 && gros < 22000:
		return 24
	case gros >= 22000 && gros <= 35200:
		return 30
	case gros >= 35200 && gros <= 60000:
		return 37
	case gros >= 60000 && gros <= 300000:
		return 45
	case gros > 300000:
		return 47
	default:
		return 19
	}

	return 19
}
