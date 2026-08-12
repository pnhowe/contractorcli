package cmd

import (
	"testing"
)

func scopeItem(percent float64, description string, elapsed string, remaining string) map[string]interface{} {
	parameters := map[string]interface{}{}
	if description != "" {
		parameters["description"] = description
	}
	if elapsed != "" {
		parameters["time_elapsed"] = elapsed
	}
	if remaining != "" {
		parameters["time_remaining"] = remaining
	}

	return map[string]interface{}{"percent": percent, "operation": "Scope", "parameters": parameters}
}

func TestJobStatusShort(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   *[]map[string]interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"empty", &[]map[string]interface{}{}, ""},
		{
			"just a scope",
			&[]map[string]interface{}{{"percent": 100.0, "operation": "Scope", "parameters": nil}},
			"100.00% Scope",
		},
		{ // the overall percent comes from the first entry, the label from the last
			"nested",
			&[]map[string]interface{}{
				scopeItem(12.5, "Install OS", "00:03", "00:12"),
				{"percent": 50.0, "operation": "Function", "parameters": map[string]interface{}{"module": "testing", "name": "remote", "dispatched": true}},
			},
			"12.50% testing.remote [dispatched]",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := jobStatusShort(test.status); got != test.expected {
				t.Errorf("got %q, want %q", got, test.expected)
			}
		})
	}
}

func TestJobStatus(t *testing.T) {
	status := &[]map[string]interface{}{
		scopeItem(12.5, "Install OS", "00:03", "00:12"),
		{"percent": 0.0, "operation": "Function", "parameters": map[string]interface{}{"module": "testing", "name": "remote", "dispatched": true}},
	}

	expected := " 12.50% Install OS Elapsed: 00:03 Remaining: 00:12\n" +
		"             " + "  0.00% testing.remote [dispatched]"

	if got := jobStatus(13, status); got != expected {
		t.Errorf("got:\n%q\nwant:\n%q", got, expected)
	}

	if got := jobStatus(13, nil); got != "" {
		t.Errorf("nil status: got %q, want empty", got)
	}
}

// a status entry that is not shaped the way we expect must degrade, not panic
func TestJobStatusMalformed(t *testing.T) {
	for _, test := range []struct {
		name string
		item map[string]interface{}
	}{
		{"empty", map[string]interface{}{}},
		{"null parameters", map[string]interface{}{"percent": 1.0, "operation": "Scope", "parameters": nil}},
		{"parameters not a map", map[string]interface{}{"percent": 1.0, "operation": "Scope", "parameters": "nope"}},
		{"percent not a number", map[string]interface{}{"percent": "nope", "operation": "Scope"}},
		{"operation not a string", map[string]interface{}{"operation": 7}},
		{"unknown operation", map[string]interface{}{"percent": 1.0, "operation": "While", "parameters": map[string]interface{}{"doing": "condition"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			jobStatusLine(test.item) // must not panic
			jobStatusShort(&[]map[string]interface{}{test.item})
		})
	}

	item := map[string]interface{}{"percent": 1.0, "operation": "While", "parameters": map[string]interface{}{"doing": "condition"}}
	if got := jobStatusLabel(item); got != "While(condition)" {
		t.Errorf("got %q, want %q", got, "While(condition)")
	}
}

func TestPrettyMap(t *testing.T) {
	for _, test := range []struct {
		name     string
		value    *map[string]interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"empty", &map[string]interface{}{}, ""},
		{
			"scalars are sorted",
			&map[string]interface{}{"zebra": "z", "alpha": "a", "count": 3.0},
			"alpha = a\ncount = 3\nzebra = z",
		},
		{ // large numbers must not come out as 1.37438953472e+11
			"big number",
			&map[string]interface{}{"total_ram": 137438953472.0, "percent": 12.5},
			"percent = 12.5\ntotal_ram = 137438953472",
		},
		{
			"nested maps flatten to dotted paths",
			&map[string]interface{}{"hardware": map[string]interface{}{"dmi": map[string]interface{}{"vendor": "Dell Inc."}}},
			"hardware.dmi.vendor = Dell Inc.",
		},
		{ // a list of scalars stays inline, a list of maps gets a line per element
			"lists",
			&map[string]interface{}{
				"tags":  []interface{}{"a", "b"},
				"disks": []interface{}{map[string]interface{}{"model": "ST2000A"}},
			},
			"disks.0.model = ST2000A\ntags = [a, b]",
		},
		{
			"empty containers and null",
			&map[string]interface{}{"a": map[string]interface{}{}, "b": []interface{}{}, "c": nil},
			"a = {}\nb = []\nc = <null>",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := prettyMap(0, test.value); got != test.expected {
				t.Errorf("got:\n%q\nwant:\n%q", got, test.expected)
			}
		})
	}
}

func TestPrettyMapIndent(t *testing.T) {
	value := &map[string]interface{}{"a": "1", "b": "2"}
	if got, expected := prettyMap(4, value), "a = 1\n    b = 2"; got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestOutputKVIsStable(t *testing.T) {
	// outputKV used to range the map directly, so key order changed between runs -- prettyMap
	// sorts, which is what makes two runs of "site config" diffable
	value := map[string]interface{}{"zebra": "z", "alpha": "a", "nested": map[string]interface{}{"b": 2.0, "a": 1.0}}

	expected := "alpha = a\nnested.a = 1\nnested.b = 2\nzebra = z"
	for i := 0; i < 20; i++ {
		if got := prettyMap(0, &value); got != expected {
			t.Fatalf("iteration %d: got:\n%q\nwant:\n%q", i, got, expected)
		}
	}
}
