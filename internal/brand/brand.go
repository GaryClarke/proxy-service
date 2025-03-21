package brand

// Brand represents a brand identifier.
type Brand string

const (
	GF  Brand = "gf"
	HEX Brand = "hex"
	RT  Brand = "rt"
	// Add additional brands as needed.
)

// FromPlatformBrandID returns the Brand corresponding to the given platformBrandID.
// It returns nil if no matching brand is found.
func FromPlatformBrandID(platformBrandID string) *Brand {
	switch platformBrandID {
	case "uk.co.bbc.goodfood2":
		b := GF
		return &b
	case "com.immediatemedia.historyextrarn":
		b := HEX
		return &b
	case "com.immediatemedia.radiotimes":
		b := RT
		return &b
	default:
		return nil
	}
}
