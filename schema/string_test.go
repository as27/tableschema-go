package schema

import "testing"

// To be in par with the python library.
func TestDecodeString_URIMustRequireScheme(t *testing.T) {
	if _, err := decodeString(stringURI, "google.com", Constraints{}); err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestDecodeString_InvalidUUIDVersion(t *testing.T) {
	// This is a uuid3: namespace DNS and python.org.
	if _, err := decodeString(stringUUID, "6fa459ea-ee8a-3ca4-894e-db77e160355e", Constraints{}); err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestDecodeString_ErrorCheckingConstraints(t *testing.T) {
	data := []struct {
		desc        string
		value       string
		format      string
		constraints Constraints
	}{
		{"InvalidMinLength_UUID", "6fa459ea-ee8a-3ca4-894e-db77e160355e", stringUUID, Constraints{MinLength: 100}},
		{"InvalidMinLength_Email", "foo@bar.com", stringEmail, Constraints{MinLength: 100}},
		{"InvalidMinLength_URI", "http://google.com", stringURI, Constraints{MinLength: 100}},
		{"InvalidMaxLength_UUID", "6fa459ea-ee8a-3ca4-894e-db77e160355e", stringUUID, Constraints{MaxLength: 1}},
		{"InvalidMaxLength_Email", "foo@bar.com", stringEmail, Constraints{MaxLength: 1}},
		{"InvalidMaxLength_URI", "http://google.com", stringURI, Constraints{MaxLength: 1}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			if _, err := decodeString(d.format, d.value, d.constraints); err == nil {
				t.Fatalf("err want:err got:nil")
			}
		})
	}
}

func TestDecodeString_Success(t *testing.T) {
	var data = []struct {
		desc        string
		value       string
		format      string
		constraints Constraints
	}{
		{"URI", "http://google.com", stringURI, Constraints{MinLength: 1}},
		{"Email", "foo@bar.com", stringEmail, Constraints{MinLength: 1}},
		{"UUID", "C56A4180-65AA-42EC-A945-5FD21DEC0538", stringUUID, Constraints{MinLength: 36, MaxLength: 36}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			v, err := decodeString(d.format, d.value, d.constraints)
			if err != nil {
				t.Errorf("want:nil got:%q", err)
			}
			if v != d.value {
				t.Errorf("want:%s got:%s", d.value, v)
			}
		})
	}
}
