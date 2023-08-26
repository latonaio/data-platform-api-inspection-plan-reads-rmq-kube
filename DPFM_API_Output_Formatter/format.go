package dpfm_api_output_formatter

import (
	"data-platform-api-inspection-plan-reads-rmq-kube/DPFM_API_Caller/requests"
	"database/sql"
	"fmt"
)

func ConvertToHeader(rows *sql.Rows) (*[]Header, error) {
	defer rows.Close()
	header := make([]Header, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.Header{}

		err := rows.Scan(
			&pm.InspectionPlan,
			&pm.InspectionPlantBusinessPartner,
			&pm.InspectionPlant,
			&pm.Product,
			&pm.ValidityStartDate,
			&pm.ValidityEndDate,
			&pm.ProductSpecification,
			&pm.InspectionSpecification,
			&pm.InspectionPlanHeaderText,
			&pm.CreationDate,
			&pm.LastChangeDate,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			return &header, err
		}

		data := pm
		header = append(header, Header{
			InspectionPlan:                 data.InspectionPlan,
			InspectionPlantBusinessPartner: data.InspectionPlantBusinessPartner,
			InspectionPlant:                data.InspectionPlant,
			Product:                        data.Product,
			ValidityStartDate:              data.ValidityStartDate,
			ValidityEndDate:                data.ValidityEndDate,
			ProductSpecification:           data.ProductSpecification,
			InspectionSpecification:        data.InspectionSpecification,
			InspectionPlanHeaderText:       data.InspectionPlanHeaderText,
			CreationDate:                   data.CreationDate,
			LastChangeDate:                 data.LastChangeDate,
			IsMarkedForDeletion:            data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		fmt.Printf("DBに対象のレコードが存在しません。")
		return &header, nil
	}

	return &header, nil
}
