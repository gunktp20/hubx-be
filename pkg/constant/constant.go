package constant

import (
	"errors"
	"fmt"
)

type ErrorCode string

type ContractSummaryStrategy string

type ContractSummaryRemainType string

type ContractStrategyName string

type WorkOrderStatus string

type MatItemType string

type MatShipIcon string

type MatStockStatus string

const (
	ErrCodeInternalServer     ErrorCode = "50000"
	ErrCodeDatabaseFailure    ErrorCode = "50001"
	ErrCodeServiceUnavailable ErrorCode = "50002"

	ErrCodeUnauthorized ErrorCode = "40100"
	ErrCodeInvalidToken ErrorCode = "40101"
	ErrCodeTokenExpire  ErrorCode = "40102"

	ErrCodeNotFound         ErrorCode = "40400"
	ErrCodeResourceNotFound ErrorCode = "40401"

	ErrCodeBadRequest          ErrorCode = "40000"
	ErrCodeInvalidInputData    ErrorCode = "40001"
	ErrCodeMissingRequireField ErrorCode = "40002"

	ErrCodeForbidden    ErrorCode = "40300"
	ErrCodeAccessDenied ErrorCode = "40301"

	StrategyNegotiatedEnum ContractSummaryStrategy = "N"
	StrategyTenderEnum     ContractSummaryStrategy = "T"
	StrategySelectionEnum  ContractSummaryStrategy = "S"

	StrategyNameNegotiated ContractStrategyName = "Negotiate"
	StrategyNameTender     ContractStrategyName = "Tender"
	StrategyNameSelection  ContractStrategyName = "Selection"

	RemainTypeFirst   ContractSummaryRemainType = "timeLtSixMValueLqtTen"
	RemainTypeSecond  ContractSummaryRemainType = "timeLtSixMValueLqtThirty"
	RemainTypeThrid   ContractSummaryRemainType = "timeLqtYearValueLqtTen"
	RemainTypeFourth  ContractSummaryRemainType = "timeLqtYearValueLqtThirty"
	RemainTypeFifth   ContractSummaryRemainType = "timeLqtYearValueMtThirty"
	RemainTypeSixth   ContractSummaryRemainType = "timeLtSixMValueMtThirty"
	RemainTypeSeventh ContractSummaryRemainType = "timeMtYearValueLqtTen"
	RemainTypeEighth  ContractSummaryRemainType = "timeMtYearValueLqtThirty"
	RemainTypeNinth   ContractSummaryRemainType = "timeMtYearValueMtThirty"

	BaseRole = "HR"

	RoleModuleProcurementDB = "PcmDb"

	RolePagePcmDbContractSummary   = "ContractSummary"
	RolePagePcmDbContractSummaryFg = "ContractSummaryByFg"

	RolePermissionSuper  = "Super"
	RolePermissionMember = "Member"
	RolePermissionAdmin  = "Admin"

	InProgressOnTime    WorkOrderStatus = "INPROT"
	InProgressDelayed   WorkOrderStatus = "INPRDL"
	MaterialNotRequired WorkOrderStatus = "MATNTRQ"
	ReadyForExecution   WorkOrderStatus = "RDEXC"

	Material         MatItemType = "Material"
	Service          MatItemType = "Service"
	DirectPurchase   MatItemType = "DirectPurchase"
	StockReservation MatItemType = "StockReserve"

	Truck MatShipIcon = "truck"
	Plane MatShipIcon = "plane"
	Ship  MatShipIcon = "ship"
	None  MatShipIcon = ""

	WaitReserv  MatStockStatus = "W_RESERV"
	WaitOnHand  MatStockStatus = "W_ONHAND"
	INADEQUATE  MatStockStatus = "I"
	FullyIssued MatStockStatus = "FULL_ISSUE"
	PRApproved  MatStockStatus = "PR_APPRV"
	POApproved  MatStockStatus = "PO_APPRV"
	NullStatus  MatStockStatus = ""
)

func (e ErrorCode) String() string {
	return string(e)
}

type ErrorMessage string

const (
	ErrMessageSomethingWentWrong ErrorMessage = "something went wrong"
	ErrMessageAccessDenied       ErrorMessage = "access denied"
)

func (e ErrorMessage) String() string {
	return string(e)
}

func (r ContractSummaryRemainType) String() string {
	return string(r)
}

func (s ContractSummaryStrategy) String() string {
	return string(s)
}

func (s ContractStrategyName) String() string {
	return string(s)
}

var (
	API_PATH_POSTS         string = fmt.Sprintf("/%s", WORD_POSTS)
	API_PATH_HOME_CONTENTS string = fmt.Sprintf("/%s", WORD_HOME_CONTENTS)
	ErrorNotFoundDocument         = errors.New("not found document")
)

const (
	WORD_POSTS         string = "posts"
	WORD_HOME_CONTENTS string = "contents"
)
