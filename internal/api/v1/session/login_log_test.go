package session

import "testing"

func TestNormalizeLoginIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{name: "ipv6 loopback", remoteAddr: "::1", want: loginIPLocalAddress},
		{name: "ipv6 loopback with port", remoteAddr: "[::1]:2019", want: loginIPLocalAddress},
		{name: "ipv4 with port", remoteAddr: "192.168.1.10:2019", want: "192.168.1.10"},
		{name: "ipv6 public", remoteAddr: "2001:db8::1", want: "2001:db8::1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeLoginIP(tt.remoteAddr); got != tt.want {
				t.Fatalf("normalizeLoginIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetLoginIPAreaFallback(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{name: "loopback", ip: loginIPLocalAddress, want: loginIPLocalArea},
		{name: "private ipv4", ip: "192.168.1.10", want: loginIPPrivateArea},
		{name: "private ipv6", ip: "fd00::1", want: loginIPPrivateArea},
		{name: "public ipv6", ip: "2001:db8::1", want: loginIPUnknownArea},
		{name: "invalid ip", ip: "not-an-ip", want: loginIPUnknownArea},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLoginIPArea(tt.ip); got != tt.want {
				t.Fatalf("getLoginIPArea() = %q, want %q", got, tt.want)
			}
		})
	}
}
