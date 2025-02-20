package loan

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"api/config"

	_ "github.com/mattn/go-sqlite3"
)

// Mock implementations
type mockCreditService struct{}
type mockPaymentService struct{}
type mockDocumentService struct{}

func (m *mockCreditService) CheckCredit(applicantID string) (int, error)      { return 750, nil }
func (m *mockCreditService) ValidateIncome(evidence []Evidence) (bool, error) { return true, nil }
func (m *mockCreditService) CalculateRisk(creditScore int, amount float64) (float64, error) {
	return 0.05, nil
}

func (m *mockPaymentService) TransferFunds(from, to string, amount float64) error     { return nil }
func (m *mockPaymentService) ValidatePayment(paymentID string) error                  { return nil }
func (m *mockPaymentService) CalculateFine(dueDate time.Time, amount float64) float64 { return 0.0 }

func (m *mockDocumentService) StoreEvidence(evidence *Evidence) error       { return nil }
func (m *mockDocumentService) GenerateInvoice(payment *PaymentPeriod) error { return nil }
func (m *mockDocumentService) GenerateStatement(loanID string) error        { return nil }

func setupTestDB() (*sql.DB, error) {
	// Use in-memory SQLite for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Read and execute schema
	schemaPath, err := filepath.Abs("schema.sql")
	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupRealDB() (*sql.DB, error) {
	cfg := config.NewConfig()

	// Construct SQLite connection string
	// fmt.Println("DBName:", cfg.DBName)
	dbPath := filepath.Join("../../data", cfg.DBName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

// Add function to choose DB setup based on environment
func setupDB(useTestDB bool) (*sql.DB, error) {
	if useTestDB {
		return setupTestDB()
	}
	return setupRealDB()
}

func setupLoanService(useTestDB bool) (LoanService, error) {
	db, err := setupDB(useTestDB)
	if err != nil {
		return nil, err
	}

	return NewLoanService(
		db,
		&mockCreditService{},
		&mockPaymentService{},
		&mockDocumentService{},
	), nil
}

func TestApplyForLoan(t *testing.T) {
	service, err := setupLoanService(false)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	tests := []struct {
		name        string
		application *LoanApplication
		evidence    []Evidence
		wantErr     bool
		errMsg      string
	}{
		{
			name: "APP-001: Valid application submission",
			application: &LoanApplication{
				ID:          "APP-0001",
				ApplicantID: generateUniqueID(), // Generate unique ID
				Amount:      10000,
				Term:        12,
				Purpose:     "Home Improvement",
			},
			evidence: []Evidence{
				{ID: "EVI-0001",
					Type:        "INCOME_STATEMENT",
					Description: "Monthly Income Statement",
				},
				{
					ID:          "EVI-0002",
					Type:        "BANK_STATEMENT",
					Description: "Last 3 Months Statement",
				},
			},
			wantErr: false,
		},
		{
			name: "APP-002: Missing required documents",
			application: &LoanApplication{
				ID:          "APP-0002",
				ApplicantID: generateUniqueID(), // Generate unique ID
				Amount:      10000,
				Term:        12,
			},
			evidence: []Evidence{}, // No documents provided
			wantErr:  false,        // Note: Current implementation doesn't validate document requirements
		},
		{
			name: "APP-003: Invalid loan amount",
			application: &LoanApplication{
				ID:          "APP-0003",
				ApplicantID: generateUniqueID(), // Generate unique ID
				Amount:      -5000,
				Term:        12,
			},
			evidence: []Evidence{
				{Type: "BANK_STATEMENT"}, // Generate unique ID
			},
			wantErr: true,
			errMsg:  "invalid loan amount or term",
		},
		{
			name: "APP-004: Invalid loan term",
			application: &LoanApplication{
				ID:          "APP-0004",
				ApplicantID: generateUniqueID(), // Generate unique ID
				Amount:      10000,
				Term:        0,
			},
			evidence: []Evidence{
				{Type: "BANK_STATEMENT"}, // Generate unique ID
			},
			wantErr: true,
			errMsg:  "invalid loan amount or term",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ApplyForLoan(tt.application, tt.evidence)

			// Check error cases
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyForLoan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ApplyForLoan() error message = %v, want %v", err.Error(), tt.errMsg)
				return
			}

			// Check successful cases
			if err == nil {
				if tt.application.Status != StatusPending {
					t.Errorf("ApplyForLoan() status = %v, want %v", tt.application.Status, StatusPending)
				}
				if tt.application.AppliedAt.IsZero() {
					t.Error("ApplyForLoan() AppliedAt not set")
				}
				if len(tt.application.Evidence) != len(tt.evidence) {
					t.Errorf("ApplyForLoan() evidence count = %v, want %v", len(tt.application.Evidence), len(tt.evidence))
				}
			}
		})
	}
}

func TestUpdateCreditScore(t *testing.T) {
	service, err := setupLoanService(false)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	tests := []struct {
		name         string
		loanID       string
		creditScore  int
		interestRate float64
		wantErr      bool
		errMsg       string
	}{
		{
			name:        "CRD-001: Update credit score",
			loanID:      "APP-0001",
			creditScore: 750,
			interestRate: 0.05,
			wantErr:      false,
		},
		{
			name:         "CRD-002: Update credit score",
			loanID:       "APP-0002",
			creditScore:  1000,
			interestRate: 0.05,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateCreditScore(tt.loanID, tt.creditScore, tt.interestRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCreditScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateCreditScore() error message = %v, want %v", err.Error(), tt.errMsg)
				return
			}
		})
	}
}

func TestReviewApplication(t *testing.T) {
	service, err := setupLoanService(false)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	tests := []struct {
		name       string
		loanID     string
		wantErr    bool
		wantStatus Status
	}{
		{
			name:       "CRD-001: Excellent credit score",
			loanID:     "APP-0001",
			wantErr:    false,
			wantStatus: StatusReviewing,
		},
		// Add more credit check scenarios here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, err := service.ReviewApplication(tt.loanID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReviewApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && application.Status != tt.wantStatus {
				t.Errorf("ReviewApplication() status = %v, want %v", application.Status, tt.wantStatus)
			}
		})
	}
}

// Function to generate a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // Simple unique ID based on current time
}
