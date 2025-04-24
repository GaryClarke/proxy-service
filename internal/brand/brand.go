package brand

import "fmt"

// Brand represents a brand identifier.
type Brand string

const (
	GF  Brand = "gf"
	HEX Brand = "hex"
	RT  Brand = "rt"
	// Add additional brands as needed.
)

// FromPlatformBrandID returns the Brand corresponding to the given platformBrandID.
// It returns an error if no matching brand is found.
func FromPlatformBrandID(platformBrandID string) (Brand, error) {
	switch platformBrandID {
	case "uk.co.bbc.goodfood2":
		b := GF
		return b, nil
	case "com.immediatemedia.historyextrarn":
		b := HEX
		return b, nil
	case "com.immediatemedia.radiotimes":
		b := RT
		return b, nil
	default:
		return "", fmt.Errorf("unknown platform brand id: %s", platformBrandID)
	}
}

// ValidateBrand todo - add supported brand check later if required
func ValidateBrand(brand Brand) error {
	return nil
}
