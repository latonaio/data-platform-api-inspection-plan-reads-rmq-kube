package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-inspection-plan-reads-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-inspection-plan-reads-rmq-kube/DPFM_API_Output_Formatter"
	"fmt"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

func (c *DPFMAPICaller) readSqlProcess(
	ctx context.Context,
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	accepter []string,
	errs *[]error,
	log *logger.Logger,
) interface{} {
	var header *[]dpfm_api_output_formatter.Header
	var inspection *[]dpfm_api_output_formatter.Inspection
	var operation *[]dpfm_api_output_formatter.Operation

	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				header = c.Header(mtx, input, output, errs, log)
			}()
		case "Inspection":
			func() {
				inspection = c.Inspection(mtx, input, output, errs, log)
			}()
		case "Operation":
			func() {
				operation = c.Operation(mtx, input, output, errs, log)
			}()
		default:
		}
		if len(*errs) != 0 {
			break
		}
	}

	data := &dpfm_api_output_formatter.Message{
		Header:     header,
		Inspection: inspection,
		Operation:  operation,
	}

	return data
}

func (c *DPFMAPICaller) Header(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.Header {
	where := fmt.Sprintf("WHERE header.InspectionPlan = %d ", input.Header.InspectionPlan)

	if input.Header.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND header.IsMarkedForDeletion = %v", where, *input.Header.IsMarkedForDeletion)
	}

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_header_data AS header
		` + where + ` ORDER BY header.IsMarkedForDeletion ASC, header.IsCancelled ASC, header.InspectionPlan DESC;`,
	)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToHeader(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}
func (c *DPFMAPICaller) Inspection(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.Inspection {
	where := fmt.Sprintf("WHERE inspection.InspectionPlan = %d ", input.Inspection.InspectionPlan)

	if input.Inspection.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND inspection.IsMarkedForDeletion = %v", where, *input.Inspection.IsMarkedForDeletion)
	}

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_inspection_data AS header
		` + where + ` ORDER BY inspection.IsMarkedForDeletion ASC, inspection.IsCancelled ASC, inspection.InspectionPlan DESC;`,
	)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToInspection(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}
func (c *DPFMAPICaller) Operation(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.Operation {
	where := fmt.Sprintf("WHERE operation.InspectionPlan = %d ", input.Operation.InspectionPlan)

	if input.Operation.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND operation.IsMarkedForDeletion = %v", where, *input.Operation.IsMarkedForDeletion)
	}

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_operation_data AS operation
		` + where + ` ORDER BY operation.IsMarkedForDeletion ASC, operation.IsCancelled ASC, operation.InspectionPlan DESC;`,
	)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToOperation(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}
