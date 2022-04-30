package certsurfer

import (
	"github.com/goccy/go-json"
	"time"
)

const (
	heartbeatType   = "heartbeat"
	certificateType = "certificate_update"
	domainType      = "dns_entries"
)

type MixedData struct {
	typ      string
	raw      []byte
	h        *HeartbeatData
	c        *CertificateData
	d        *DomainData
	pingStat PingStatus
}

func (d *MixedData) UnmarshalJSON(data []byte) error {
	var m struct {
		Type      string          `json:"message_type"`
		Data      json.RawMessage `json:"data"`
		Timestamp float64         `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	d.typ = m.Type
	switch m.Type {
	case heartbeatType:
		d.h = &HeartbeatData{}
		d.h.Timestamp = m.Timestamp
	case certificateType:
		d.c = &CertificateData{}
		return json.Unmarshal(m.Data, &d.c.Data)
	case domainType:
		d.d = &DomainData{}
		return json.Unmarshal(m.Data, &d.d.Data)
	default:
		d.raw = data
	}
	return nil
}

func (d MixedData) AsHeartbeat() (HeartbeatData, bool) {
	if d.h != nil {
		return *d.h, true
	}
	return HeartbeatData{}, false
}

func (d MixedData) AsCertificate() (CertificateData, bool) {
	if d.c != nil {
		return *d.c, true
	}
	return CertificateData{}, false
}

func (d MixedData) AsDomain() (DomainData, bool) {
	if d.d != nil {
		return *d.d, true
	}
	return DomainData{}, false
}

//AsUnknown cast mixed data into raw unknown data
//type in will always be provided
//raw is the raw data in bytes, it's only available for unknown data
//bool returns if the data is unknown
func (d MixedData) AsUnknown() (typ string, raw []byte, ok bool) {
	return d.typ, d.raw, len(d.raw) != 0
}

func (d MixedData) Ping() (PingStatus, bool) {
	return d.pingStat, !d.pingStat.IsZero()
}

type PingStatus struct {
	Latency time.Duration
	Time    time.Time
}

func (p PingStatus) IsZero() bool {
	return p.Latency == 0 && p.Time.IsZero()
}
