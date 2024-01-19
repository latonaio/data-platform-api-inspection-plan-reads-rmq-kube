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
	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				header = c.Header(mtx, input, output, errs, log)
			}()
		case "SpecGeneral":
			func() {
				specGeneral = c.SpecGeneral(mtx, input, output, errs, log)
			}()
		case "SpecDetail":
			func() {
				specDetail = c.SpecDetail(mtx, input, output, errs, log)
		case "ComponentComposition":
			func() {
				componentComposition = c.ComponentComposition(mtx, input, output, errs, log)	
		case "Inspection":
			func() 
				inspection = c.Inspection(mtx, input, output, errs, log)
		case "Operation":
			func() {
				operation = c.Operation(mtx, input, output, errs, log)
		default:
		}
		if len(*errs) != 0 {
			break
		}
	}

	data := &dpfm_api_output_formatter.Message{
		Header: 				header,
		SpecGeneral:          	specGeneral,
		SpecDetail:           	specDetail,
		ComponentComposition:	componentComposition,
		Inspection:           	inspection,
		Operation:            	operation,
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

func (c *DPFMAPICaller) SpecGeneral(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.SpecGeneral {
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	groupBy := "\nGROUP BY InspectionPlan"

	rows, err := c.db.Query(
		`SELECT
    InspectionPlan
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_spec_general_data
		` + where + groupBy + `;`)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToSpecGeneral(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) SpecGeneral(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.SpecGeneral {
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	if input.Header.SpecGeneral.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND IsMarkedForDeletion = %v", where, *input.Header.SpecGeneral.IsMarkedForDeletion)
	}
	groupBy := "\nGROUP BY InspectionPlan, IsMarkedForDeletion"

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_spec_general_data
		` + where + groupBy + `;`)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToSpecGeneral(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) SpecDetail(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.SpecDetail {
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	if input.Header.SpecDetail.SpecType != nil {
		where = fmt.Sprintf("%s\nAND SpecType = \"%s\"", where, *input.Header.SpecDetail.SpecType)
	}
	if input.Header.SpecDetail.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND IsMarkedForDeletion = %v", where, *input.Header.SpecDetail.IsMarkedForDeletion)
	}
	groupBy := "\nGROUP BY InspectionPlan, SpecType, IsMarkedForDeletion"

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_spec_detail_data
		` + where + groupBy + `;`)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToSpecDetail(rows)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) ComponentComposition(
	mtx *sync.Mutex,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	errs *[]error,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.ComponentComposition {
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	if input.Header.ComponentComposition.ComponentCompositionType != nil {
		where = fmt.Sprintf("%s\nAND ComponentCompositionType = \"%s\"", where, *input.Header.ComponentComposition.ComponentCompositionType)
	}
	if input.Header.ComponentComposition.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND IsMarkedForDeletion = %v", where, *input.Header.ComponentComposition.IsMarkedForDeletion)
	}
	groupBy := "\nGROUP BY InspectionPlan, ComponentCompositionType, IsMarkedForDeletion"

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_component_composition_data
		` + where + groupBy + `;`)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToComponentComposition(rows)
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
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	if input.Header.Inspection.Inspection != nil {
		where = fmt.Sprintf("%s\nAND Inspection = %d", where, *input.Header.Inspection.Inspection)
	}
	if input.Header.Inspection.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND IsMarkedForDeletion = %v", where, *input.Header.Inspection.IsMarkedForDeletion)
	}
	groupBy := "\nGROUP BY InspectionPlan, Inspection, IsMarkedForDeletion"

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_inspection_data
		` + where + groupBy + `;`)
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
	where := "WHERE 1 = 1"

	if input.Header.InspectionPlan != nil {
		where = fmt.Sprintf("%s\nAND InspectionPlan = %d", where, *input.Header.InspectionPlan)
	}
	if input.Header.Operation.Operations != nil {
		where = fmt.Sprintf("%s\nAND Operations = %d", where, *input.Header.Operation.Operations)
	}
	if input.Header.Operation.OperationsItem != nil {
		where = fmt.Sprintf("%s\nAND OperationsItem = %d", where, *input.Header.Operation.OperationsItem)
	}
	if input.Header.Operation.OperationID != nil {
		where = fmt.Sprintf("%s\nAND OperationID = %d", where, *input.Header.Operation.OperationID)
	}
	if input.Header.Operation.IsMarkedForDeletion != nil {
		where = fmt.Sprintf("%s\nAND IsMarkedForDeletion = %v", where, *input.Header.Operation.IsMarkedForDeletion)
	}
	groupBy := "\nGROUP BY InspectionPlan, Operations, OperationsItem, OperationID, IsMarkedForDeletion"

	rows, err := c.db.Query(
		`SELECT *
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_inspection_plan_operation_data
		` + where + groupBy + `;`)
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
