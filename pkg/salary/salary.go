package salary

import (
	"encoding/json"
	"log"
	"net/http"
)

type NetSalary struct {
	GrossSalaryPerYear  int `json:"gross_salary_per_year"`
	GrossSalaryPerMonth int `json:"gross_salary_per_month"`

	NetSalaryPerYear  int `json:"net_salary_per_year"`
	NetSalaryPerMonth int `json:"net_salary_per_month"`

	IRPFApplied int `json:"irpf_applied"`
}

type Salary struct {
	YearlyGrossSalary     int    `json:"yearly_gross_salary"`
	PaymentsPerYear       int    `json:"payments_per_year"`
	Age                   string `json:"age"`
	ProfessionalCategory  string `json:"professional_category"`
	ContractType          string `json:"contract_type"`
	FamilySituation       string `json:"family_situation"`
	ChildrenExclusive     bool   `json:"children_exclusive"`
	ChildrenYoungerThan25 int    `json:"children_younger_than_25"`
	ChildrenYoungerThan3  int    `json:"children_younger_than_3"`
}

// CreateHandler will add a user into the database
func (s *Salary) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var netSalary NetSalary

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
	//w.Write(j)

	log.Println("*********************************")

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

	s.RestSS()
	log.Println("Discount SS:", s.RestSS())

	netSalary.GrossSalaryPerYear = grossSalary
	netSalary.GrossSalaryPerMonth = grossPerMonth
	netSalary.NetSalaryPerMonth = s.SplitSalaryByPayments() - s.RestIRPPerMonth()
	netSalary.NetSalaryPerYear = s.YearlyGrossSalary - s.RestIRPPerYear()
	netSalary.IRPFApplied = irpf

	log.Println("netSalary:", netSalary)
	log.Println("*********************************")

	n, err := json.Marshal(netSalary)
	if err != nil {
		log.Println("err", err)
	}

	w.Write(n)

}

func (s *Salary) RestSS() int {
	return s.SplitSalaryByPayments() / 100 * 6
}

func (s *Salary) RestIRPPerMonth() int {
	return s.SplitSalaryByPayments() / 100 * s.CalculateIRPF()
}

func (s *Salary) RestIRPPerYear() int {
	return s.YearlyGrossSalary / 100 * s.CalculateIRPF()
}

func (s *Salary) SplitSalaryByPayments() int {
	return s.YearlyGrossSalary / s.PaymentsPerYear
}

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
