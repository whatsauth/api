package helper

type QRStatus struct {
	Status  bool   `json:"status"`
	QRCode  string `json:"qrcode"`
	Message string `json:"message"`
}
