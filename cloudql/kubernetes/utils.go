package kubernetes

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"math"
	"strconv"
	"strings"
)

func v1TimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case metav1.Time:
		return v.ToUnstructured(), nil
	case *metav1.Time:
		if v == nil {
			return nil, nil
		}
		return v.ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("invalid time format %T! ", v)
	}
}

func v1MicroTimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case metav1.MicroTime:
		return metav1.NewTime(v.Time).ToUnstructured(), nil
	case *metav1.MicroTime:
		if v == nil {
			return nil, nil
		}
		return metav1.NewTime(v.Time).ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("invalid time format %T! ", v)
	}
}

func labelSelectorToString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	selector := d.Value.(*metav1.LabelSelector)

	ss, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return nil, err
	}

	return ss.String(), nil
}

func selectorMapToString(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("selectorMapToString")

	selector_map := d.Value.(map[string]string)

	if len(selector_map) == 0 {
		return nil, nil
	}

	selector_string := labels.SelectorFromSet(selector_map).String()

	return selector_string, nil
}

// normalizeCPUToMilliCores converts CPU quantities to millicores (m), rounding up if necessary.
func normalizeCPUToMilliCores(cpu string) (int64, error) {
	if strings.HasSuffix(cpu, "m") {
		// Already in millicores
		value, err := strconv.ParseFloat(strings.TrimSuffix(cpu, "m"), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid CPU value: %s", cpu)
		}
		return int64(math.Ceil(value)), nil
	}

	// Convert cores to millicores
	value, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid CPU value: %s", cpu)
	}

	milliCores := value * 1000
	return int64(math.Ceil(milliCores)), nil
}

// normalizeMemoryToBytes converts memory quantities to bytes, rounding up if necessary.
func normalizeMemoryToBytes(memory string) (int64, error) {
	// Set default value to arg value and unit to Bytes
	valuePart := memory
	unitPart := "B"
	for i, r := range memory {
		if r < '0' || r > '9' {
			valuePart = memory[:i]
			unitPart = memory[i:]
			break
		}
	}
	if valuePart == "" {
		return 0, fmt.Errorf("invalid memory value: %s", memory)
	}

	value, err := strconv.ParseFloat(valuePart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory value: %s", memory)
	}

	unitPart = strings.TrimSpace(unitPart)
	multiplier, exists := memoryUnits[unitPart]
	if !exists {
		return 0, fmt.Errorf("unknown unit: %s", unitPart)
	}

	bytes := value * multiplier
	return int64(math.Ceil(bytes)), nil
}

// Unit multipliers for memory
var memoryUnits = map[string]float64{
	"B":  1,
	"Ki": math.Pow(2, 10),
	"Mi": math.Pow(2, 20),
	"Gi": math.Pow(2, 30),
	"Ti": math.Pow(2, 40),
	"Pi": math.Pow(2, 50),
	"Ei": math.Pow(2, 60),
	"k":  1e3,
	"M":  1e6,
	"G":  1e9,
	"T":  1e12,
	"P":  1e15,
	"E":  1e18,
}

func mergeTags(labels map[string]string, annotations map[string]string) map[string]string {
	tags := make(map[string]string)
	for k, v := range annotations {
		tags[k] = v
	}
	for k, v := range labels {
		tags[k] = v
	}
	return tags
}
