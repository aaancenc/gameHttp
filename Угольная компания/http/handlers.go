package http

import (
	"coal_company/company"
	"coal_company/company/equipment"
	"coal_company/company/miners"
	input_dto "coal_company/http/dto/input"
	output_dto "coal_company/http/dto/output"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPHandlers struct {
	company *company.Company

	closeServer func() error
}

func NewHTTPHandlers(copmany *company.Company) *HTTPHandlers {
	return &HTTPHandlers{
		company: copmany,
	}
}

func (handlers *HTTPHandlers) SetCloseServerFunc(f func() error) {
	handlers.closeServer = f
}

// Handler for hire new miner
func (h *HTTPHandlers) HandleCreateNewMiner(w http.ResponseWriter, r *http.Request) {
	var minerDTO input_dto.MinerDTO
	if err := json.NewDecoder(r.Body).Decode(&minerDTO); err != nil {
		http.Error(
			w,
			NewErrorDTO(err).ToString(),
			http.StatusBadRequest,
		)

		return
	}

	minerType := miners.MinerType(minerDTO.MinerType)

	miner, err := h.company.HireMiner(minerType)
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, company.ErrUnknokwnMinerType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, company.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(miner.Info(), "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// Handler for get info about miners (may consume query params)
func (h *HTTPHandlers) HandlerGetMiners(w http.ResponseWriter, r *http.Request) {
	minerType := r.URL.Query().Get("type")
	if minerType != "" {
		minersByType := h.company.GetMinersByType(miners.MinerType(minerType))
		outputDTO := output_dto.NewMinersByTypeDTO(minersByType)

		b, err := json.MarshalIndent(outputDTO, "", "    ")
		if err != nil {
			panic(err)
		}

		if _, err := w.Write(b); err != nil {
			http.Error(w, NewErrorDTO(err).ToString(), http.StatusInternalServerError)
		}

		return
	}

	allMiners := h.company.GetAllMiners()
	outputDTO := output_dto.NewAllMinersDTO(allMiners)

	b, err := json.MarshalIndent(outputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		http.Error(w, NewErrorDTO(err).ToString(), http.StatusInternalServerError)
	}
}

// Handler for buy new equipment
func (h *HTTPHandlers) HandleBuyEquipment(w http.ResponseWriter, r *http.Request) {
	var equipmentInputDTO input_dto.EquipmentDTO
	if err := json.NewDecoder(r.Body).Decode(&equipmentInputDTO); err != nil {
		http.Error(
			w,
			NewErrorDTO(err).ToString(),
			http.StatusBadRequest,
		)

		return
	}

	equipmentType := equipment.EquipmentType(equipmentInputDTO.EquipmentType)
	equipment, err := h.company.BuyEquipment(equipmentType)
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, company.ErrUnknokwnEquipmentType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, company.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	equipmentOutputDTO := output_dto.NewEquipmentDTO(equipment)
	b, err := json.MarshalIndent(equipmentOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
		return
	}

}

// Handler for check equipment current status
func (h *HTTPHandlers) HandleCheckEquipment(w http.ResponseWriter, r *http.Request) {
	equipment := h.company.GetEquipment()
	equipmentOutputDTO := output_dto.NewEquipmentDTO(equipment)

	b, err := json.MarshalIndent(equipmentOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// Handler for check company current statistics
func (h *HTTPHandlers) HandleGetCompanyStatistics(w http.ResponseWriter, r *http.Request) {
	statistics := h.company.GetStatistics()
	statisticsOutputDTO := output_dto.NewCompanyStatisticsDTO(statistics)

	b, err := json.MarshalIndent(statisticsOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// Handler for check miners salaries
func (h *HTTPHandlers) HandleGetMinersSalaries(w http.ResponseWriter, r *http.Request) {
	minersSalaresOutputDTO := output_dto.NewMinersSalariesDTO(
		miners.LittleMinerSalary,
		miners.BasicMinerSalary,
		miners.PowerfulMinerSalary,
	)

	b, err := json.MarshalIndent(minersSalaresOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// Handler for check equipment prices
func (h *HTTPHandlers) HandleGetEquipmentPrices(w http.ResponseWriter, r *http.Request) {
	equipmentPricesOutputDTO := output_dto.NewEquipmentPricesDTO(
		equipment.EquipmentPickaxeCost,
		equipment.EquipmentVentilationCost,
		equipment.EquipmentTrolleysCost,
	)

	b, err := json.MarshalIndent(equipmentPricesOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// Handler for request to game completion
func (h *HTTPHandlers) HandleCompleteGame(w http.ResponseWriter, r *http.Request) {
	statistics, err := h.company.Complete()
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, company.ErrNotAllEquipmentPurchased) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	statisticsOutputDTO := output_dto.NewCompanyStatisticsDTO(statistics)
	b, err := json.MarshalIndent(statisticsOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}

	go func() {
		if err := h.closeServer(); err != nil {
			fmt.Println("failed to close HTTP server:", err)
		}
	}()
}
