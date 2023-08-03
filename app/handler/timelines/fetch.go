package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type QueryParams struct {
	SinceID int64
	MaxID   int64
	Limit   int64
}

func (h *handler) FetchPublicTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams, err := queryParamsParser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeline, err := h.sr.FindStatusesByRange(ctx, queryParams.SinceID, queryParams.MaxID, queryParams.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryParamsParser(r *http.Request) (*QueryParams, error) {
	integerParamParserWithDefualt := func(paramKey string, dflt int64) (int64, error) {
		paramValue := r.FormValue(paramKey)
		if paramValue == "" {
			return dflt, nil
		}

		v, err := strconv.ParseInt(paramValue, 10 /* base */, 64 /* bitSize */)
		if err != nil {
			return -1, fmt.Errorf("Invalid query parameter for %s: %s", paramKey, err)
		}

		return v, nil
	}

	limit, err := integerParamParserWithDefualt("limit", 40)
	if err != nil {
		return nil, err
	}

	sinceID, err := integerParamParserWithDefualt("since_id", 0)
	if err != nil {
		return nil, err
	}

	maxID, err := integerParamParserWithDefualt("max_id", sinceID+limit)
	if err != nil {
		return nil, err
	}

	return &QueryParams{SinceID: sinceID, MaxID: maxID, Limit: limit}, nil
}
