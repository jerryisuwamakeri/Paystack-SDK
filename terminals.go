package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// TerminalService groups the Paystack Terminal (POS) API operations.
// Access it via Client.Terminals.
type TerminalService struct {
	client *Client
}

// Terminal represents a registered Paystack Terminal device.
type Terminal struct {
	ID           int64  `json:"id"`
	SerialNumber string `json:"serial_number"`
	DeviceMake   string `json:"device_make,omitempty"`
	TerminalID   string `json:"terminal_id"`
	Status       string `json:"status,omitempty"`
}

// Fetch retrieves a terminal by its terminal ID.
func (s *TerminalService) Fetch(ctx context.Context, terminalID string, opts ...RequestOption) (*Terminal, *Response, error) {
	if terminalID == "" {
		return nil, nil, ErrMissingTerminalID
	}

	var out Terminal
	resp, err := s.client.do(ctx, "GET", "/terminal/"+url.PathEscape(terminalID), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SendTerminalEventRequest is the payload for TerminalService.SendEvent.
type SendTerminalEventRequest struct {
	// Type is the event type, for example "invoice" or "transaction".
	// Required.
	Type string `json:"type"`

	// Action is the action to take, for example "process", "view", or
	// "print". Required.
	Action string         `json:"action"`
	Data   map[string]any `json:"data,omitempty"`
}

// TerminalEvent identifies a queued terminal event.
type TerminalEvent struct {
	ID string `json:"id"`
}

// SendEvent pushes an event to a terminal for it to act on.
func (s *TerminalService) SendEvent(ctx context.Context, terminalID string, req SendTerminalEventRequest, opts ...RequestOption) (*TerminalEvent, *Response, error) {
	if terminalID == "" {
		return nil, nil, ErrMissingTerminalID
	}
	if req.Type == "" {
		return nil, nil, ErrMissingEventType
	}
	if req.Action == "" {
		return nil, nil, ErrMissingEventAction
	}

	var out TerminalEvent
	resp, err := s.client.do(ctx, "POST", "/terminal/"+url.PathEscape(terminalID)+"/event", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// TerminalEventStatus reports whether a queued terminal event has been
// delivered.
type TerminalEventStatus struct {
	Delivered bool `json:"delivered"`
}

// FetchEventStatus checks the delivery status of an event previously sent
// with SendEvent.
func (s *TerminalService) FetchEventStatus(ctx context.Context, terminalID, eventID string, opts ...RequestOption) (*TerminalEventStatus, *Response, error) {
	if terminalID == "" {
		return nil, nil, ErrMissingTerminalID
	}
	if eventID == "" {
		return nil, nil, ErrMissingEventID
	}

	var out TerminalEventStatus
	resp, err := s.client.do(ctx, "GET", "/terminal/"+url.PathEscape(terminalID)+"/event/"+url.PathEscape(eventID), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// TerminalStatus reports a terminal's connectivity.
type TerminalStatus struct {
	Online    bool `json:"online"`
	Available bool `json:"available"`
}

// FetchStatus checks whether a terminal is online and available.
func (s *TerminalService) FetchStatus(ctx context.Context, terminalID string, opts ...RequestOption) (*TerminalStatus, *Response, error) {
	if terminalID == "" {
		return nil, nil, ErrMissingTerminalID
	}

	var out TerminalStatus
	resp, err := s.client.do(ctx, "GET", "/terminal/"+url.PathEscape(terminalID)+"/presence", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateTerminalRequest is the payload for TerminalService.Update.
type UpdateTerminalRequest struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

// Update modifies a terminal's name or address.
func (s *TerminalService) Update(ctx context.Context, terminalID string, req UpdateTerminalRequest, opts ...RequestOption) (*Response, error) {
	if terminalID == "" {
		return nil, ErrMissingTerminalID
	}
	return s.client.do(ctx, "PUT", "/terminal/"+url.PathEscape(terminalID), req, nil, opts...)
}

// ListTerminalsParams filters TerminalService.List and
// TerminalService.ListAll.
type ListTerminalsParams struct {
	Page    int
	PerPage int
}

func (p ListTerminalsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of terminals.
func (s *TerminalService) List(ctx context.Context, params ListTerminalsParams, opts ...RequestOption) ([]Terminal, *Response, error) {
	path := "/terminal"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Terminal
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every terminal, fetching additional
// pages on demand as the caller advances it.
func (s *TerminalService) ListAll(ctx context.Context, params ListTerminalsParams, opts ...RequestOption) *Iterator[Terminal] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Terminal, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
