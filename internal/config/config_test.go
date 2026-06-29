package config

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr string
	}{
		{
			name: "all set",
			env: map[string]string{
				"GCC":           "/usr/bin/gcc",
				"CLANG":         "/usr/bin/clang",
				"CLANG_TBLGEN":  "/usr/bin/clang-tblgen",
			},
			want: Config{
				GCC:         "/usr/bin/gcc",
				Clang:       "/usr/bin/clang",
				ClangTblgen: "/usr/bin/clang-tblgen",
			},
		},
		{
			name: "missing gcc",
			env: map[string]string{
				"CLANG":        "/usr/bin/clang",
				"CLANG_TBLGEN": "/usr/bin/clang-tblgen",
			},
			wantErr: `GCC`,
		},
		{
			name: "missing clang",
			env: map[string]string{
				"GCC":          "/usr/bin/gcc",
				"CLANG_TBLGEN": "/usr/bin/clang-tblgen",
			},
			wantErr: `CLANG`,
		},
		{
			name: "missing clang tblgen",
			env: map[string]string{
				"GCC":   "/usr/bin/gcc",
				"CLANG": "/usr/bin/clang",
			},
			wantErr: `CLANG_TBLGEN`,
		},
		{
			name: "empty gcc",
			env: map[string]string{
				"GCC":          "",
				"CLANG":        "/usr/bin/clang",
				"CLANG_TBLGEN": "/usr/bin/clang-tblgen",
			},
			wantErr: `GCC`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, val := range tt.env {
				t.Setenv(key, val)
			}
			for _, key := range []string{"GCC", "CLANG", "CLANG_TBLGEN"} {
				if _, ok := tt.env[key]; !ok {
					t.Setenv(key, "")
				}
			}

			got, err := Load()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("Load() error = nil, want error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Load() error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Load() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
