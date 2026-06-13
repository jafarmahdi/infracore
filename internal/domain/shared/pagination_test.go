package shared

import "testing"

func TestPageNormalized(t *testing.T) {
	tests := []struct {
		name string
		page Page
		want Page
	}{
		{name: "zero value", page: Page{}, want: Page{Number: 1, Size: DefaultPageSize}},
		{name: "negative values", page: Page{Number: -2, Size: -10}, want: Page{Number: 1, Size: DefaultPageSize}},
		{name: "caps oversized pages", page: Page{Number: 2, Size: 1000}, want: Page{Number: 2, Size: MaxPageSize}},
		{name: "preserves valid values", page: Page{Number: 3, Size: 50}, want: Page{Number: 3, Size: 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.page.Normalized(); got != tt.want {
				t.Fatalf("Normalized() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestNewPageResultNormalizesZeroValuePage(t *testing.T) {
	result := NewPageResult([]string{"asset-1"}, 21, Page{})

	if result.Page != 1 || result.PageSize != DefaultPageSize {
		t.Fatalf("unexpected normalized page metadata: %+v", result)
	}
	if result.TotalPages != 2 {
		t.Fatalf("TotalPages = %d, want 2", result.TotalPages)
	}
}

func TestPageOffsetUsesNormalizedValues(t *testing.T) {
	if got := (Page{}).Offset(); got != 0 {
		t.Fatalf("Offset() = %d, want 0", got)
	}
	if got := (Page{Number: 3, Size: 25}).Offset(); got != 50 {
		t.Fatalf("Offset() = %d, want 50", got)
	}
}
