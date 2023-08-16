package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type QueryParams struct {
	MediaOnly bool
	SinceID   int64
	MaxID     int64
	Limit     int64
}

func (h *handler) FetchPublicTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams, err := queryParamsParser(r, 40 /* defulat limit */, 80 /* max limit */)
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

func queryParamsParser(r *http.Request, dlftLimit int64, maxLimit int64) (*QueryParams, error) {
	booleanParamParserWithDefualt := func(paramKey string, dflt bool) (bool, error) {
		paramValue := r.FormValue(paramKey)
		if paramValue == "" {
			return dflt, nil
		}

		v, err := strconv.ParseBool(paramValue)
		if err != nil {
			return false, fmt.Errorf("invalid query parameter for %s: %s", paramKey, err)
		}

		return v, nil
	}

	integerParamParserWithDefualt := func(paramKey string, dflt int64) (int64, error) {
		paramValue := r.FormValue(paramKey)
		if paramValue == "" {
			return dflt, nil
		}

		v, err := strconv.ParseInt(paramValue, 10 /* base */, 64 /* bitSize */)
		if err != nil {
			return -1, fmt.Errorf("invalid query parameter for %s: %s", paramKey, err)
		}

		return v, nil
	}

	mediaOnly, err := booleanParamParserWithDefualt("media_only", false)
	if err != nil {
		return nil, err
	}

	limit, err := integerParamParserWithDefualt("limit", dlftLimit)
	if err != nil {
		return nil, err
	} else if !(0 <= limit && limit <= maxLimit) {
		return nil, fmt.Errorf("exceed the limitation of public time lines: [defulat limit] \"%d\", [max limit] \"%d\", [current limit]: \"%d\"", dlftLimit, maxLimit, limit)
	}

	sinceID, err := integerParamParserWithDefualt("since_id", 0)
	if err != nil {
		return nil, err
	}

	maxID, err := integerParamParserWithDefualt("max_id", sinceID+limit)
	if err != nil {
		return nil, err
	}

	return &QueryParams{MediaOnly: mediaOnly, SinceID: sinceID, MaxID: maxID, Limit: limit}, nil
}
