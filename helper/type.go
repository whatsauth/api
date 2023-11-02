package helper

type QRStatus struct {
	Status  bool   `json:"status,omitempty"`
	QRCode  string `json:"qrcode,omitempty"`
	Message string `json:"message,omitempty"`
}
