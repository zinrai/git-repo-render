package main

import "testing"

func TestIsBinaryContent(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "empty data is not binary",
			data: []byte{},
			want: false,
		},
		{
			name: "plain text",
			data: []byte("Hello, world!\n"),
			want: false,
		},
		{
			name: "Go source code",
			data: []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n"),
			want: false,
		},
		{
			name: "JSON content",
			data: []byte(`{"key": "value", "number": 42}`),
			want: false,
		},
		{
			name: "PNG header bytes",
			data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00},
			want: true,
		},
		{
			name: "GIF header bytes",
			data: []byte("GIF89a" + string([]byte{0x01, 0x00, 0x01, 0x00, 0x80, 0x00, 0x00})),
			want: true,
		},
		{
			name: "null bytes indicate binary",
			data: []byte{0x00, 0x00, 0x00, 0x00},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBinaryContent(tt.data)
			if got != tt.want {
				t.Errorf("isBinaryContent(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}
